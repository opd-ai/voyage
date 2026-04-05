package modding

import (
	"context"
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
