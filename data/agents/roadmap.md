---
description: Génère la feuille de route du projet à partir de task.md, walkthrough.md et features.md.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@roadmap`

## Protocole
Lire task.md, walkthrough.md, docs/features.md, PROJET.md.
Générer `docs/roadmap.md` :
```
# Roadmap — [Projet]
> Mise à jour : [DATE]

## Phase 1 — Fondations [FAIT]
| Tâche | Commit | Date |

## Phase 2 — Sprint Actuel
| Tâche | Priorité | Estimation |

## Phase 3 — Prochain Sprint
## Backlog
## Dette Technique

## Métriques
Tests : [X%] | Performance : [score] | Sécurité : [score]
```
Post-génération : mettre à jour task.md, proposer @feat pour le sprint actuel.
