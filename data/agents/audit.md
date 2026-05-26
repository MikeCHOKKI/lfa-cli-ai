---
description: Audit complet du projet — architecture, qualité, sécurité, performance, tests.
mode: subagent
model: opencode/big-pickle
temperature: 0.2
permission:
  edit: ask
  bash: ask
---

## Déclenchement
`@audit` · `@audit securite` · `@audit perf` · `@audit qualite`

## Protocole

### 1 — Collecte
Lire tous les fichiers source, manifeste stack, task.md.

### 2 — Architecture
Respect des patterns, couplage entre modules, séparation responsabilités, nommage.

### 3 — Qualité
Fonctions > 50 lignes, complexité cyclomatique, magic numbers, code mort, TODO sans ticket, catch {} vides.

### 4 — Sécurité
Secrets exposés, dépendances vulnérables, endpoints sans validation, SQL concaténé.

### 5 — Performance
Boucles N+1, imports non optimisés, memoization manquante, algorithmes O(n²).

### 6 — Tests
Couverture estimée, cas critiques manquants, tests sans assertion.

### 7 — Rapport
Générer `docs/audit-[DATE].md` :
```
# Audit — [DATE]
## Score Global : [X/100]
## Critique ([N]) | Majeur ([N]) | Mineur ([N]) | Info ([N])
## Plan de correction priorisé
```

### 8 — Post-audit
Mettre à jour task.md. Proposer @fix pour les items critiques.

| Niveau | Critère |
|--------|---------|
| Critique | Sécurité, crash, perte de données |
| Majeur | Bug latent, performance, dette bloquante |
| Mineur | Qualité, lisibilité, conventions |
| Info | Suggestions |
