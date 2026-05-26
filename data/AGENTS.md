# Token Efficiency Rules

## Réponses
- Max 3 lignes sauf si demande explicite de détail
- Zéro phrase d'introduction ("Bien sûr", "Je vais", "Voici")
- Zéro reformulation de la tâche avant de la faire
- Répondre directement à ce qui est demandé

## Lecture avant action
- Toujours `read` un fichier avant `edit` ou `write`
- Toujours `glob` ou `grep` pour trouver les conventions existantes avant d'écrire du code nouveau
- Vérifier les imports/dépendances existants avant d'en ajouter

## Outils — ordre de priorité
1. `glob` / `grep` / `read` pour explorer
2. `edit` pour modifier un fichier existant
3. `write` pour créer un nouveau fichier
4. `bash` seulement si impossible autrement

## Code
- Pas de commentaires sauf demande explicite
- Pas de README, CHANGELOG, documentation sauf demande
- Respecter le style du fichier ciblé (indentation, naming, quotes)
- Pas de refactoring non demandé

## Agents et skills
- Invoquer `@agent` ou charger un `skill` uniquement si la tâche le requiert explicitement
- Ne pas chaîner des sous-agents inutilement
- `task` = sous-agent isolé : n'utiliser que si la complexité le justifie

## Commits
- Jamais de `git commit` sans demande explicite
- Format : `[type(scope)] - Titre court en français`

## Sons (notifications)
- **Autorisation requise** → `printf '\a'` (BIP)
- **Tâche terminée** → `printf '\a'` (BIP)
- Un BIP = attention, deux BIP = terminé

## Sécurité
- Ne jamais logger, afficher ou inclure dans le code : clés API, tokens, secrets
- Si un secret est détecté dans le diff → signaler en 1 ligne et stopper
