package modding

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestDefaultWASMConfig(t *testing.T) {
	config := DefaultWASMConfig()

	if config.MaxMemoryPages != 256 {
		t.Errorf("expected MaxMemoryPages 256, got %d", config.MaxMemoryPages)
	}
	if config.ExecutionTimeout != 5*time.Second {
		t.Errorf("expected ExecutionTimeout 5s, got %v", config.ExecutionTimeout)
	}
	if config.Capabilities != MinimalCapabilities {
		t.Errorf("expected MinimalCapabilities, got %d", config.Capabilities)
	}
	if config.DebugMode {
		t.Error("expected DebugMode false")
	}
}

func TestCapabilityConstants(t *testing.T) {
	// Ensure capabilities are unique powers of 2
	caps := []Capability{
		CapReadEvents, CapWriteEvents,
		CapReadGenres, CapWriteGenres,
		CapReadResources, CapWriteResources,
		CapReadCrew, CapModifyCrew,
		CapTriggerEvents, CapAccessRNG,
	}

	seen := make(map[Capability]bool)
	for _, cap := range caps {
		if seen[cap] {
			t.Errorf("duplicate capability value: %d", cap)
		}
		seen[cap] = true

		// Check it's a power of 2
		if cap&(cap-1) != 0 {
			t.Errorf("capability %d is not a power of 2", cap)
		}
	}
}

func TestMinimalCapabilities(t *testing.T) {
	// Minimal should have read caps but not write caps
	if MinimalCapabilities&CapReadEvents == 0 {
		t.Error("MinimalCapabilities should include CapReadEvents")
	}
	if MinimalCapabilities&CapWriteEvents != 0 {
		t.Error("MinimalCapabilities should not include CapWriteEvents")
	}
	if MinimalCapabilities&CapReadGenres == 0 {
		t.Error("MinimalCapabilities should include CapReadGenres")
	}
	if MinimalCapabilities&CapWriteGenres != 0 {
		t.Error("MinimalCapabilities should not include CapWriteGenres")
	}
}

func TestAllCapabilities(t *testing.T) {
	// All capabilities should include everything
	caps := []Capability{
		CapReadEvents, CapWriteEvents,
		CapReadGenres, CapWriteGenres,
		CapReadResources, CapWriteResources,
		CapReadCrew, CapModifyCrew,
		CapTriggerEvents, CapAccessRNG,
	}

	for _, cap := range caps {
		if AllCapabilities&cap == 0 {
			t.Errorf("AllCapabilities missing capability: %d", cap)
		}
	}
}

func TestNewWASMLoader(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	if loader == nil {
		t.Fatal("NewWASMLoader returned nil")
	}
	if loader.mods == nil {
		t.Error("mods map not initialized")
	}
	if loader.modOrder == nil {
		t.Error("modOrder slice not initialized")
	}
	if loader.Count() != 0 {
		t.Error("expected empty loader")
	}
}

func TestWASMLoaderLoadMissingFile(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	_, err := loader.LoadFromFile("/nonexistent/path.wasm", MinimalCapabilities)
	if err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}

func TestWASMLoaderList(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	mods := loader.List()
	if len(mods) != 0 {
		t.Errorf("expected empty list, got %d", len(mods))
	}

	enabled := loader.ListEnabled()
	if len(enabled) != 0 {
		t.Errorf("expected empty enabled list, got %d", len(enabled))
	}
}

func TestWASMLoaderGetNotFound(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	_, err := loader.Get("nonexistent")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

func TestWASMLoaderUnloadNotFound(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	err := loader.Unload("nonexistent")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

func TestWASMLoaderClose(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	err := loader.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if loader.Count() != 0 {
		t.Error("expected empty loader after close")
	}
}

func TestWASMLoaderGetAllEvents(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	events := loader.GetAllEvents()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestWASMLoaderGetAllGenres(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	genres := loader.GetAllGenres()
	if len(genres) != 0 {
		t.Errorf("expected no genres, got %d", len(genres))
	}
}

func TestCapabilityString(t *testing.T) {
	testCases := []struct {
		cap      Capability
		expected string
	}{
		{0, "none"},
		{CapReadEvents, "read_events"},
		{CapWriteEvents, "write_events"},
		{CapReadEvents | CapWriteEvents, "read_events, write_events"},
		{MinimalCapabilities, "read_events, read_genres, read_resources, read_crew"},
	}

	for _, tc := range testCases {
		result := CapabilityString(tc.cap)
		if result != tc.expected {
			t.Errorf("CapabilityString(%d) = %q, expected %q", tc.cap, result, tc.expected)
		}
	}
}

func TestHostData(t *testing.T) {
	hd := &hostData{
		events:       make([]EventDef, 0),
		genres:       make([]GenreDef, 0),
		resources:    make(map[string]float64),
		outputBuffer: make([]byte, 0),
		inputBuffer:  make([]byte, 0),
	}

	// Test adding events
	hd.events = append(hd.events, EventDef{Title: "Test"})
	if len(hd.events) != 1 {
		t.Error("failed to add event")
	}

	// Test adding genres
	hd.genres = append(hd.genres, GenreDef{ID: "test"})
	if len(hd.genres) != 1 {
		t.Error("failed to add genre")
	}

	// Test resources
	hd.resources["food"] = 100
	if hd.resources["food"] != 100 {
		t.Error("failed to set resource")
	}
}

func TestWASMModHasCapability(t *testing.T) {
	mod := &WASMMod{
		Capabilities: CapReadEvents | CapWriteEvents,
	}

	if !mod.HasCapability(CapReadEvents) {
		t.Error("expected mod to have CapReadEvents")
	}
	if !mod.HasCapability(CapWriteEvents) {
		t.Error("expected mod to have CapWriteEvents")
	}
	if mod.HasCapability(CapReadGenres) {
		t.Error("expected mod to not have CapReadGenres")
	}
}

func TestWASMModGetAddedEventsEmpty(t *testing.T) {
	hd := &hostData{
		events: make([]EventDef, 0),
	}
	mod := &WASMMod{
		hostData: hd,
	}

	events := mod.GetAddedEvents()
	if len(events) != 0 {
		t.Errorf("expected empty events, got %d", len(events))
	}
}

func TestWASMModGetAddedGenresEmpty(t *testing.T) {
	hd := &hostData{
		genres: make([]GenreDef, 0),
	}
	mod := &WASMMod{
		hostData: hd,
	}

	genres := mod.GetAddedGenres()
	if len(genres) != 0 {
		t.Errorf("expected empty genres, got %d", len(genres))
	}
}

func TestWASMModGetLogs(t *testing.T) {
	hd := &hostData{
		outputBuffer: []byte("test log output"),
	}
	mod := &WASMMod{
		hostData: hd,
	}

	logs := mod.GetLogs()
	if logs != "test log output" {
		t.Errorf("expected 'test log output', got %q", logs)
	}
}

func TestWASMModClearLogs(t *testing.T) {
	hd := &hostData{
		outputBuffer: []byte("test log output"),
	}
	mod := &WASMMod{
		hostData: hd,
	}

	mod.ClearLogs()
	if len(mod.hostData.outputBuffer) != 0 {
		t.Error("expected empty output buffer after clear")
	}
}

func TestWASMModCloseNilRuntime(t *testing.T) {
	mod := &WASMMod{
		runtime: nil,
	}

	err := mod.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWASMConfigCustom(t *testing.T) {
	config := WASMConfig{
		MaxMemoryPages:   128,
		ExecutionTimeout: 10 * time.Second,
		Capabilities:     AllCapabilities,
		DebugMode:        true,
	}

	loader := NewWASMLoader(config)
	if loader.config.MaxMemoryPages != 128 {
		t.Error("custom config not applied")
	}
	if loader.config.ExecutionTimeout != 10*time.Second {
		t.Error("custom timeout not applied")
	}
}

func TestWASMModInitializeNoInit(t *testing.T) {
	hd := &hostData{
		events: make([]EventDef, 0),
	}
	mod := &WASMMod{
		hostData: hd,
		config:   DefaultWASMConfig(),
	}

	// With nil module, Initialize should return nil
	err := mod.Initialize(context.Background())
	if err != nil {
		t.Errorf("expected nil error for nil module, got %v", err)
	}
}

func TestWASMModOnTurnStartNoHook(t *testing.T) {
	hd := &hostData{
		events: make([]EventDef, 0),
	}
	mod := &WASMMod{
		hostData: hd,
		config:   DefaultWASMConfig(),
	}

	// With nil module, OnTurnStart should return nil
	err := mod.OnTurnStart(context.Background(), 1)
	if err != nil {
		t.Errorf("expected nil error for nil module, got %v", err)
	}
}

func TestWASMModOnEventNoHook(t *testing.T) {
	hd := &hostData{
		events: make([]EventDef, 0),
	}
	mod := &WASMMod{
		hostData: hd,
		config:   DefaultWASMConfig(),
	}

	// With nil module, OnEvent should return nil
	err := mod.OnEvent(context.Background(), "encounter")
	if err != nil {
		t.Errorf("expected nil error for nil module, got %v", err)
	}
}

// TestLoadFromBytesInvalidWASM tests that invalid WASM bytecode is rejected.
func TestLoadFromBytesInvalidWASM(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	// Invalid WASM bytecode (not starting with magic bytes)
	invalidBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	_, err := loader.LoadFromBytes(invalidBytes, MinimalCapabilities)
	if err == nil {
		t.Error("expected error for invalid WASM bytecode")
	}
	// Should be a compile error
	if !errors.Is(err, ErrWASMCompileFailed) {
		t.Errorf("expected ErrWASMCompileFailed, got %v", err)
	}
}

// TestLoadFromBytesEmptyWASM tests that empty WASM bytecode is rejected.
func TestLoadFromBytesEmptyWASM(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	emptyBytes := []byte{}
	_, err := loader.LoadFromBytes(emptyBytes, MinimalCapabilities)
	if err == nil {
		t.Error("expected error for empty WASM bytecode")
	}
}

// TestLoadFromReaderError tests LoadFromReader with a failing reader.
func TestLoadFromReaderError(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	// Create a reader that always fails
	failReader := &failingReader{}
	_, err := loader.LoadFromReader(failReader, MinimalCapabilities)
	if err == nil {
		t.Error("expected error from failing reader")
	}
}

// failingReader is an io.Reader that always returns an error.
type failingReader struct{}

func (f *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

// TestWASMLoaderGetNotFoundAfterUnload tests getting a mod after unloading.
func TestWASMLoaderGetNotFoundAfterUnload(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	// Manually add a mod
	mod := &WASMMod{ID: "test-mod", Name: "Test", hostData: newHostData()}
	loader.mu.Lock()
	loader.mods["test-mod"] = mod
	loader.modOrder = append(loader.modOrder, "test-mod")
	loader.mu.Unlock()

	// Unload it
	err := loader.Unload("test-mod")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test getting a non-existent mod
	_, err = loader.Get("test-mod")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

// TestCapabilityDenialScenarios tests that capability checks work correctly.
func TestCapabilityDenialScenarios(t *testing.T) {
	// Test mod with no write capabilities
	mod := &WASMMod{
		Capabilities: CapReadEvents | CapReadGenres,
	}

	// Should have read caps
	if !mod.HasCapability(CapReadEvents) {
		t.Error("expected mod to have CapReadEvents")
	}
	if !mod.HasCapability(CapReadGenres) {
		t.Error("expected mod to have CapReadGenres")
	}

	// Should NOT have write caps
	if mod.HasCapability(CapWriteEvents) {
		t.Error("expected mod to NOT have CapWriteEvents")
	}
	if mod.HasCapability(CapWriteGenres) {
		t.Error("expected mod to NOT have CapWriteGenres")
	}
	if mod.HasCapability(CapWriteResources) {
		t.Error("expected mod to NOT have CapWriteResources")
	}
	if mod.HasCapability(CapModifyCrew) {
		t.Error("expected mod to NOT have CapModifyCrew")
	}
}

// TestHostDataConcurrency tests that host data is safe for concurrent access.
func TestHostDataConcurrency(t *testing.T) {
	hd := newHostData()

	done := make(chan bool, 2)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			hd.mu.Lock()
			hd.events = append(hd.events, EventDef{Title: "Test"})
			hd.mu.Unlock()
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			hd.mu.RLock()
			_ = len(hd.events)
			hd.mu.RUnlock()
		}
		done <- true
	}()

	<-done
	<-done
}

// TestWASMLoaderListEnabled tests listing enabled mods only.
func TestWASMLoaderListEnabledEmpty(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	enabled := loader.ListEnabled()
	if len(enabled) != 0 {
		t.Errorf("expected empty enabled list, got %d", len(enabled))
	}
}

// TestWASMModDisabled tests behavior when mod is disabled.
func TestWASMModDisabled(t *testing.T) {
	hd := &hostData{
		events: []EventDef{{Title: "Test"}},
	}
	mod := &WASMMod{
		ID:       "test",
		Name:     "Test Mod",
		Enabled:  false,
		hostData: hd,
	}

	// Disabled mod should still return its events
	events := mod.GetAddedEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

// TestWASMModGetAddedEventsWithData tests retrieving added events.
func TestWASMModGetAddedEventsWithData(t *testing.T) {
	hd := &hostData{
		events: []EventDef{
			{Title: "Event 1", Description: "Desc 1"},
			{Title: "Event 2", Description: "Desc 2"},
		},
	}
	mod := &WASMMod{
		hostData: hd,
	}

	events := mod.GetAddedEvents()
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
	if events[0].Title != "Event 1" {
		t.Errorf("expected 'Event 1', got %q", events[0].Title)
	}
}

// TestWASMModGetAddedGenresWithData tests retrieving added genres.
func TestWASMModGetAddedGenresWithData(t *testing.T) {
	hd := &hostData{
		genres: []GenreDef{
			{ID: "custom1", Name: "Custom Genre 1"},
			{ID: "custom2", Name: "Custom Genre 2"},
		},
	}
	mod := &WASMMod{
		hostData: hd,
	}

	genres := mod.GetAddedGenres()
	if len(genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(genres))
	}
	if genres[0].ID != "custom1" {
		t.Errorf("expected 'custom1', got %q", genres[0].ID)
	}
}

// TestWASMLoaderReloadMod tests that reloading a mod with same ID fails.
func TestWASMLoaderDuplicateModID(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	// Manually add a mod to simulate already loaded
	mod := &WASMMod{ID: "test-mod", Name: "Test"}
	loader.mu.Lock()
	loader.mods["test-mod"] = mod
	loader.modOrder = append(loader.modOrder, "test-mod")
	loader.mu.Unlock()

	// Verify we can get it
	retrieved, err := loader.Get("test-mod")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved.Name != "Test" {
		t.Error("retrieved wrong mod")
	}
}

// TestWASMLoaderUnload tests unloading a mod.
func TestWASMLoaderUnloadMod(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	// Manually add a mod
	mod := &WASMMod{ID: "test-mod", Name: "Test", hostData: newHostData()}
	loader.mu.Lock()
	loader.mods["test-mod"] = mod
	loader.modOrder = append(loader.modOrder, "test-mod")
	loader.mu.Unlock()

	// Unload it
	err := loader.Unload("test-mod")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should no longer exist
	_, err = loader.Get("test-mod")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

// TestNewHostData tests that newHostData initializes all fields.
func TestNewHostData(t *testing.T) {
	hd := newHostData()

	if hd.events == nil {
		t.Error("events should be initialized")
	}
	if hd.genres == nil {
		t.Error("genres should be initialized")
	}
	if hd.resources == nil {
		t.Error("resources should be initialized")
	}
	if hd.outputBuffer == nil {
		t.Error("outputBuffer should be initialized")
	}
	if hd.inputBuffer == nil {
		t.Error("inputBuffer should be initialized")
	}
}

// TestCapabilityStringAll tests the string representation of all capabilities.
func TestCapabilityStringAll(t *testing.T) {
	result := CapabilityString(AllCapabilities)
	expected := []string{
		"read_events", "write_events",
		"read_genres", "write_genres",
		"read_resources", "write_resources",
		"read_crew", "modify_crew",
		"trigger_events", "access_rng",
	}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("CapabilityString(AllCapabilities) missing %q", exp)
		}
	}
}

// TestWASMModMetadataFields tests that metadata fields can be accessed.
func TestWASMModMetadataFields(t *testing.T) {
	mod := &WASMMod{
		ID:          "test-id",
		Name:        "Test Name",
		Version:     "1.0.0",
		Author:      "Test Author",
		Description: "Test Description",
	}

	if mod.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %q", mod.ID)
	}
	if mod.Name != "Test Name" {
		t.Errorf("expected Name 'Test Name', got %q", mod.Name)
	}
	if mod.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got %q", mod.Version)
	}
	if mod.Author != "Test Author" {
		t.Errorf("expected Author 'Test Author', got %q", mod.Author)
	}
	if mod.Description != "Test Description" {
		t.Errorf("expected Description 'Test Description', got %q", mod.Description)
	}
}

// TestWASMLoaderCount tests mod counting.
func TestWASMLoaderCount(t *testing.T) {
	config := DefaultWASMConfig()
	loader := NewWASMLoader(config)

	if loader.Count() != 0 {
		t.Error("expected count 0")
	}

	// Manually add mods
	loader.mu.Lock()
	loader.mods["mod1"] = &WASMMod{ID: "mod1"}
	loader.mods["mod2"] = &WASMMod{ID: "mod2"}
	loader.modOrder = append(loader.modOrder, "mod1", "mod2")
	loader.mu.Unlock()

	if loader.Count() != 2 {
		t.Errorf("expected count 2, got %d", loader.Count())
	}
}
