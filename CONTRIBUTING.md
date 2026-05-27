# Contributing to LFA CLI

Thank you for your interest in contributing to LFA CLI!

## Getting Started

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Run tests and linting
5. Commit with [Conventional Commits](https://www.conventionalcommits.org/)
6. Push and open a Pull Request

## Development

### Prerequisites

- Go 1.22+
- Node.js 18+ (for UI)

### Build

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run go vet + staticcheck
make clean      # Remove binary
```

### Cross-compilation

```bash
make cross      # Build for all 6 targets
```

## Code Standards

### Go

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `fmt.Errorf("context: %w", err)` for error wrapping
- Write table-driven tests
- Run `go vet` and `staticcheck` before committing

### Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add Ollama detection caching
fix: prevent tar path traversal
docs: update README with new commands
chore: remove dead code
test: add unit tests for config package
```

## Pull Request Process

1. Update `CHANGELOG.md` with your changes
2. Ensure CI passes (lint, test, build)
3. Request review from a maintainer
4. Address review feedback

## Reporting Issues

- **Bug reports**: include OS, Go version, steps to reproduce
- **Feature requests**: describe the use case and expected behavior
- **Security issues**: see [SECURITY.md](SECURITY.md)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
