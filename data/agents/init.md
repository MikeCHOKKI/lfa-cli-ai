---
description: Initialise un projet nouveau ou existant — détecte l'état réel, crée uniquement ce qui manque, ne touche jamais à ce qui existe.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: ask
---

## Usage
`@init` — détection automatique + questions pour ce qui manque
`@init [description]` — idem avec contexte fourni en entrée

---

## Protocole

### 1 — Audit de l'existant (TOUJOURS en premier)

Avant toute question, avant toute création, scanner le répertoire courant :

```bash
# Présence de fichiers de configuration connus
ls -la
find . -maxdepth 3 -name "package.json" -o -name "composer.json" -o -name "pom.xml" \
       -o -name "Cargo.toml" -o -name "go.mod" -o -name "pubspec.yaml" \
       -o -name "build.gradle" -o -name "pyproject.toml" | head -20

# Git
git status 2>/dev/null || echo "NO_GIT"
git log --oneline -5 2>/dev/null || echo "NO_COMMITS"
git remote -v 2>/dev/null || echo "NO_REMOTE"
```

Vérifier la présence de chaque fichier géré par @init :

| Fichier | Présent ? | Action |
|---------|-----------|--------|
| `PROJET.md` | ✓ / ✗ | lire si présent |
| `task.md` | ✓ / ✗ | lire si présent |
| `walkthrough.md` | ✓ / ✗ | lire si présent |
| `AGENTS.md` (local) | ✓ / ✗ | lire si présent |
| `.env.example` | ✓ / ✗ | lire si présent |
| `.gitignore` | ✓ / ✗ | lire si présent |
| `docs/architecture.md` | ✓ / ✗ | — |
| `docs/features.md` | ✓ / ✗ | — |
| `docs/DESIGN_SYSTEM.md` | ✓ / ✗ | — |
| `docs/api.md` | ✓ / ✗ | — |

Si `PROJET.md` existe → le lire entièrement. Il est la source de vérité du projet.

---

### 2 — Déterminer le mode

#### Mode A — Projet existant sans fichiers @init
Critères : fichiers source présents (code, manifestes) mais `PROJET.md` / `task.md` absents.

Actions :
1. Analyser la stack depuis les manifestes (`package.json`, `pom.xml`, etc.)
2. Analyser la structure des dossiers pour déduire le pattern architectural
3. Lire quelques fichiers source représentatifs pour comprendre les conventions
4. Inférer le contexte métier depuis les noms de modules, routes, modèles
5. Poser uniquement les questions non déductibles (max 3)
6. Annoncer : `Projet existant détecté — initialisation des fichiers de gouvernance manquants.`

#### Mode B — Projet partiellement initialisé
Critères : certains fichiers @init présents, d'autres manquants.

Actions :
1. Lire tous les fichiers @init existants
2. Extraire le contexte (stack, uml_engine, conventions) depuis `PROJET.md` s'il existe
3. Identifier précisément ce qui manque
4. Créer uniquement les fichiers manquants, en cohérence avec l'existant
5. Annoncer : `Initialisation partielle détectée — complétion des [N] fichiers manquants.`

#### Mode C — Projet vierge
Critères : répertoire vide ou quasi-vide (pas de code source, pas de manifeste).

Actions :
1. Si description fournie → extraire stack, type, contexte
2. Poser les questions essentielles (max 5)
3. Créer la structure complète
4. Annoncer : `Nouveau projet — initialisation complète.`

**Règle absolue** : ne jamais écraser un fichier existant. Si un fichier existe → le lire, éventuellement le compléter si du contenu manque (section absente), jamais le remplacer.

---

### 3 — Questions (uniquement ce qui ne peut pas être déduit)

```
1. Quel est le nom du projet ?              [si non trouvable dans package.json / pom.xml / dossier]
2. Stack principale ?                       [si non déduite des manifestes]
3. Type de BDD ?                            [si applicable et non détectable]
4. UML engine ? [mermaid | plantuml]        [défaut : mermaid — demander uniquement si non défini dans PROJET.md]
5. Déploiement cible ?                      [Docker, bare metal, cloud, k8s — si non documenté]
```

⛔ Attendre les réponses avant toute création ou modification.

---

### 4 — Inférence depuis un projet existant

Quand le code source est présent, extraire activement :

**Stack** — depuis : `package.json` (dependencies), `composer.json` (require), `pom.xml` (dependencies), `go.mod` (require), `pubspec.yaml` (dependencies), `Cargo.toml` (dependencies)

**Framework** — depuis : présence de `next.config.*`, `angular.json`, `artisan`, `symfony.lock`, `spring-boot`, structure `cmd/` (Go), `lib/main.dart` (Flutter)

**BDD** — depuis : variables dans `.env.example` ou `.env` (si lisible), config files (`database.yml`, `datasource.*`), noms de packages ORM (TypeORM, Prisma, Eloquent, Hibernate, GORM, Diesel, SQLx)

**Pattern architectural** — depuis : structure des dossiers (`controller/service/repository` → MVC/Clean, `features/` → feature-first, `domain/application/infrastructure` → Hexagonal)

**Contexte métier** — depuis : noms de routes, noms de modèles/entités, noms de services

**Conventions existantes** — depuis : analyse de 3-5 fichiers source représentatifs (indentation, naming, structure des fonctions)

**Git** — remote URL pour déduire le nom de projet et l'organisation

---

### 5 — Fichiers à créer ou compléter

Ne créer un fichier que s'il est **absent**. Ne compléter une section que si elle est **manquante dans un fichier existant**.

#### `PROJET.md` (créer si absent, compléter les sections vides si présent)
```markdown
# [Nom du projet]
> [Date] | [Statut : En cours]

## Description
[Description fournie ou inférée depuis le code]

## Stack
[Inférée ou fournie]
- Runtime : [Node.js / PHP / JVM / Go / Dart / Rust]
- Framework : [Next.js / Slim / Spring Boot / Gin / Flutter / etc.]
- BDD : [PostgreSQL / MariaDB / MongoDB / SQLite / etc.]
- Outils : [Docker, Redis, NATS, etc.]

## Architecture
Pattern : [MVC / Clean / Hexagonal / MVVM / Feature-first / etc.]
[Justification depuis la structure observée]

## Modules
[Liste inférée depuis les dossiers et routes existants]

## Configuration
uml_engine: mermaid
deploy_target: [cible]

## Conventions
- Commits : type(scope): titre (conventional commits, français)
- Branches : feat/, fix/, refactor/, release/
- Nommage : [inféré depuis le code existant]
```

#### `task.md` (créer si absent)
```markdown
# Tasks — [Nom du projet]
> Mise à jour : [DATE]

## En cours
[Inféré depuis le code : fonctionnalités partielles, TODO visibles]

## Backlog
[Inféré depuis la description ou les fichiers]

## Fait
- [x] @init — [DATE] ([mode : nouveau | existant | complétion])
```

#### `walkthrough.md` (créer si absent)
```markdown
# Walkthrough — [Nom du projet]
> Historique des décisions techniques et livrables.

---

## [DATE] — @init [mode]
[Mode A/B/C — description de ce qui a été fait]
Stack détectée : [stack]
Fichiers créés/complétés : [liste]
[Si projet existant : résumé de ce qui existait déjà]
```

#### `AGENTS.md` local (créer si absent)
```markdown
# AGENTS.md — [Nom du projet]
> Règles spécifiques à ce projet. Complète les règles globales.

## Contexte projet
- Stack : [stack]
- UML Engine : [mermaid | plantuml]
- Pattern : [pattern]
- BDD : [type et ORM]

## Conventions spécifiques
[Inférées depuis le code — indentation, naming, structure]

## Modules principaux
[Liste inférée avec description courte]

## Commandes utiles
\`\`\`bash
# Développement
[inférée depuis package.json scripts / Makefile / README]

# Tests
[inférée]

# Build
[inférée]
\`\`\`
```

#### `.env.example` (créer si absent — ne jamais toucher au `.env` réel)
Variables inférées depuis :
- `.env` existant → copier les clés sans les valeurs
- Config files → variables référencées (`process.env.X`, `$_ENV['X']`, `os.Getenv("X")`)
- Dépendances → variables standards connues (Redis `REDIS_URL`, S3 `AWS_*`, etc.)

#### `.gitignore` (créer si absent, compléter si des entrées manquent)
Toujours vérifier que ces entrées sont présentes :
```
.env
.env.local
.opencode/
.agents/
```
Ajouter les entrées manquantes selon la stack. Ne pas dupliquer l'existant.

#### `docs/` (créer le dossier et les fichiers manquants)
- `docs/architecture.md` — si absent
- `docs/features.md` — si absent
- `docs/DESIGN_SYSTEM.md` — si absent et projet avec UI
- `docs/api.md` — si absent et projet avec API/endpoints

#### Fichiers stack-specific (créer si absents et pertinents)
| Stack | Fichiers |
|-------|----------|
| Next.js + Tailwind | `tailwind.config.ts`, `src/lib/utils.ts` |
| PHP | `src/bootstrap.php` si absent |
| Spring Boot | `src/main/resources/application.yml` si absent |
| Go | `Makefile` si absent |
| Flutter | `lib/app/app.dart` si absent |
| Docker détecté ou requis | `Dockerfile`, `docker-compose.yml`, `.dockerignore` |
| Gitea/GitHub CI non configuré | `.gitea/workflows/ci.yml` ou `.github/workflows/ci.yml` |

---

### 6 — Git

**Projet sans git** → `git init`, puis commit initial.
**Projet avec git existant** → ne pas réinitialiser, ne pas modifier l'historique.
**Après créations/modifications** → proposer un commit :
```
chore(init): Ajout des fichiers de gouvernance [liste courte]
```
Ne commiter qu'après confirmation explicite.

---

### 7 — Rapport final

```
@init — [MODE : Nouveau | Existant | Complétion partielle]
Projet  : [Nom]
Stack   : [stack détectée/confirmée]
UML     : [mermaid | plantuml]
Git     : [initialisé | existant — N commits — remote: url | existant — pas de remote]

Existant conservé (non modifié) :
  • [fichiers présents et intacts]

Créé :
  ✓ [fichier] — [raison si non évident]

Complété (sections ajoutées) :
  ~ [fichier] — [sections ajoutées]

Ignoré (non applicable) :
  — [fichier] — [raison]

Commit proposé : chore(init): [description]

Prochaines étapes :
  @list        → inventorier les fonctionnalités existantes  [si projet existant]
  @usecases    → définir les cas d'utilisation               [si nouveau]
  @feat [...]  → première fonctionnalité                     [si nouveau]
  @db schema   → modéliser le schéma BDD                     [si applicable]
  @conformite  → vérifier la conformité globale              [si projet existant]
```
