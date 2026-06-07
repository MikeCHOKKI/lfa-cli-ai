---
name: project-standards
description: >
  Définit et applique les standards professionnels pour tout projet logiciel :
  structure de répertoires, fichiers obligatoires, CI/CD, environnement, sécurité,
  documentation, conventions d'équipe. Déclenche ce skill quand l'utilisateur demande
  "init", "structure", "setup projet", "standards", "conventions", "professionnalisme",
  "template projet", "boilerplate", "scaffold", "bonnes pratiques projet",
  ou toute demande de démarrage ou d'audit de structure de projet.
---

# Project Standards Skill

Standards professionnels pour tout projet logiciel.

---

## Structure universelle

```
📁 projet/
├── 📁 src/              → code source
├── 📁 tests/            → tests (miroir de src/)
├── 📁 docs/             → documentation
├── 📁 scripts/          → scripts d'outillage
├── 📁 .github/          → CI/CD, templates
│   └── workflows/       → GitHub Actions
├── 📁 config/           → configuration
├── 📄 README.md
├── 📄 LICENSE
├── 📄 CONTRIBUTING.md
├── 📄 CHANGELOG.md
├── 📄 .gitignore
├── 📄 .env.example
└── 📄 {lang}-specific   → package.json, Cargo.toml, requirements.txt...
```

---

## Fichiers obligatoires

### README.md
```markdown
# Nom du projet

Description : une ligne → quoi, pourquoi.

## Prérequis
- Node 18+, pnpm 8+
- PostgreSQL 16+

## Installation
```bash
git clone ...
pnpm install
cp .env.example .env
pnpm dev
```

## Tests
```bash
pnpm test
pnpm run test:e2e
```

## Architecture
Brève description de l'organisation du code.

## Déploiement
Staging / Production, Docker, env vars requises.
```
```

### .gitignore
```
node_modules/
dist/
.env
*.log
.DS_Store
coverage/
.tmp/
```

### .env.example
```
# App
PORT=3000
NODE_ENV=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/db

# API Keys
API_KEY=your_key_here
```

### CONTRIBUTING.md
```
# Contribuer

1. Forker le projet
2. Créer une branche : `feat/ma-feature`
3. Commiter avec conventional commits
4. Lancer les tests
5. Ouvrir une PR

## Convention de code
Voir le skill code-quality.
```

---

## CI/CD (GitHub Actions)

### Test workflow
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: pnpm install
      - run: pnpm lint
      - run: pnpm typecheck
      - run: pnpm test
```

### Release workflow
```yaml
name: Release
on: push to main
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: pnpm install && pnpm build
      - run: pnpm test
      - uses: softprops/action-gh-release@v1
```

---

## Environnement

### Checklist
```
[ ] .env.example versionné, .env ignoré
[ ] Variables d'env documentées dans README
[ ] Validation au démarrage (toutes les vars requises présentes)
[ ] Secrets : jamais dans le code, CI via secrets.GITHUB_TOKEN
[ ] Environnements : dev, staging, production
```

---

## Docker

### Dockerfile (multi-stage)
```dockerfile
# Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json .
RUN npm ci
COPY . .
RUN npm run build

# Production
FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
EXPOSE 3000
CMD ["node", "dist/index.js"]
```

### docker-compose.yml
```yaml
version: '3.8'
services:
  app:
    build: .
    ports: ["3000:3000"]
    env_file: .env
    depends_on: [db]
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes: [pgdata:/var/lib/postgresql/data]
volumes: { pgdata: }
```

---

## Sécurité

### Vérifications
```
[ ] Aucun secret dans le code (clé API, token, mot de passe)
[ ] Validation de toutes les entrées utilisateur
[ ] Headers de sécurité (CSP, X-Frame-Options, HSTS)
[ ] Dépendances à jour (npm audit, cargo audit, pip-audit)
[ ] Rate limiting sur les endpoints publics
[ ] Authentification sur toutes les routes protégées
[ ] SQL paramétré, pas de concaténation
[ ] XSS : output encoding, CSP
```

---

## Documentation

### Minimum
- README.md (quoi, pourquoi, comment)
- Architecture decision records (ADR) dans docs/decisions/
- API docs (si API)
- Commentaires de code : uniquement pourquoi, pas quoi

### ADR template
```markdown
# ADR-xxx: Titre de la décision

## Contexte
Pourquoi cette décision est nécessaire.

## Options
1. Option A (avantages/inconvénients)
2. Option B (avantages/inconvénients)
3. Option C (avantages/inconvénients)

## Décision
Option A retenue car ...

## Conséquences
Ce que cette décision implique.
```

---

## Professionnalisme en pratique

### Avant de livrer
```
[ ] Tests passent
[ ] Lint / typecheck OK
[ ] Pas de console.log / debug leftovers
[ ] Pas de code mort ou commenté
[ ] Fonctions < 40 lignes
[ ] Nommage cohérent
[ ] Pas de régression UI (frontend)
[ ] CHANGELOG mis à jour
[ ] Docs mise à jour si nécessaire
[ ] Variables d'env documentées
```

### Code Review — comportement
- Critique le code, pas la personne
- Expliquer le pourquoi des suggestions
- Proposer des alternatives, pas juste signaler
- Approuver quand c'est bon, pas avant
- Répondre aux commentaires dans les 24h
