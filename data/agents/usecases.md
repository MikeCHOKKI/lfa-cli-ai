---
description: Identifie, structure et priorise les cas d'utilisation d'un module ou du système entier.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: deny
---

## Usage
`@usecases [module]` · `@usecases` pour le système complet

---

## Protocole

### 1 — Identification
- Lire `PROJET.md`, `task.md`, `docs/features.md` s'ils existent
- Scanner les routes / endpoints → chaque endpoint est un cas d'utilisation potentiel
- Scanner les services métier → chaque méthode publique significative est un cas d'utilisation
- Compléter avec les fonctionnalités mentionnées dans `task.md` mais sans code

### 2 — Structuration
Pour chaque cas d'utilisation significatif, documenter :

```
## UC-[N] : [Titre en action — verbe + complément]
**Acteur(s)** : [rôle(s) déclencheur(s)]
**Préconditions** : [état requis du système avant déclenchement]
**Scénario principal** :
  1. [étape]
  2. [étape]
  3. [résultat]
**Scénarios alternatifs** :
  - [condition] → [comportement alternatif]
**Exceptions** :
  - [erreur] → [comportement système]
**Postconditions** : [état du système après succès]
```

Ne documenter en détail que les UC de priorité Haute et Moyenne. Les UC Faible → ligne dans le tableau uniquement.

### 3 — Priorisation
Évaluer chaque UC sur :
- **Valeur métier** : impact direct sur l'objectif produit (Haute / Moyenne / Faible)
- **Fréquence d'usage** : quotidien / occasionnel / rare
- **Complexité** : XS / S / M / L / XL
- **Dépendances** : bloque d'autres UC ?

### 4 — Génération
Créer ou mettre à jour `docs/usecases-[module].md` :

```markdown
# Cas d'utilisation — [Module]
> Date : [DATE] | [N] cas recensés

## Tableau de bord
| UC | Titre | Acteur | Priorité | Statut | Complexité |
|----|-------|--------|----------|--------|------------|
| UC-01 | Créer un utilisateur | Admin | Haute | Implémenté | S |
| UC-02 | Réinitialiser mot de passe | Utilisateur | Haute | Partiel | M |

## Détail des UC prioritaires

### UC-01 : Créer un utilisateur
[structure complète]

### UC-02 : ...
```

Statuts : `Implémenté` · `Partiel` · `Planifié` · `À spécifier`

### 5 — Validation & suivi
- Soumettre le document à l'utilisateur pour validation
- Après validation : mettre à jour `task.md` pour les UC `Planifié` non tracés
- Proposer `@feat UC-[N]` pour les UC haute priorité non implémentés
