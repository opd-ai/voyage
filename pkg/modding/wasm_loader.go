package modding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Capability represents a permission that can be granted to a WASM mod.
type Capability uint32

const (
	// CapReadEvents allows the mod to read existing event data.
	CapReadEvents Capability = 1 << iota
	// CapWriteEvents allows the mod to add new events.
	CapWriteEvents
	// CapReadGenres allows the mod to read genre configurations.
	CapReadGenres
	// CapWriteGenres allows the mod to add custom genres.
	CapWriteGenres
	// CapReadResources allows the mod to read resource data.
	CapReadResources
	// CapWriteResources allows the mod to modify resource values.
	CapWriteResources
	// CapReadCrew allows the mod to read crew member data.
	CapReadCrew
	// CapModifyCrew allows the mod to modify crew members.
	CapModifyCrew
	// CapTriggerEvents allows the mod to trigger custom events.
	CapTriggerEvents
	// CapAccessRNG allows the mod to access the game's RNG.
	CapAccessRNG
)

// AllCapabilities is a convenience mask for all permissions.
const AllCapabilities = CapReadEvents | CapWriteEvents | CapReadGenres |
	CapWriteGenres | CapReadResources | CapWriteResources |
	CapReadCrew | CapModifyCrew | CapTriggerEvents | CapAccessRNG

// MinimalCapabilities provides read-only access.
const MinimalCapabilities = CapReadEvents | CapReadGenres | CapReadResources | CapReadCrew

// WASMModErrors for the WASM mod loader.
var (
	ErrWASMCompileFailed   = errors.New("failed to compile WASM module")
	ErrWASMInstantiate     = errors.New("failed to instantiate WASM module")
	ErrCapabilityDenied    = errors.New("capability denied")
	ErrModExecutionTimeout = errors.New("mod execution timed out")
	ErrInvalidWASMExport   = errors.New("WASM module missing required exports")
	ErrModPanicked         = errors.New("mod execution panicked")
)

// WASMConfig configures the WASM sandbox environment.
type WASMConfig struct {
	// MaxMemoryPages limits memory allocation (1 page = 64KB).
	MaxMemoryPages uint32
	// ExecutionTimeout is the maximum time a mod function can run.
	ExecutionTimeout time.Duration
	// Capabilities is the set of permissions granted to the mod.
	Capabilities Capability
	// DebugMode enables verbose logging.
	DebugMode bool
}

// DefaultWASMConfig returns safe defaults for WASM execution.
func DefaultWASMConfig() WASMConfig {
	return WASMConfig{
		MaxMemoryPages:   256, // 16MB max
		ExecutionTimeout: 5 * time.Second,
		Capabilities:     MinimalCapabilities,
		DebugMode:        false,
	}
}

// WASMMod represents a loaded WASM mod.
type WASMMod struct {
	ID           string
	Name         string
	Version      string
	Author       string
	Description  string
	Capabilities Capability
	LoadedAt     time.Time
	Enabled      bool

	runtime  wazero.Runtime
	module   api.Module
	config   WASMConfig
	mu       sync.RWMutex
	hostData *hostData
}

// hostData holds shared data accessible to the WASM module.
type hostData struct {
	mu           sync.RWMutex
	events       []EventDef
	genres       []GenreDef
	resources    map[string]float64
	outputBuffer []byte
	inputBuffer  []byte
}

// WASMLoader manages WASM mod loading and execution.
type WASMLoader struct {
	mu       sync.RWMutex
	mods     map[string]*WASMMod
	modOrder []string
	config   WASMConfig
}

// NewWASMLoader creates a new WASM mod loader.
func NewWASMLoader(config WASMConfig) *WASMLoader {
	return &WASMLoader{
		mods:     make(map[string]*WASMMod),
		modOrder: make([]string, 0),
		config:   config,
	}
}

// LoadFromFile loads a WASM mod from a file.
func (l *WASMLoader) LoadFromFile(path string, caps Capability) (*WASMMod, error) {
	wasmBytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return l.LoadFromBytes(wasmBytes, caps)
}

// LoadFromReader loads a WASM mod from an io.Reader.
func (l *WASMLoader) LoadFromReader(r io.Reader, caps Capability) (*WASMMod, error) {
	wasmBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return l.LoadFromBytes(wasmBytes, caps)
}

// LoadFromBytes loads and compiles a WASM mod from bytes.
func (l *WASMLoader) LoadFromBytes(wasmBytes []byte, caps Capability) (*WASMMod, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ctx := context.Background()
	runtime, hd, err := l.initializeRuntime(ctx, caps)
	if err != nil {
		return nil, err
	}

	module, err := l.compileAndInstantiate(ctx, runtime, wasmBytes)
	if err != nil {
		runtime.Close(ctx)
		return nil, err
	}

	mod := l.createWASMMod(runtime, module, hd, caps)
	if err := l.finalizeMod(ctx, runtime, mod); err != nil {
		return nil, err
	}

	return mod, nil
}

// initializeRuntime creates the WASM runtime with host functions.
func (l *WASMLoader) initializeRuntime(ctx context.Context, caps Capability) (wazero.Runtime, *hostData, error) {
	runtimeConfig := wazero.NewRuntimeConfig().
		WithMemoryLimitPages(l.config.MaxMemoryPages)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		runtime.Close(ctx)
		return nil, nil, fmt.Errorf("%w: %v", ErrWASMInstantiate, err)
	}

	hd := newHostData()
	if _, err := l.registerHostFunctions(ctx, runtime, caps, hd); err != nil {
		runtime.Close(ctx)
		return nil, nil, err
	}

	return runtime, hd, nil
}

// newHostData creates initialized host data for a mod.
func newHostData() *hostData {
	return &hostData{
		events:       make([]EventDef, 0),
		genres:       make([]GenreDef, 0),
		resources:    make(map[string]float64),
		outputBuffer: make([]byte, 0, 4096),
		inputBuffer:  make([]byte, 0, 4096),
	}
}

// compileAndInstantiate compiles WASM bytes and creates a module instance.
func (l *WASMLoader) compileAndInstantiate(ctx context.Context, runtime wazero.Runtime, wasmBytes []byte) (api.Module, error) {
	compiled, err := runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWASMCompileFailed, err)
	}

	modConfig := wazero.NewModuleConfig().
		WithStdout(io.Discard).
		WithStderr(io.Discard).
		WithStartFunctions()

	module, err := runtime.InstantiateModule(ctx, compiled, modConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWASMInstantiate, err)
	}

	return module, nil
}

// createWASMMod constructs a WASMMod from initialized components.
func (l *WASMLoader) createWASMMod(runtime wazero.Runtime, module api.Module, hd *hostData, caps Capability) *WASMMod {
	return &WASMMod{
		runtime:      runtime,
		module:       module,
		config:       l.config,
		Capabilities: caps,
		LoadedAt:     time.Now(),
		Enabled:      true,
		hostData:     hd,
	}
}

// finalizeMod loads metadata and registers the mod.
func (l *WASMLoader) finalizeMod(ctx context.Context, runtime wazero.Runtime, mod *WASMMod) error {
	if err := mod.loadMetadata(ctx); err != nil {
		runtime.Close(ctx)
		return err
	}

	if _, exists := l.mods[mod.ID]; exists {
		runtime.Close(ctx)
		return ErrModAlreadyLoaded
	}

	l.mods[mod.ID] = mod
	l.modOrder = append(l.modOrder, mod.ID)
	return nil
}

// registerHostFunctions adds capability-guarded host functions.
func (l *WASMLoader) registerHostFunctions(ctx context.Context, runtime wazero.Runtime, caps Capability, hd *hostData) (api.Module, error) {
	builder := runtime.NewHostModuleBuilder("voyage")

	registerEventFunctions(builder, caps, hd)
	registerGenreFunctions(builder, caps, hd)
	registerResourceFunctions(builder, caps, hd)
	registerUtilityFunctions(builder, caps, hd)

	return builder.Instantiate(ctx)
}

// registerEventFunctions adds event reading and writing host functions.
func registerEventFunctions(builder wazero.HostModuleBuilder, caps Capability, hd *hostData) {
	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
			if caps&CapReadEvents == 0 {
				return 0
			}
			hd.mu.RLock()
			defer hd.mu.RUnlock()
			data, _ := json.Marshal(hd.events)
			return writeToMemory(m, ptr, len, data)
		}).
		Export("voyage_read_events")

	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
			if caps&CapWriteEvents == 0 {
				return 0
			}
			data := readFromMemory(m, ptr, len)
			var event EventDef
			if json.Unmarshal(data, &event) == nil {
				hd.mu.Lock()
				hd.events = append(hd.events, event)
				hd.mu.Unlock()
				return 1
			}
			return 0
		}).
		Export("voyage_add_event")
}

// registerGenreFunctions adds genre reading and writing host functions.
func registerGenreFunctions(builder wazero.HostModuleBuilder, caps Capability, hd *hostData) {
	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
			if caps&CapReadGenres == 0 {
				return 0
			}
			hd.mu.RLock()
			defer hd.mu.RUnlock()
			data, _ := json.Marshal(hd.genres)
			return writeToMemory(m, ptr, len, data)
		}).
		Export("voyage_read_genres")

	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
			if caps&CapWriteGenres == 0 {
				return 0
			}
			data := readFromMemory(m, ptr, len)
			var genre GenreDef
			if json.Unmarshal(data, &genre) == nil {
				hd.mu.Lock()
				hd.genres = append(hd.genres, genre)
				hd.mu.Unlock()
				return 1
			}
			return 0
		}).
		Export("voyage_add_genre")
}

// registerResourceFunctions adds resource reading host function.
func registerResourceFunctions(builder wazero.HostModuleBuilder, caps Capability, hd *hostData) {
	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
			if caps&CapReadResources == 0 {
				return 0
			}
			hd.mu.RLock()
			defer hd.mu.RUnlock()
			data, _ := json.Marshal(hd.resources)
			return writeToMemory(m, ptr, len, data)
		}).
		Export("voyage_read_resources")
}

// registerUtilityFunctions adds log and capability check host functions.
func registerUtilityFunctions(builder wazero.HostModuleBuilder, caps Capability, hd *hostData) {
	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) {
			data := readFromMemory(m, ptr, len)
			hd.mu.Lock()
			hd.outputBuffer = append(hd.outputBuffer, data...)
			hd.mu.Unlock()
		}).
		Export("voyage_log")

	builder.NewFunctionBuilder().
		WithFunc(func(cap uint32) uint32 {
			if caps&Capability(cap) != 0 {
				return 1
			}
			return 0
		}).
		Export("voyage_has_capability")
}

// loadMetadata extracts mod information from WASM exports.
func (m *WASMMod) loadMetadata(ctx context.Context) error {
	// Look for standard export functions
	getID := m.module.ExportedFunction("mod_get_id")
	getName := m.module.ExportedFunction("mod_get_name")
	getVersion := m.module.ExportedFunction("mod_get_version")
	getAuthor := m.module.ExportedFunction("mod_get_author")

	// ID is required
	if getID == nil {
		return fmt.Errorf("%w: mod_get_id", ErrInvalidWASMExport)
	}

	// Call with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.ExecutionTimeout)
	defer cancel()

	m.ID = m.callStringFunc(timeoutCtx, getID)
	if m.ID == "" {
		return ErrMissingID
	}

	m.Name = m.callStringFunc(timeoutCtx, getName)
	if m.Name == "" {
		m.Name = m.ID
	}

	m.Version = m.callStringFunc(timeoutCtx, getVersion)
	if m.Version == "" {
		m.Version = "0.0.0"
	}

	m.Author = m.callStringFunc(timeoutCtx, getAuthor)

	return nil
}

// callStringFunc calls a WASM function that returns a string pointer and length.
func (m *WASMMod) callStringFunc(ctx context.Context, fn api.Function) string {
	if fn == nil {
		return ""
	}

	results, err := fn.Call(ctx)
	if err != nil || len(results) < 2 {
		return ""
	}

	ptr := uint32(results[0])
	length := uint32(results[1])

	if length == 0 {
		return ""
	}

	mem := m.module.Memory()
	if mem == nil {
		return ""
	}

	data, ok := mem.Read(ptr, length)
	if !ok {
		return ""
	}

	return string(data)
}

// Initialize calls the mod's init function if present.
func (m *WASMMod) Initialize(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.module == nil {
		return nil
	}

	initFn := m.module.ExportedFunction("mod_init")
	if initFn == nil {
		return nil // Init is optional
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.ExecutionTimeout)
	defer cancel()

	_, err := initFn.Call(timeoutCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrModExecutionTimeout
		}
		return fmt.Errorf("%w: %v", ErrModPanicked, err)
	}

	return nil
}

// OnTurnStart calls the mod's turn start hook.
func (m *WASMMod) OnTurnStart(ctx context.Context, turn int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.module == nil {
		return nil
	}

	hookFn := m.module.ExportedFunction("mod_on_turn_start")
	if hookFn == nil {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.ExecutionTimeout)
	defer cancel()

	_, err := hookFn.Call(timeoutCtx, uint64(turn))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrModExecutionTimeout
		}
		return fmt.Errorf("%w: %v", ErrModPanicked, err)
	}

	return nil
}

// OnEvent calls the mod's event hook.
func (m *WASMMod) OnEvent(ctx context.Context, eventCategory string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.module == nil {
		return nil
	}

	hookFn := m.module.ExportedFunction("mod_on_event")
	if hookFn == nil {
		return nil
	}

	// Write event category to memory
	mem := m.module.Memory()
	if mem == nil {
		return nil
	}

	// Use a fixed buffer location for input
	const inputPtr = 1024
	categoryBytes := []byte(eventCategory)
	mem.Write(inputPtr, categoryBytes)

	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.ExecutionTimeout)
	defer cancel()

	_, err := hookFn.Call(timeoutCtx, inputPtr, uint64(len(categoryBytes)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrModExecutionTimeout
		}
		return fmt.Errorf("%w: %v", ErrModPanicked, err)
	}

	return nil
}

// GetAddedEvents returns events added by this mod.
func (m *WASMMod) GetAddedEvents() []EventDef {
	m.hostData.mu.RLock()
	defer m.hostData.mu.RUnlock()
	result := make([]EventDef, len(m.hostData.events))
	copy(result, m.hostData.events)
	return result
}

// GetAddedGenres returns genres added by this mod.
func (m *WASMMod) GetAddedGenres() []GenreDef {
	m.hostData.mu.RLock()
	defer m.hostData.mu.RUnlock()
	result := make([]GenreDef, len(m.hostData.genres))
	copy(result, m.hostData.genres)
	return result
}

// GetLogs returns log output from the mod.
func (m *WASMMod) GetLogs() string {
	m.hostData.mu.RLock()
	defer m.hostData.mu.RUnlock()
	return string(m.hostData.outputBuffer)
}

// ClearLogs clears the log buffer.
func (m *WASMMod) ClearLogs() {
	m.hostData.mu.Lock()
	defer m.hostData.mu.Unlock()
	m.hostData.outputBuffer = m.hostData.outputBuffer[:0]
}

// Close releases all resources for this mod.
func (m *WASMMod) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.runtime != nil {
		return m.runtime.Close(context.Background())
	}
	return nil
}

// HasCapability checks if a capability is granted.
func (m *WASMMod) HasCapability(cap Capability) bool {
	return m.Capabilities&cap != 0
}

// Get returns a loaded WASM mod by ID.
func (l *WASMLoader) Get(modID string) (*WASMMod, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	mod, exists := l.mods[modID]
	if !exists {
		return nil, ErrModNotFound
	}

	return mod, nil
}

// Unload removes a WASM mod.
func (l *WASMLoader) Unload(modID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	mod, exists := l.mods[modID]
	if !exists {
		return ErrModNotFound
	}

	mod.Close()
	delete(l.mods, modID)

	// Remove from order
	for i, id := range l.modOrder {
		if id == modID {
			l.modOrder = append(l.modOrder[:i], l.modOrder[i+1:]...)
			break
		}
	}

	return nil
}

// List returns all loaded WASM mods in load order.
func (l *WASMLoader) List() []*WASMMod {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make([]*WASMMod, 0, len(l.modOrder))
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists {
			result = append(result, mod)
		}
	}

	return result
}

// ListEnabled returns only enabled WASM mods.
func (l *WASMLoader) ListEnabled() []*WASMMod {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make([]*WASMMod, 0)
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			result = append(result, mod)
		}
	}

	return result
}

// Count returns the number of loaded WASM mods.
func (l *WASMLoader) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.mods)
}

// Close releases all loaded WASM mods.
func (l *WASMLoader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lastErr error
	for _, mod := range l.mods {
		if err := mod.Close(); err != nil {
			lastErr = err
		}
	}

	l.mods = make(map[string]*WASMMod)
	l.modOrder = make([]string, 0)

	return lastErr
}

// GetAllEvents returns all events from all enabled WASM mods.
func (l *WASMLoader) GetAllEvents() []EventDef {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var events []EventDef
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			events = append(events, mod.GetAddedEvents()...)
		}
	}

	return events
}

// GetAllGenres returns all genres from all enabled WASM mods.
func (l *WASMLoader) GetAllGenres() []GenreDef {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var genres []GenreDef
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			genres = append(genres, mod.GetAddedGenres()...)
		}
	}

	return genres
}

// Helper functions for memory access

func readFromMemory(m api.Module, ptr, length uint32) []byte {
	mem := m.Memory()
	if mem == nil {
		return nil
	}

	data, ok := mem.Read(ptr, length)
	if !ok {
		return nil
	}

	result := make([]byte, len(data))
	copy(result, data)
	return result
}

func writeToMemory(m api.Module, ptr, maxLen uint32, data []byte) uint32 {
	mem := m.Memory()
	if mem == nil {
		return 0
	}

	writeLen := uint32(len(data))
	if writeLen > maxLen {
		writeLen = maxLen
	}

	if !mem.Write(ptr, data[:writeLen]) {
		return 0
	}

	return writeLen
}

// capabilityNames maps each capability flag to its string name.
var capabilityNames = []struct {
	cap  Capability
	name string
}{
	{CapReadEvents, "read_events"},
	{CapWriteEvents, "write_events"},
	{CapReadGenres, "read_genres"},
	{CapWriteGenres, "write_genres"},
	{CapReadResources, "read_resources"},
	{CapWriteResources, "write_resources"},
	{CapReadCrew, "read_crew"},
	{CapModifyCrew, "modify_crew"},
	{CapTriggerEvents, "trigger_events"},
	{CapAccessRNG, "access_rng"},
}

// CapabilityString returns a human-readable capability name.
func CapabilityString(cap Capability) string {
	names := collectCapabilityNames(cap)
	if len(names) == 0 {
		return "none"
	}
	return joinNames(names)
}

// collectCapabilityNames returns the names of all capabilities set in cap.
func collectCapabilityNames(cap Capability) []string {
	names := make([]string, 0, len(capabilityNames))
	for _, cn := range capabilityNames {
		if cap&cn.cap != 0 {
			names = append(names, cn.name)
		}
	}
	return names
}

// joinNames concatenates names with comma separator.
func joinNames(names []string) string {
	result := names[0]
	for i := 1; i < len(names); i++ {
		result += ", " + names[i]
	}
	return result
}
