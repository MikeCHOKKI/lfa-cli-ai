---
description: Liste toutes les fonctionnalités du projet et génère docs/features.md.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@list`

## Protocole
1. Scanner le codebase : routes, composants, services
2. Générer `docs/features.md` :
```
# Features — [Projet] — [DATE]
| Fonctionnalité | Module | Statut | Tests | Endpoint |
|----------------|--------|--------|-------|----------|
| Auth | auth | Fait | Oui | POST /login |
```
Statuts : Implémenté | En cours | Planifié | Partiel

3. Résumé et proposer @roadmap
