---
description: Gestion base de données — ERD, migrations, seeds, audit qualité schéma.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: ask
---

## Usage
`@db schema` · `@db migration [description]` · `@db seed [module]` · `@db audit`

---

## `@db schema`

1. Lire `PROJET.md` → `uml_engine` (mermaid ou plantuml)
2. Scanner tous les modèles / entités / migrations existants
3. Extraire : tables, colonnes (nom + type), clés primaires, clés étrangères, index, contraintes
4. Générer `docs/uml/erd.[mmd|puml]`
5. Si > 15 tables : découper en sous-diagrammes par module fonctionnel + un ERD global simplifié
6. Commit : `uml(docs): Schéma ERD mis à jour`

---

## `@db migration [description]`

1. Lire les migrations existantes pour identifier la convention de nommage et le format utilisé
2. Vérifier que la migration ne duplique pas une existante (même table, même colonne)
3. Créer `migrations/[TIMESTAMP]_[description].sql`

Structure obligatoire :
```sql
-- Migration : [description]
-- Date : [DATE]

-- UP
[instructions de création/modification]

-- DOWN
[instructions de rollback exact et complet]
```

Règles strictes :
- **Jamais** `DROP TABLE` ou `DROP COLUMN` sans confirmation explicite de l'utilisateur
- Le `DOWN` doit annuler exactement le `UP`, rien de plus
- Toute nouvelle colonne NOT NULL doit avoir une valeur DEFAULT
- Tout ajout de FK doit vérifier l'existence de l'index correspondant

4. Commit : `build(db): Migration [description]`

---

## `@db seed [module]`

1. Lire le schéma des tables concernées
2. Respecter les contraintes FK (ordre d'insertion, IDs cohérents entre tables)
3. Données réalistes : noms, emails, montants, dates cohérentes avec le contexte métier
4. Volume : 10 à 50 entrées par table principale, proportionnel pour les tables de liaison
5. Créer `seeds/[module]-seed.sql`
6. Commit : `build(db): Seed [module]`

---

## `@db audit`

Scanner l'ensemble du schéma et signaler :

| Catégorie | Problème détecté |
|-----------|-----------------|
| Index | Colonnes FK sans index, colonnes fréquemment filtrées sans index |
| Intégrité | Relations sans contrainte FK explicite |
| Nullable | Colonnes NOT NULL sans valeur DEFAULT |
| Timestamps | Tables sans `created_at` / `updated_at` |
| Nommage | Incohérences (snake_case vs camelCase, pluriel vs singulier) |
| Types | VARCHAR sans longueur, TEXT pour des données courtes, INT pour des IDs qui devraient être UUID |
| Sécurité | Colonnes `password`/`token` en clair (non hashé détectable par le type TEXT court) |

Générer `docs/db-audit-[DATE].md` avec : finding, table concernée, impact, correctif recommandé.
Commit : `docs(db): Audit schéma [DATE]`
