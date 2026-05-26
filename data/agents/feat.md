---
description: Implémente une nouvelle fonctionnalité avec plan, validation et documentation.
mode: subagent
model: opencode/big-pickle
temperature: 0.4
permission:
  edit: ask
  bash: ask
---

## Usage
`@feat [description]`

## Protocole

### 1 — Clarification (max 3 questions)
Identifier les ambiguïtés. Comprendre le "quoi" et le "pourquoi".

### 2 — Analyse d'impact
- Impact sur architecture et design existant
- Fichiers modifiés, risques de régression
- Lire docs/DESIGN_SYSTEM.md avant tout ajout UI

### 3 — Plan
Générer `implementation_plan.md` :
```
# Plan : [Feature]
## User Stories
## Fichiers à créer/modifier
## Tests à écrire
## Impact design
## Dépendances/Risques
```

### 4 — Attendre validation utilisateur
Ne pas toucher au filesystem avant confirmation.

### 5 — Implémentation incrémentale
Commits atomiques par sous-fonctionnalité.

### 6 — Vérification
git diff → seuls fichiers prévus modifiés. Tests passants.

### 7 — Documentation
Mettre à jour docs/api.md, docs/features.md, task.md, walkthrough.md.
Commit : `[feat(scope)] - Ajout [fonctionnalité]`
