package modding

import "errors"

// Sentinel errors for mod validation.
var (
	ErrMissingID          = errors.New("mod ID is required")
	ErrMissingName        = errors.New("mod name is required")
	ErrMissingVersion     = errors.New("mod version is required")
	ErrMissingTitle       = errors.New("event title is required")
	ErrMissingDescription = errors.New("event description is required")
	ErrMissingCategory    = errors.New("event category is required")
	ErrInvalidCategory    = errors.New("invalid event category")
	ErrNoChoices          = errors.New("event must have at least one choice")
	ErrModNotFound        = errors.New("mod not found")
	ErrModAlreadyLoaded   = errors.New("mod already loaded")
	ErrInvalidJSON        = errors.New("invalid JSON format")
	ErrFileNotFound       = errors.New("mod file not found")
	ErrCircularDependency = errors.New("circular mod dependency detected")
)

// ValidationError provides details about a validation failure.
type ValidationError struct {
	Field   string
	Index   int
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Index >= 0 {
		return e.Field + "[" + itoa(e.Index) + "]: " + e.Message
	}
	return e.Field + ": " + e.Message
}

// itoa converts an int to string without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	result := make([]byte, 0, 12)
	negative := i < 0
	if negative {
		i = -i
	}

	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}

	if negative {
		result = append([]byte{'-'}, result...)
	}

	return string(result)
}
