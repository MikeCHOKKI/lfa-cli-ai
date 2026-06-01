---
description: Génère des diagrammes UML (séquence, activité, composant, déploiement, état) basés sur le code source réel.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@uml sequence [module]` · `@uml activite [feature]` · `@uml composant` · `@uml deploiement` · `@uml etat [entité]`

---

## Règle fondamentale
Lire `PROJET.md` → `uml_engine` **avant toute génération**.
Ne jamais mélanger Mermaid et PlantUML dans un même projet. 
Si `uml_engine` est absent de `PROJET.md` → demander à l'utilisateur de le définir avant de continuer.

---

## Par type de diagramme

### `@uml sequence [module]`
→ Déléguer à `@sequence [module]` pour le protocole complet.

### `@uml activite [feature]`
→ Déléguer à `@activity [feature]` pour le protocole complet.

### `@uml composant`
Représenter l'architecture des composants du système :
- Identifier les modules principaux et leurs interfaces exposées
- Relations : dépendance `-->`, fourniture `--`, usage `- ->`
- Grouper par couche (Présentation, Application, Infrastructure, Externe)
- Inclure les dépendances externes (BDD, APIs tierces, queues)

### `@uml deploiement`
Représenter l'infrastructure de déploiement :
- Nœuds physiques ou virtuels (serveurs, containers, cloud)
- Artefacts déployés sur chaque nœud
- Protocoles de communication entre nœuds (HTTP, gRPC, AMQP, etc.)
- Basé sur `PROJET.md` (infra définie) ou `docs/architecture.md`

### `@uml etat [entité]`
Modéliser le cycle de vie d'une entité :
- Lire le code source de l'entité (modèle + service)
- Identifier tous les états possibles (champs `status`, `state`, enums)
- Identifier toutes les transitions (méthodes, events, triggers BDD)
- États initial et final clairement marqués
- Gardes (`[condition]`) sur les transitions conditionnelles

---

## Sauvegarde
- Mermaid → `docs/uml/[type]-[module]-[DATE].mmd`
- PlantUML → `docs/uml/[type]-[module]-[DATE].puml`

Post-génération :
- Mettre à jour `docs/architecture.md` → section "Diagrammes disponibles"
- Commit : `uml(docs): Diagramme [type] [module]`
