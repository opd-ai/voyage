# Contributing to Voyage

Thank you for your interest in contributing to Voyage! This document provides guidelines for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git
- OpenGL development libraries (Linux)

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/voyage.git
   cd voyage
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Verify the build:
   ```bash
   go build ./...
   go test ./...
   ```

### Linux Dependencies

On Ubuntu/Debian:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

On Fedora:
```bash
sudo dnf install mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel
```

## Code Style

### Go Standards

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (run automatically with most editors)
- Run `go vet` before committing

### Project Conventions

- **Documentation**: All exported symbols must have doc comments
- **Error Handling**: Return errors rather than panicking; wrap with context
- **Testing**: Aim for ≥40% coverage per package
- **Complexity**: Keep functions under 30 lines and cyclomatic complexity under 10

### GenreSwitcher Interface

All ECS Systems must implement the `GenreSwitcher` interface:

```go
type GenreSwitcher interface {
    SetGenre(genreID GenreID)
}
```

Use the `BaseSystem` embed for default implementation:

```go
type MySystem struct {
    engine.BaseSystem
    // ...
}
```

### No Bundled Assets

This project enforces **100% procedural generation**. Do not add:
- Image files (`.png`, `.jpg`, `.gif`, `.svg`, etc.)
- Audio files (`.mp3`, `.wav`, `.ogg`, etc.)
- Pre-written narrative content

All visual and audio content must be generated at runtime.

## Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following the code style guidelines

3. Write or update tests as needed

4. Ensure all checks pass:
   ```bash
   go build ./...
   go test -race ./...
   go vet ./...
   ./scripts/validate-no-assets.sh
   ```

5. Commit with a clear message:
   ```bash
   git commit -m "Add feature: brief description"
   ```

6. Push and create a pull request on GitHub

7. Wait for CI to pass and address any review feedback

## Testing Guidelines

### Unit Tests

- Place tests in `*_test.go` files alongside the code
- Use table-driven tests for multiple cases
- Test both success and error paths

### Determinism Tests

For procedural generation, verify determinism:

```go
func TestDeterminism(t *testing.T) {
    seed := int64(12345)
    result1 := Generate(seed)
    result2 := Generate(seed)
    if !reflect.DeepEqual(result1, result2) {
        t.Error("Same seed should produce same result")
    }
}
```

### Benchmarks

Add benchmarks for performance-critical code:

```go
func BenchmarkGenerate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Generate(int64(i))
    }
}
```

Run with: `go test -bench=. ./pkg/...`

## Reporting Issues

When reporting bugs, please include:

- Go version (`go version`)
- Operating system and version
- Steps to reproduce
- Expected vs actual behavior
- Seed value (if relevant to reproduction)

## Feature Requests

Check the [ROADMAP.md](ROADMAP.md) for planned features. If your idea isn't listed:

1. Open an issue describing the feature
2. Explain how it fits the project's procedural generation philosophy
3. Discuss implementation approach if you plan to contribute

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
