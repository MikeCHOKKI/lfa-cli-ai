---
description: Vérifie la conformité du projet aux standards (structure, code, commits, sécurité).
mode: subagent
model: opencode/big-pickle
temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@conformite`

## Checklist

### Structure
- [ ] AGENTS.md présent, dans .gitignore si règles globales
- [ ] PROJET.md, task.md, walkthrough.md présents
- [ ] docs/ avec architecture.md, api.md, DESIGN_SYSTEM.md
- [ ] .env.example commité, .env ignoré
- [ ] .gitignore couvre .env, .opencode/, .agents/, node_modules
- [ ] uml_engine défini dans PROJET.md

### Code
- [ ] Fonctions ≤ 50 lignes, pas de magic numbers
- [ ] Pas de catch {} vide, console.log en prod, TODO sans ticket
- [ ] Typage strict activé

### Commits
Lire les 10 derniers commits. Vérifier `[type(scope)] - Titre`.

## Rapport
```
Structure : [X/6] | Code : [X/3] | Commits : [X/10] | Sécurité : [X/3]
Score global : [XX]%
```
Score < 80% → proposer @fix conformite. Score ≥ 80% → noter dans task.md.
