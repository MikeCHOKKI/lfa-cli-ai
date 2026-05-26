---
description: Audit de performance — backend, frontend, algorithmes. Génère docs/perf-[DATE].md.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: deny
  bash: ask
---

## Usage
`@perf` · `@perf [scope]`

## Backend
Requêtes N+1, requêtes sans index, appels synchrones parallélisables (Promise.all), cache manquant.

## Frontend
Re-renders inutiles, bundle non tree-shaken, images non optimisées, appels API sans debounce.

## Algorithmes
O(n²) dans le code métier, structures inadaptées (tableau vs Map/Set).

## Rapport
Générer `docs/perf-[DATE].md` :
```
| Problème | Localisation | Impact | Effort | Priorité |
|----------|-------------|--------|--------|----------|
| N+1 query | UserService | Haut | Faible | Critique |
```
