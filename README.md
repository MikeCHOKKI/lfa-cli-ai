# LFA CLI — OpenCode AI Configuration Tool

[![CI](https://github.com/MikeCHOKKI/lfa-cli-ai/actions/workflows/ci.yml/badge.svg)](https://github.com/MikeCHOKKI/lfa-cli-ai/actions/workflows/ci.yml)
[![Release](https://github.com/MikeCHOKKI/lfa-cli-ai/actions/workflows/release.yml/badge.svg)](https://github.com/MikeCHOKKI/lfa-cli-ai/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/lfa-cli/lfa-cli-ai)](https://goreportcard.com/report/github.com/lfa-cli/lfa-cli-ai)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

CLI Go qui automatise la détection, l'installation et la configuration d'OpenCode et de ses agents/skills IA.

## Aperçu

LFA CLI analyse votre système, détecte les outils IA installés (OpenCode, Ollama), puis déploie automatiquement :
- **22 agents** pré-configurés
- **18 skills** spécialisés
- **Configuration MCP** (filesystem, memory, github, fetch, ollama)
- **Permissions** OpenCode adaptées

## Installation

```bash
curl -fsSL https://lfa-cli.vercel.app/install.sh | sh
```

Ou téléchargez le binaire depuis les [GitHub Releases](https://github.com/MikeCHOKKI/lfa-cli-ai/releases).

## Commandes

| Commande | Description |
|----------|-------------|
| `lfa version` | Affiche la version |
| `lfa doctor` | Diagnostic système (OS, OpenCode, Ollama, chemins) |
| `lfa setup` | Déploie la configuration OpenCode |
| `lfa dashboard [-y]` | TUI interactif (ou `-y` pour mode non-interactif) |

### Flags

| Commande | Flag | Description |
|----------|------|-------------|
| `setup` | `--ollama` | Active l'intégration Ollama (par défaut: true) |
| `setup` | `--dry-run` | Simule sans écrire de fichiers |
| *(global)* | `-y, --yes` | Mode non-interactif (répond oui à tout) |

### Exemples

```bash
# Diagnostic rapide
lfa doctor

# Installation complète avec Ollama
lfa setup --ollama

# Simulation sans écriture
lfa setup --dry-run

# Mode automatique
lfa dashboard -y
```

## Architecture

```
lfa-cli-ai/
├── main.go                     # Point d'entrée
├── cmd/
│   ├── root.go                 # Commande racine (Cobra)
│   ├── version.go              # lfa version
│   ├── doctor.go               # lfa doctor
│   ├── setup.go                # lfa setup
│   ├── dashboard.go            # lfa dashboard
│   └── setup_shared.go         # Logique partagée setup/dashboard
├── internal/
│   ├── config/
│   │   └── config.go           # Génération/écriture config OpenCode
│   ├── detect/
│   │   └── detect.go           # Détection OS, OpenCode, Ollama
│   ├── installer/
│   │   └── installer.go        # Téléchargement, extraction, déploiement
│   └── ui/
│       ├── tui.go              # TUI Bubbletea (dashboard interactif)
│       ├── styles.go           # Thème Lipgloss
│       └── prompts.go          # Prompts Huh (confirm, select)
└── data/
    ├── agents/                 # 22 fichiers .md (agents IA)
    ├── skills/                 # 18 dossiers (skills spécialisés)
    ├── opencode.jsonc          # Template de configuration
    └── AGENTS.md               # Règles globales
```

## Stack technique

| Composant | Bibliothèque |
|-----------|-------------|
| CLI framework | [Cobra](https://github.com/spf13/cobra) v1.8.1 |
| TUI | [Bubbletea](https://github.com/charmbracelet/bubbletea) v1.1.0 |
| Styles terminal | [Lipgloss](https://github.com/charmbracelet/lipgloss) v0.13.0 |
| Formulaires | [Huh](https://github.com/charmbracelet/huh) v0.6.0 |
| Logging | [charmbracelet/log](https://github.com/charmbracelet/log) v0.4.0 |
| Widgets TUI | [Bubbles](https://github.com/charmbracelet/bubbles) v0.20.0 |

## Build

```bash
make build      # Build pour la plateforme courante
make test       # Exécute les tests
make lint       # go vet + staticcheck
make coverage   # Couverture de tests (HTML)
make clean      # Supprime le binaire
make run        # Build et exécute
```

### Cross-compilation

```bash
make cross      # Build les 6 cibles
```

| OS | Architecture |
|----|-------------|
| linux | amd64, arm64 |
| darwin | amd64, arm64 |
| windows | amd64, arm64 |

Le binaire est strippé (`-ldflags "-s -w"`) pour une taille de ~9.6 MB.

## CI/CD

Trois workflows GitHub Actions :

| Workflow | Déclencheur | Action |
|----------|------------|--------|
| `ci.yml` | push/PR sur `main` | lint, test, build, cross-compilation |
| `release.yml` | tag `v*` | build multiplateforme + checksums SHA256 + GitHub Release |
| `release-please.yml` | push sur `main` | auto-tag et changelog automatique |

## Configuration déployée

Après `lfa setup`, les fichiers sont installés dans :

| Platform | Chemin |
|----------|--------|
| Linux | `~/.config/opencode/` |
| macOS | `~/Library/Application Support/opencode/` |
| Windows | `%APPDATA%\opencode\` |

Structure déployée :
```
~/.config/opencode/
├── opencode.jsonc        # Configuration principale
├── AGENTS.md             # Règles globales
├── agents/               # 22 agents .md
└── skills/               # 18 skills complets
```

## Variables d'environnement

| Variable | Description | Requis |
|----------|-------------|--------|
| `GITHUB_TOKEN` | Token pour le MCP GitHub | Optionnel (requis pour les fonctionnalités GitHub) |

## Sécurité

- Téléchargements via HTTPS uniquement
- `io.LimitReader` (500 MB) contre les decompression bombs
- Versions pinnées pour les packages MCP (npx @0.6.2, uvx @0.1.4)
- Aucun secret hardcodé — tokens via variables d'environnement
- `.env` exclu du version control

Voir [SECURITY.md](SECURITY.md) pour plus de détails.

## Contribuer

Voir [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — voir [LICENSE](LICENSE).
