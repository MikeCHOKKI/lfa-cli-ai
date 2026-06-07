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

## Notifications — OBLIGATOIRE (ne jamais ignorer)

> ⚠️ **RÈGLE ABSOULUE** : tout agent DOIT exécuter ces commandes sans exception.

- **Fin de tâche** : `lfa notify "OpenCode" "Tâche terminée"` après la dernière action
- **Action requise** : `lfa notify "OpenCode" "Action requise"` **AVANT** d'attendre une réponse humaine
- **Erreur bloquante** : `lfa notify "OpenCode" "Erreur bloquante: <description courte>"`

Ces `lfa notify` sont déjà autorisés par la config (permission: allow). Multi-plateforme : notify-send (Linux), osascript (macOS), PowerShell (Windows). Aucune excuse pour les sauter.

## Commits

- Jamais de `git commit` sans demande explicite
- Format : `type(scope): Titre court en français` (conventional commits)
- Types autorisés : `feat`, `fix`, `refactor`, `ci`, `docs`, `test`, `chore`, `perf`, `secu`

## Sécurité

- Ne jamais logger, afficher ou inclure dans le code : clés API, tokens, mots de passe, secrets
- Si un secret est détecté dans le diff ou un fichier → **signaler en une ligne et stopper immédiatement**
- Ne jamais commiter de `.env` contenant des valeurs réelles

---

## Sous-agents disponibles

La liste complète des agents et leurs descriptions est dans `opencode.jsonc` (section `agent`).
Les agents principaux : @feat, @fix, @ui, @test, @audit, @code, @standards, @svg, @mockup, @poster, @animation, @palette, @import, @question, @reset, @commit, @deploy, @verify.

Pour la délégation avancée, charger le skill `delegation`.
