# AGENTS.md — Directives globales opencode

> Ce fichier est lu par l'agent à chaque session. Il définit le comportement par défaut, les conventions, et les règles de travail.

---

## Comportement général

- Lire les fichiers pertinents **avant** d'agir — jamais d'hypothèse sur le contenu
- Comprendre la demande dans son intégralité avant de commencer
- En cas d'ambiguïté bloquante : poser **une seule question**, précise, puis attendre
- Travailler de manière incrémentale : une sous-tâche à la fois, vérifier, continuer

---

## Lecture avant action

- Toujours `read` un fichier avant `edit` ou `write`
- Toujours `glob` ou `grep` pour trouver les conventions existantes avant d'écrire du nouveau code
- Vérifier les imports et dépendances existants avant d'en ajouter de nouveaux

---

## Outils — ordre de priorité

1. `glob` / `grep` / `read` — exploration et compréhension
2. `edit` — modifier un fichier existant
3. `write` — créer un nouveau fichier
4. `bash` — uniquement si impossible autrement

---

## Code

- Respecter strictement le style du fichier ciblé (indentation, naming, quotes, structure)
- Pas de commentaires sauf si explicitement demandé
- Pas de README, CHANGELOG, documentation sauf si demandé
- Pas de refactoring non demandé : corriger le bug, rien de plus
- Pas d'imports, dépendances ou packages non utilisés

---

## Agents et skills

- Invoquer `@agent` ou charger un `skill` uniquement si la tâche le requiert
- Ne pas chaîner des sous-agents inutilement
- Utiliser `task` (sous-agent isolé) uniquement si la complexité le justifie réellement

---

## Réponses — format et concision

> ⚠️ Ces règles s'appliquent au **style de réponse**, pas à la **profondeur du travail produit**.
> L'agent doit traiter **toutes les instructions** reçues, même courtes dans sa réponse.

- Pas de phrase d'introduction ("Bien sûr", "Je vais", "Voici", "D'accord")
- Pas de reformulation de la tâche avant de l'exécuter
- Pas de conclusion bavarde ("N'hésite pas à me dire si...")
- **Si une tâche nécessite un fichier détaillé, un plan complet, un rapport** → le produire en entier, sans tronquer
- **Si la réponse conversationnelle suffit** → 1 à 3 lignes maximum
- La concision s'applique aux **échanges**, pas aux **livrables**

---

## Notifications

- Tâche terminée → `notify-send "OpenCode" "Tâche terminée"` (notification bureau)
- Demande d'autorisation → `notify-send "OpenCode" "Action requise"`

## Commits

- Jamais de `git commit` sans demande explicite
- Format : `type(scope): Titre court en français` (conventional commits)
- Types autorisés : `feat`, `fix`, `refactor`, `ci`, `docs`, `test`, `chore`, `perf`, `secu`
- Exemples :
  - `feat(auth): Ajout du flux de réinitialisation de mot de passe`
  - `fix(api): Correction de la pagination sur /products`
  - `secu(deps): Mise à jour lodash suite CVE-2024-XXXX`

---

## Sécurité

- Ne jamais logger, afficher ou inclure dans le code : clés API, tokens, mots de passe, secrets
- Si un secret est détecté dans le diff ou un fichier → **signaler en une ligne et stopper immédiatement**
- Ne jamais commiter de `.env` contenant des valeurs réelles

---

---

## Sous-agents disponibles

### Développement
| Commande              | Rôle                                                               |
|-----------------------|--------------------------------------------------------------------|
| `@feat [description]` | Nouvelle fonctionnalité — plan validé, implémentation incrémentale |
| `@fix [bug/erreur]`   | Correction ciblée — cause racine, plan minimal, vérification       |
| `@ui [composant]`     | Composant UI — designer mindset, conforme au design system         |
| `@test [module]`      | Tests unitaires, intégration, E2E — couverture complète            |

### Qualité & Audit
| Commande          | Rôle                                                                |
|-------------------|---------------------------------------------------------------------|
| `@audit [scope]`  | Audit complet — architecture, qualité, sécurité, performance, tests |
| `@securite`       | Audit sécurité — secrets, CVE, auth, validation, BDD, headers       |
| `@perf [scope]`   | Audit performance — N+1, bundle, algorithmes, cache                 |
| `@verify [scope]` | Vérification fonctionnelle, technique et conformité                 |
| `@conformite`     | Conformité projet — structure, code, commits, sécurité de base      |

### Architecture & Documentation
| Commande              | Rôle                                                            |
|-----------------------|-----------------------------------------------------------------|
| `@question [dilemme]` | Dilemme architectural — 3 approches, recommandation, ADR        |
| `@details [sujet]`    | Analyse approfondie — faisabilité, risques, plan d'action       |
| `@roadmap`            | Feuille de route — phases, sprint, backlog, dette technique     |
| `@list`               | Inventaire des fonctionnalités → `docs/features.md`             |
| `@usecases [module]`  | Cas d'utilisation — identification, priorisation, documentation |

### UML
| Commande                | Rôle                                                              |
|-------------------------|-------------------------------------------------------------------|
| `@uml [type] [module]`  | Dispatcher UML — séquence, activité, composant, déploiement, état |
| `@sequence [flux]`      | Diagramme de séquence — interactions système basées sur le code   |
| `@activity [processus]` | Diagramme d'activité — workflows et processus métier              |
| `@classe [module]`      | Diagramme de classes — entités, relations, stéréotypes            |

### Base de données
| Commande               | Rôle                         |
|------------------------|------------------------------|
| `@db schema`           | ERD depuis les modèles       |
| `@db migration [desc]` | Migration SQL avec UP/DOWN   |
| `@db seed [module]`    | Seeds réalistes et cohérents |
| `@db audit`            | Audit qualité du schéma      |

### Ops & Workflow
| Commande                        | Rôle                                                          |
|---------------------------------|---------------------------------------------------------------|
| `@init [description]`           | Initialisation projet — structure, fichiers, git, conventions |
| `@commit`                       | Message de commit professionnel depuis le diff                |
| `@deploy [staging\|production]` | Checklist pré-déploiement — bloque sur item critique          |
| `@reset`                        | Réinitialisation du contexte agent                            |
