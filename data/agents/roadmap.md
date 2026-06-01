---
description: Génère la feuille de route du projet — phases, sprint actuel, backlog, dette technique, métriques.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@roadmap`

---

## Protocole

### 1 — Collecte
Lire dans l'ordre :
1. `task.md` → état actuel des tâches, backlog, blockers
2. `walkthrough.md` → historique des livraisons
3. `docs/features.md` → fonctionnalités prévues et leur statut
4. `PROJET.md` → objectifs, contraintes, vision produit
5. Tout `docs/audit-*.md` présent → dette technique connue

### 2 — Synthèse
- Identifier ce qui est **terminé** (livré, testé, documenté)
- Identifier ce qui est **en cours** (commencé mais non livré)
- Identifier ce qui est **planifié** (spécifié mais non commencé)
- Identifier la **dette technique** (issues audit, TODO non traités)
- Estimer les efforts restants par item (XS/S/M/L/XL)

### 3 — Génération
Créer ou écraser `docs/roadmap.md` :

```markdown
# Roadmap — [Nom du projet]
> Mise à jour : [DATE] | Horizon : [N semaines]

---

## Phase 1 — Fondations ✅ LIVRÉ
| Tâche | Commit | Date de livraison |
|-------|--------|-------------------|
| ...   | ...    | ...               |

---

## Phase 2 — Sprint Actuel 🔄 EN COURS
> Objectif : [ce que ce sprint doit livrer]

| Tâche | Priorité | Effort | Assigné |
|-------|----------|--------|---------|
| ...   | P0       | M      | —       |

---

## Phase 3 — Prochain Sprint 📋 PLANIFIÉ
| Tâche | Priorité | Effort |
|-------|----------|--------|

---

## Backlog
| Tâche | Origine | Effort estimé |
|-------|---------|---------------|

---

## Dette Technique
| Item | Réf Audit | Impact | Effort |
|------|-----------|--------|--------|

---

## Métriques
| Indicateur | Valeur actuelle | Cible |
|------------|-----------------|-------|
| Couverture tests | X% | 80% |
| Score sécurité | X/100 | 90/100 |
| Score performance | X/100 | 85/100 |
```

### 4 — Post-génération
- Mettre à jour `task.md` : marquer les tâches terminées, ajouter les nouvelles
- Proposer `@feat [tâche]` pour les items P0 du sprint actuel non commencés
