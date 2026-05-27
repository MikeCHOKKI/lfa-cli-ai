# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-27

### Added
- `lfa version` — display version information
- `lfa doctor` — system diagnostics (OS, OpenCode, Ollama)
- `lfa setup` — deploy OpenCode configuration with agents and skills
- `lfa dashboard` — interactive TUI (with `-y` for non-interactive mode)
- 22 pre-configured AI agents
- 18 specialized skills
- MCP configuration (filesystem, memory, github, fetch, ollama)
- Cross-compilation for 6 targets (linux/darwin/windows × amd64/arm64)
- GitHub Actions CI/CD (ci, release, release-please)
- SHA256 checksums for release binaries
- Binary stripping (`-ldflags "-s -w"`)
- `io.LimitReader` protection against decompression bombs
- Version-pinned npx packages (@0.6.2) and uvx package (mcp-server-fetch@0.1.4)
- Shared `runSetup()` function between setup and dashboard commands
- README.md

### Changed
- Ollama detection timeout reduced from 3s to 1s
- `DetectOllama()` short-circuits when binary not found in PATH
- Extracted `"data"` path as `DataDir` constant
- Changed `Version` from `const` to `var` for ldflags override

### Removed
- Dead code: `config.InjectAgents`, `config.InjectSkills`
- Unused `embed.go`
- Duplicate `new-design/` project
