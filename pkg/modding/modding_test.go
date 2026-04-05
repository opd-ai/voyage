package modding

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestModValidate(t *testing.T) {
	tests := []struct {
		name    string
		mod     Mod
		wantErr bool
	}{
		{
			name:    "empty mod",
			mod:     Mod{},
			wantErr: true,
		},
		{
			name:    "missing name",
			mod:     Mod{ID: "test"},
			wantErr: true,
		},
		{
			name:    "missing version",
			mod:     Mod{ID: "test", Name: "Test"},
			wantErr: true,
		},
		{
			name:    "valid minimal mod",
			mod:     Mod{ID: "test", Name: "Test", Version: "1.0.0"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mod.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventDefValidate(t *testing.T) {
	tests := []struct {
		name    string
		event   EventDef
		wantErr bool
	}{
		{
			name:    "empty event",
			event:   EventDef{},
			wantErr: true,
		},
		{
			name:    "missing description",
			event:   EventDef{Title: "Test", Category: "encounter"},
			wantErr: true,
		},
		{
			name:    "missing choices",
			event:   EventDef{Title: "Test", Description: "Desc", Category: "encounter"},
			wantErr: true,
		},
		{
			name:    "invalid category",
			event:   EventDef{Title: "Test", Description: "Desc", Category: "invalid", Choices: []ChoiceDef{{Text: "ok"}}},
			wantErr: true,
		},
		{
			name: "valid event",
			event: EventDef{
				Title:       "Test Event",
				Description: "Something happens",
				Category:    "encounter",
				Choices:     []ChoiceDef{{Text: "Do something"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenreDefValidate(t *testing.T) {
	tests := []struct {
		name    string
		genre   GenreDef
		wantErr bool
	}{
		{
			name:    "empty genre",
			genre:   GenreDef{},
			wantErr: true,
		},
		{
			name:    "missing name",
			genre:   GenreDef{ID: "test"},
			wantErr: true,
		},
		{
			name:    "valid genre",
			genre:   GenreDef{ID: "steampunk", Name: "Steampunk"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genre.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoaderRegister(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{ID: "test", Name: "Test Mod", Version: "1.0.0"}
	registered, err := loader.Register(mod)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if registered.ID != mod.ID {
		t.Errorf("expected ID %s, got %s", mod.ID, registered.ID)
	}
	if !registered.Enabled {
		t.Error("mod should be enabled by default")
	}
	if registered.LoadedAt.IsZero() {
		t.Error("LoadedAt should be set")
	}

	// Try to register again
	_, err = loader.Register(mod)
	if err != ErrModAlreadyLoaded {
		t.Errorf("expected ErrModAlreadyLoaded, got %v", err)
	}
}

func TestLoaderGet(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{ID: "test", Name: "Test", Version: "1.0.0"}
	_, _ = loader.Register(mod)

	got, err := loader.Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != "test" {
		t.Errorf("expected ID test, got %s", got.ID)
	}

	_, err = loader.Get("nonexistent")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

func TestLoaderUnload(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{ID: "test", Name: "Test", Version: "1.0.0"}
	_, _ = loader.Register(mod)

	err := loader.Unload("test")
	if err != nil {
		t.Fatalf("Unload() error = %v", err)
	}

	_, err = loader.Get("test")
	if err != ErrModNotFound {
		t.Error("mod should be removed after unload")
	}

	err = loader.Unload("nonexistent")
	if err != ErrModNotFound {
		t.Errorf("expected ErrModNotFound, got %v", err)
	}
}

func TestLoaderList(t *testing.T) {
	loader := NewLoader()

	mod1 := &Mod{ID: "mod1", Name: "Mod 1", Version: "1.0.0"}
	mod2 := &Mod{ID: "mod2", Name: "Mod 2", Version: "1.0.0"}
	_, _ = loader.Register(mod1)
	_, _ = loader.Register(mod2)

	list := loader.List()
	if len(list) != 2 {
		t.Errorf("expected 2 mods, got %d", len(list))
	}

	// Check order
	if list[0].ID != "mod1" || list[1].ID != "mod2" {
		t.Error("mods should be in load order")
	}
}

func TestLoaderEnableDisable(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{ID: "test", Name: "Test", Version: "1.0.0"}
	_, _ = loader.Register(mod)

	err := loader.Disable("test")
	if err != nil {
		t.Fatalf("Disable() error = %v", err)
	}

	enabled := loader.ListEnabled()
	if len(enabled) != 0 {
		t.Error("disabled mod should not be in enabled list")
	}

	err = loader.Enable("test")
	if err != nil {
		t.Fatalf("Enable() error = %v", err)
	}

	enabled = loader.ListEnabled()
	if len(enabled) != 1 {
		t.Error("enabled mod should be in enabled list")
	}
}

func TestLoaderLoadFromBytes(t *testing.T) {
	loader := NewLoader()

	jsonData := []byte(`{
		"id": "test-mod",
		"name": "Test Mod",
		"version": "1.0.0",
		"author": "TestAuthor",
		"events": [
			{
				"category": "encounter",
				"genre": "fantasy",
				"title": "Test Event",
				"description": "A test event",
				"choices": [
					{"text": "Option 1", "outcome": {"morale_delta": 5}}
				]
			}
		]
	}`)

	mod, err := loader.LoadFromBytes(jsonData)
	if err != nil {
		t.Fatalf("LoadFromBytes() error = %v", err)
	}

	if mod.ID != "test-mod" {
		t.Errorf("expected ID test-mod, got %s", mod.ID)
	}
	if len(mod.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(mod.Events))
	}
}

func TestLoaderLoadFromReader(t *testing.T) {
	loader := NewLoader()

	jsonData := `{
		"id": "reader-mod",
		"name": "Reader Mod",
		"version": "1.0.0"
	}`

	reader := bytes.NewReader([]byte(jsonData))
	mod, err := loader.LoadFromReader(reader)
	if err != nil {
		t.Fatalf("LoadFromReader() error = %v", err)
	}

	if mod.ID != "reader-mod" {
		t.Errorf("expected ID reader-mod, got %s", mod.ID)
	}
}

func TestLoaderLoadFromBytesInvalidJSON(t *testing.T) {
	loader := NewLoader()

	_, err := loader.LoadFromBytes([]byte(`{invalid json`))
	if err != ErrInvalidJSON {
		t.Errorf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestLoaderLoadFromBytesValidationFails(t *testing.T) {
	loader := NewLoader()

	// Missing required fields
	_, err := loader.LoadFromBytes([]byte(`{"id": "test"}`))
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestLoaderGetAllEvents(t *testing.T) {
	loader := NewLoader()

	mod1 := &Mod{
		ID:      "mod1",
		Name:    "Mod 1",
		Version: "1.0.0",
		Events: []EventDef{
			{Title: "Event 1", Description: "Desc", Category: "encounter", Choices: []ChoiceDef{{Text: "ok"}}},
		},
	}
	mod2 := &Mod{
		ID:      "mod2",
		Name:    "Mod 2",
		Version: "1.0.0",
		Events: []EventDef{
			{Title: "Event 2", Description: "Desc", Category: "discovery", Choices: []ChoiceDef{{Text: "ok"}}},
			{Title: "Event 3", Description: "Desc", Category: "hazard", Choices: []ChoiceDef{{Text: "ok"}}},
		},
	}

	_, _ = loader.Register(mod1)
	_, _ = loader.Register(mod2)

	events := loader.GetAllEvents()
	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}

func TestLoaderGetEventsForGenre(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{
		ID:      "test",
		Name:    "Test",
		Version: "1.0.0",
		Events: []EventDef{
			{Title: "Fantasy Event", Description: "Desc", Category: "encounter", Genre: "fantasy", Choices: []ChoiceDef{{Text: "ok"}}},
			{Title: "Scifi Event", Description: "Desc", Category: "encounter", Genre: "scifi", Choices: []ChoiceDef{{Text: "ok"}}},
			{Title: "Generic Event", Description: "Desc", Category: "encounter", Genre: "", Choices: []ChoiceDef{{Text: "ok"}}},
		},
	}
	_, _ = loader.Register(mod)

	fantasyEvents := loader.GetEventsForGenre("fantasy")
	if len(fantasyEvents) != 2 { // fantasy + generic
		t.Errorf("expected 2 fantasy events, got %d", len(fantasyEvents))
	}

	scifiEvents := loader.GetEventsForGenre("scifi")
	if len(scifiEvents) != 2 { // scifi + generic
		t.Errorf("expected 2 scifi events, got %d", len(scifiEvents))
	}
}

func TestLoaderGetCustomGenres(t *testing.T) {
	loader := NewLoader()

	mod := &Mod{
		ID:      "test",
		Name:    "Test",
		Version: "1.0.0",
		Genres: []GenreDef{
			{ID: "steampunk", Name: "Steampunk"},
			{ID: "western", Name: "Wild West"},
		},
	}
	_, _ = loader.Register(mod)

	genres := loader.GetCustomGenres()
	if len(genres) != 2 {
		t.Errorf("expected 2 genres, got %d", len(genres))
	}
}

func TestLoaderCount(t *testing.T) {
	loader := NewLoader()

	if loader.Count() != 0 {
		t.Error("expected 0 mods initially")
	}

	_, _ = loader.Register(&Mod{ID: "1", Name: "1", Version: "1.0.0"})
	_, _ = loader.Register(&Mod{ID: "2", Name: "2", Version: "1.0.0"})

	if loader.Count() != 2 {
		t.Errorf("expected 2 mods, got %d", loader.Count())
	}
}

func TestLoaderClear(t *testing.T) {
	loader := NewLoader()

	_, _ = loader.Register(&Mod{ID: "1", Name: "1", Version: "1.0.0"})
	_, _ = loader.Register(&Mod{ID: "2", Name: "2", Version: "1.0.0"})

	loader.Clear()

	if loader.Count() != 0 {
		t.Error("expected 0 mods after clear")
	}
}

func TestLoaderLoadFromFile(t *testing.T) {
	loader := NewLoader()

	// Create temp file
	tmpDir := t.TempDir()
	modPath := filepath.Join(tmpDir, "test-mod.json")

	modJSON := `{
		"id": "file-mod",
		"name": "File Mod",
		"version": "1.0.0"
	}`

	err := os.WriteFile(modPath, []byte(modJSON), 0o644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	mod, err := loader.LoadFromFile(modPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if mod.ID != "file-mod" {
		t.Errorf("expected ID file-mod, got %s", mod.ID)
	}
}

func TestLoaderLoadFromFileNotFound(t *testing.T) {
	loader := NewLoader()

	_, err := loader.LoadFromFile("/nonexistent/path/mod.json")
	if err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}

func TestLoaderLoadDirectory(t *testing.T) {
	loader := NewLoader()

	tmpDir := t.TempDir()

	// Create some mod files
	mod1JSON := `{"id": "mod1", "name": "Mod 1", "version": "1.0.0"}`
	mod2JSON := `{"id": "mod2", "name": "Mod 2", "version": "1.0.0"}`
	invalidJSON := `{invalid`

	_ = os.WriteFile(filepath.Join(tmpDir, "mod1.json"), []byte(mod1JSON), 0o644)
	_ = os.WriteFile(filepath.Join(tmpDir, "mod2.json"), []byte(mod2JSON), 0o644)
	_ = os.WriteFile(filepath.Join(tmpDir, "invalid.json"), []byte(invalidJSON), 0o644)
	_ = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not a mod"), 0o644)

	mods, err := loader.LoadDirectory(tmpDir)
	if err != nil {
		t.Fatalf("LoadDirectory() error = %v", err)
	}

	// Should load 2 valid mods, skip invalid and non-json
	if len(mods) != 2 {
		t.Errorf("expected 2 mods, got %d", len(mods))
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "events",
		Index:   0,
		Message: "missing title",
	}

	expected := "events[0]: missing title"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}

	// Without index
	err2 := &ValidationError{
		Field:   "name",
		Index:   -1,
		Message: "required",
	}

	expected2 := "name: required"
	if err2.Error() != expected2 {
		t.Errorf("expected %q, got %q", expected2, err2.Error())
	}
}

func TestValidCategories(t *testing.T) {
	cats := ValidCategories()
	if len(cats) != 7 {
		t.Errorf("expected 7 categories, got %d", len(cats))
	}
}

func TestValidGenres(t *testing.T) {
	genres := ValidGenres()
	if len(genres) != 5 {
		t.Errorf("expected 5 genres, got %d", len(genres))
	}
}

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		cat   string
		valid bool
	}{
		{"weather", true},
		{"encounter", true},
		{"discovery", true},
		{"hardship", true},
		{"windfall", true},
		{"hazard", true},
		{"crew", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		if isValidCategory(tt.cat) != tt.valid {
			t.Errorf("isValidCategory(%q) = %v, want %v", tt.cat, !tt.valid, tt.valid)
		}
	}
}
