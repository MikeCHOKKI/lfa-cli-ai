---
description: Audit de performance — backend, frontend, algorithmes. Localise, quantifie et priorise les problèmes.
mode: subagent

temperature: 0.3
permission:
  edit: deny
  bash: ask
---

## Usage
`@perf` · `@perf [scope]`

---

## Protocole

### 1 — Collecte
- Lire `PROJET.md` → stack (ORM, framework, renderer)
- Si scope précisé → restreindre l'analyse au module/couche concerné
- Lire les fichiers source dans l'ordre : services → repositories → contrôleurs → composants UI

### 2 — Analyse Backend

**Requêtes N+1**
- Boucles contenant des appels à la BDD ou à l'ORM
- Relations chargées sans `eager loading` / `JOIN` / `includes`
- Pattern : `foreach($items) { $item->relation }` ou équivalent

**Requêtes sans index**
- `WHERE`, `ORDER BY`, `JOIN ON` sur des colonnes sans index identifié dans le schéma
- Comparer avec le schéma BDD (`docs/uml/erd.*` ou migrations)

**Appels synchrones parallélisables**
- Appels API / BDD séquentiels sans dépendance entre eux → candidats `Promise.all` / goroutines / async
- Appels en boucle vers des services externes

**Cache manquant**
- Données statiques ou peu changeantes recalculées à chaque requête
- Résultats de requêtes coûteuses non mis en cache (Redis, mémoire)

### 3 — Analyse Frontend

**Re-renders inutiles**
- Composants React/Vue sans `memo` / `useMemo` / `computed` sur des données stables
- Props objets recréés à chaque render (ex: `style={{ }}` inline, callbacks non mémorisés)

**Bundle**
- Imports de librairies entières au lieu d'imports ciblés (`import _ from 'lodash'` vs `import debounce`)
- Code non tree-shakable
- Chunks trop volumineux sans lazy loading (`import()` dynamique)

**Réseau**
- Appels API sans debounce sur des événements à haute fréquence (input, scroll, resize)
- Appels dupliqués à la même ressource sans deduplication
- Images sans `lazy` loading, sans format optimisé (WebP), sans `srcset`

### 4 — Analyse Algorithmes
- Boucles imbriquées O(n²) évitables par une structure Map/Set ou un tri préalable
- Recherches linéaires répétées sur des collections non indexées
- Calculs répétés sans memoization dans la logique métier

### 5 — Rapport
Générer `docs/perf-[DATE].md` :

```markdown
# Audit Performance — [DATE]
> Scope : [module ou "global"] | Stack : [stack détectée]

## Résumé exécutif
[2-3 phrases sur les problèmes les plus impactants]

## Findings

| # | Problème | Type | Localisation | Impact | Effort | Priorité |
|---|----------|------|-------------|--------|--------|----------|
| P-01 | N+1 sur UserRepository | Backend | `src/services/UserService.ts:45` | Élevé | Faible | Critique |
| P-02 | Bundle lodash complet | Frontend | `src/utils/date.ts:2` | Moyen | Faible | Majeur |

## Détail des Critiques

### P-01 : [Description]
**Localisation** : `fichier:ligne`  
**Problème** : [explication technique]  
**Solution recommandée** : [correction concrète]  
**Gain estimé** : [ordre de grandeur]
```

Post-génération : proposer `@fix P-01` pour les items Critique.
