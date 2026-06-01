# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0](https://github.com/MikeCHOKKI/lfa-cli-ai/compare/v0.1.0...v0.2.0) (2026-06-01)


### Features

* add project files — LICENSE, CHANGELOG, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, .editorconfig, .golangci.yml ([84488c3](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/84488c34b6f411b6f832d21f1340863a9e5704fe))
* **data:** Sync agents et skills depuis la config opencode ([65e5a87](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/65e5a874536a182466a7a01d878d066ade833315))
* **release:** Génération automatique du changelog dans les notes de release ([69df000](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/69df0006698dd1c913763a7611932c0c21071d2d))


### Bug Fixes

* **ci:** lint + test sur develop, build + cross sur main/PR ([c91e590](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/c91e590bdee83520137d507578d62f2a3d237f53))
* **ci:** Remplacement de GITHUB_TOKEN par RELEASE_PLEASE_TOKEN ([f56749b](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/f56749bca1e5b72504732c515b664876f16ae5b5))
* **ci:** Remplacement de GITHUB_TOKEN par RELEASE_PLEASE_TOKEN ([9c035c4](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/9c035c4483e3522a607604635af2a9401111d5d9))
* **ci:** Restreindre les builds à la branche main uniquement ([95b773c](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/95b773c7c2ff71bf4cad70cd052aba84fc8d5fec))
* **ci:** Token explicite pour release-please ([23d868b](https://github.com/MikeCHOKKI/lfa-cli-ai/commit/23d868b83842c59db213fffd62775d29a7d923d6))

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
