---
description: Gestion base de données — schéma ERD, migrations, seeds, audit BDD.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: ask
  bash: ask
---

## Usage
`@db schema` · `@db migration [description]` · `@db seed` · `@db audit`

## @db schema
Lire les modèles → ERD dans `docs/uml/erd.mmd` (ou .puml selon uml_engine).
Commit : `[uml(docs)] - Schéma ERD mis à jour`

## @db migration
1. Fichier : `migrations/[TIMESTAMP]_[description].sql`
2. Obligatoire : `-- UP` et `-- DOWN`
3. Jamais DROP TABLE/COLUMN sans confirmation explicite
4. Commit : `[build(db)] - Migration [description]`

## @db seed
Données réalistes et cohérentes entre tables, 10-50 entrées.
Fichier : `seeds/[module]-seed.sql`

## @db audit
Signaler : colonnes sans index, relations sans FK, NOT NULL sans défaut, tables sans timestamps, nommage incohérent.
Rapport : `docs/db-audit-[DATE].md`
