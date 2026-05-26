---
description: Analyse une question d'architecture et propose 3 approches avec recommandation.
mode: subagent
model: opencode/big-pickle
temperature: 0.7
permission:
  edit: deny
  bash: deny
---

## Usage
`@question [sujet ou dilemme]`

## Protocole

### 1 — Contextualisation
Lire docs/architecture.md + codebase. Comprendre les contraintes existantes.

### 2 — 3 Approches
Pour chaque : description, avantages, inconvénients, complexité, compatibilité stack.

### 3 — Recommandation
Basée sur le projet réel. Justification concise.

### 4 — Décision
Si validée → créer `docs/decisions/ADR-[N]-[sujet].md` :
```
# ADR-[N] : [Titre]
Date : [DATE] | Statut : Accepté
Contexte / Décision / Conséquences / Alternatives rejetées
```
Commit : `[docs(decisions)] - ADR-[N] [sujet]`
