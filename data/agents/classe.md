---
description: Génère un diagramme de classes UML précis et lisible à partir du code source d'un module.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@classe [module ou chemin]`

---

## Protocole

### 1 — Contextualisation
- Lire `PROJET.md` → `uml_engine`. Ne jamais mélanger mermaid et PlantUML dans un projet.
- Lire `docs/architecture.md` → comprendre le découpage en couches (Domain, Application, Infrastructure)
- Scanner le module ciblé : `glob [module]/**/*.{ts,php,java,go,rs,dart,py}`

### 2 — Extraction des éléments
Pour chaque fichier :
- **Classes, interfaces, types, enums, traits** → nom + stéréotype (`<<interface>>`, `<<abstract>>`, `<<service>>`, `<<repository>>`, `<<entity>>`, `<<value-object>>`)
- **Attributs** : nom, type, visibilité (`+` public, `-` privé, `#` protégé)
- **Méthodes** : nom, paramètres typés, type de retour, visibilité
- Exclure : getters/setters triviaux, méthodes de framework non significatives

### 3 — Relations
| Notation | Relation      | Critère |
|----------|---------------|---------|
| `<\|--` | Héritage      | `extends` |
| `<\|..` | Implémentation | `implements` |
| `*--` | Composition   | Cycle de vie dépendant |
| `o--` | Agrégation    | Référence sans ownership |
| `-->` | Association   | Utilisation directe |
| `..>` | Dépendance    | Import ou injection ponctuelle |

Indiquer les multiplicités (`1`, `0..1`, `*`, `1..*`) quand pertinent.

### 4 — Découpage si > 10 classes
Produire des sous-diagrammes par couche :
- `classes-[module]-domain-[DATE]`
- `classes-[module]-application-[DATE]`
- `classes-[module]-infrastructure-[DATE]`
- Un diagramme d'overview avec relations inter-couches uniquement

### 5 — Génération
- Mermaid → `docs/uml/classes-[module]-[DATE].mmd`
- PlantUML → `docs/uml/classes-[module]-[DATE].puml`

Règles de lisibilité :
- Regrouper par couche dans le fichier
- Limiter les attributs affichés aux champs métier significatifs (pas les métadonnées framework)
- Si PlantUML : utiliser `skinparam` pour la lisibilité (couleurs par stéréotype)

### 6 — Documentation
- Mettre à jour `docs/architecture.md` → section "Modèle de classes"
- Commit : `uml(docs): Diagramme de classes [module]`
