---
description: Génère des diagrammes d'activité UML pour modéliser les processus métier et workflows.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@activity [processus]` — modélise un processus métier en diagramme d'activité.

## Protocole

### 1 — Analyse du processus
Lire les fichiers pertinents. Identifier acteurs, activités, décisions, flux parallèles.

### 2 — Conception du diagramme
- Nœuds initial et final
- Activités et actions
- Flots de contrôle
- Nœuds de décision (branchement)
- Fourches et jonctions (parallélisme)
- Couloirs (swimlanes) si plusieurs acteurs

### 3 — Génération
Mermaid → `docs/uml/activity-[module]-[DATE].mmd`
PlantUML → `docs/uml/activity-[module]-[DATE].puml`

### 4 — Validation
Vérifier complétude, logique métier, cohérence avec la documentation existante.

### 5 — Documentation
Commit : `[uml(docs)] - Diagramme d'activité [module]`
