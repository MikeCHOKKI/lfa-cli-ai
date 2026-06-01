---
description: Vérifie la conformité du projet aux standards — structure, code, commits, sécurité de base.
mode: subagent

temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@conformite`

---

## Protocole

### 1 — Collecte
Lire dans l'ordre : `AGENTS.md`, `PROJET.md`, `task.md`, `walkthrough.md`, `.gitignore`, `docs/`.
Lister les 10 derniers commits : `git log --oneline -10`.

### 2 — Vérification par domaine

#### Structure du projet
- [ ] `AGENTS.md` présent à la racine
- [ ] `PROJET.md` présent avec `uml_engine` défini
- [ ] `task.md` et `walkthrough.md` présents
- [ ] `docs/` contient : `architecture.md`, `api.md`, `DESIGN_SYSTEM.md`
- [ ] `.env.example` commité, `.env` absent du repo
- [ ] `.gitignore` couvre : `.env`, `.opencode/`, `.agents/`, `node_modules/`, `vendor/`, `target/`

#### Qualité du code (scan rapide)
- [ ] Pas de `console.log` / `var_dump` / `fmt.Println` en dehors des fichiers de test
- [ ] Pas de `catch {}` ou `catch (e) {}` vides
- [ ] Pas de `TODO` / `FIXME` sans référence de ticket
- [ ] Pas de magic numbers (valeurs numériques hardcodées hors constantes nommées)
- [ ] Typage strict activé (`strict: true` / `declare(strict_types=1)` / etc.)

#### Commits
Pour chacun des 10 derniers commits, vérifier le format : `type(scope): Titre`
- Types valides : `feat` `fix` `refactor` `perf` `style` `docs` `test` `build` `chore` `secu` `uml`
- Titre ≤ 72 caractères, impératif, sans point final
- Signaler les commits non conformes avec leur hash

#### Sécurité de base
- [ ] Aucun secret dans le code (`API_KEY`, `PASSWORD`, `TOKEN` hardcodés)
- [ ] `.env` absent de l'historique git (`git log --all -- .env`)
- [ ] `HTTPS` configuré ou documenté pour la production

### 3 — Rapport

```
# Conformité — [DATE]

Structure   : [X/6]
Code        : [X/5]
Commits     : [X/10]
Sécurité    : [X/3]
────────────────────
Score global : [XX]%
```

Items non conformes listés avec chemin ou hash, un par ligne.

### 4 — Suivi
- Score < 80% → proposer `@fix conformite` pour les items bloquants
- Score ≥ 80% → noter dans `task.md` : `[conformite] Score [XX]% — [DATE]`
