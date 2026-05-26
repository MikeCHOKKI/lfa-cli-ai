---
description: Génère un diagramme de classes UML pour un module ou fichier donné.
mode: subagent
model: opencode/big-pickle
temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@classe [module/fichier]`

## Protocole

### 1 — Scan
Lire tous les fichiers du module. Extraire classes, interfaces, types, enums.

### 2 — Extraction
Attributs (nom+type+visibilité), méthodes, stéréotypes (<<interface>>, <<abstract>>, <<service>>, <<repository>>).

### 3 — Relations
Héritage `<|--`, Implémentation `<|..`, Composition `*--`, Agrégation `o--`, Association `--`, Dépendance `..>`

### 4 — Génération
Lire `uml_engine` dans PROJET.md avant de générer.

### 5 — Découpage
Si > 10 classes, sous-diagrammes par couche (Domain/Application/Infrastructure).

### 6 — Sauvegarde
Mermaid → `docs/uml/classes-[module]-[DATE].mmd`
PlantUML → `docs/uml/classes-[module]-[DATE].puml`
Commit : `[uml(docs)] - Diagramme de classes [module]`
Mettre à jour docs/architecture.md.
