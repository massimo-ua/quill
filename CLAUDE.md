# Quill Project Development Guide

## Build & Test Commands
- Build: `go build ./...`
- Test all: `go test ./...`
- Test specific package: `go test github.com/massimo-ua/quill/internal/domain`
- Test specific test: `go test -run TestID_Time ./internal/domain/common`
- Coverage: `go test -cover ./...`
- Format code: `go fmt ./...`
- Lint: `go vet ./...`

## Code Style Guidelines
- **Naming**: Use CamelCase for types/funcs; avoid abbreviations
- **Imports**: Group standard, external, internal packages in this order
- **Error Handling**: Always check errors; return early on errors
- **Types**: Use meaningful types; prefer strong typing over primitive types
- **Testing**: Use table-driven tests with descriptive names
- **Documentation**: Add comments for exported types and functions
- **Package Structure**: Follow Domain-Driven Design (DDD) principles
- **File Structure**: One entity/value object per file with its tests