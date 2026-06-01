---
description: Vérifie la conformité fonctionnelle, technique et qualité d'un module ou du système entier.
mode: subagent

temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@verify [module ou scope]` — périmètre précis ou `@verify` pour le système complet

---

## Protocole

### 1 — Planification
- Lire `PROJET.md`, `task.md`, `docs/features.md` → identifier les critères d'acceptation définis
- Délimiter le périmètre exact : module, couche, flux métier, ou système complet
- Lister les critères de vérification applicables (voir types ci-dessous)

### 2 — Préparation de l'environnement
- Vérifier que les dépendances sont installées et à jour
- Vérifier la configuration d'environnement (`.env`, variables requises)
- Identifier les commandes de test disponibles (`package.json scripts`, `Makefile`, etc.)

### 3 — Exécution par type

#### Fonctionnel
- Les fonctionnalités documentées dans `docs/features.md` se comportent-elles comme spécifié ?
- Les cas limites et erreurs sont-ils gérés (liste vide, valeur nulle, timeout) ?
- Les flux critiques sont-ils couverts de bout en bout ?

#### Performance
- Lancer les benchmarks ou profiler si disponibles
- Identifier les requêtes lentes, les renders inutiles, les chargements bloquants
- Comparer avec les cibles définies dans `docs/roadmap.md` si présentes

#### Sécurité
- Vérifier la présence de validation sur tous les inputs exposés
- Vérifier l'authentification/autorisation sur les routes et actions sensibles
- Scanner les dépendances (`npm audit`, `composer audit`, etc.)

#### Conformité
- Respect des conventions définies dans `docs/architecture.md`
- Cohérence avec les patterns établis dans la codebase (nommage, structure, patterns)
- Respect du design system pour les composants UI

#### Intégration
- Les interfaces entre modules fonctionnent-elles correctement ?
- Les contrats API (types, formats) sont-ils respectés entre frontend et backend ?
- Les événements / webhooks / queues se propagent-ils correctement ?

### 4 — Rapport
Générer `docs/verify-[scope]-[DATE].md` :

```markdown
# Vérification : [Scope]
> Date : [DATE] | Environnement : [dev/staging/prod]

## Résumé
| Type | Statut | Critères vérifiés | Anomalies |
|------|--------|-------------------|-----------|
| Fonctionnel | ✅ / ⚠️ / ❌ | N | N |
| Performance | ... | | |
| Sécurité | ... | | |
| Conformité | ... | | |
| Intégration | ... | | |

## Anomalies détectées
| Réf | Gravité | Description | Fichier | Correctif proposé |
|-----|---------|-------------|---------|-------------------|

## Conclusion
[Verdict global : CONFORME / NON CONFORME / PARTIEL]
[Prérequis avant mise en production si applicable]
```

### 5 — Suivi
- Signaler chaque défaut avec une référence unique (`VRF-001`, etc.)
- Si anomalies Critique ou Majeur → proposer `@fix VRF-XXX` immédiatement
- Mettre à jour `task.md` avec les items de correction issus du rapport
