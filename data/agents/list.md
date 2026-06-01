---
description: Inventorie toutes les fonctionnalités du projet et génère docs/features.md à partir du code source réel.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@list`

---

## Protocole

### 1 — Collecte
- Lire `PROJET.md`, `task.md`, `docs/architecture.md` — vue macro existante
- Scanner le codebase pour extraire les fonctionnalités **réellement implémentées** :
  - Routes / endpoints : `glob **/routes/**`, `grep -r "@Route\|router\.\|app\.\(get\|post\|put\|delete\|patch\)"` 
  - Composants UI : `glob **/components/**/*.{tsx,vue,dart}`
  - Services métier : `glob **/services/**`
  - Tâches planifiées / jobs : `glob **/jobs/**`, `glob **/cron/**`
  - Événements / listeners : `glob **/events/**`, `glob **/listeners/**`

### 2 — Classification
Pour chaque fonctionnalité identifiée, déterminer :
- **Module** : sous-système fonctionnel (auth, produits, commandes, notifications, etc.)
- **Statut** :
  - `Implémenté` — code présent, logique complète
  - `Partiel` — code présent mais incomplet (TODO visible, cas manquants)
  - `En cours` — branch ou PR active (si détectable)
  - `Planifié` — mentionné dans `task.md` ou `docs/` mais sans code
- **Tests** : `Oui` / `Non` / `Partiel` (présence de fichiers de test pour le module)
- **Endpoint** : route principale si API, `—` si logique interne

### 3 — Génération
Créer ou écraser `docs/features.md` :

```markdown
# Features — [Nom du projet]
> Généré le [DATE] | [N] fonctionnalités recensées

## Résumé
| Statut | Nombre |
|--------|--------|
| Implémenté | N |
| Partiel | N |
| En cours | N |
| Planifié | N |

## Inventaire

### [Module 1]
| # | Fonctionnalité | Statut | Tests | Endpoint / Localisation |
|---|---------------|--------|-------|--------------------------|
| 1 | Connexion utilisateur | Implémenté | Oui | POST /auth/login |
| 2 | Réinitialisation mot de passe | Partiel | Non | POST /auth/reset |

### [Module 2]
...

## Fonctionnalités planifiées (sans code)
| Fonctionnalité | Source | Priorité |
|---------------|--------|----------|
```

### 4 — Post-génération
- Mettre à jour `task.md` si des fonctionnalités `Partiel` ou `Planifié` ne sont pas tracées
- Proposer `@roadmap` pour construire la feuille de route à partir de cet inventaire
