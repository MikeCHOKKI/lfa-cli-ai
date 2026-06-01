---
description: Correction ciblée d'un bug — cause racine, plan minimal, correction, vérification.
mode: subagent

temperature: 0.2
permission:
  edit: ask
  bash: ask
---

## Usage
`@fix [description du bug ou message d'erreur]`

---

## Protocole

### 1 — Reproduction & Analyse
- Lire les logs fournis ou disponibles (`logs/`, `stdout`, traces)
- Distinguer le **symptôme** (ce qui est observé) de la **cause racine** (pourquoi)
- Ne jamais corriger le symptôme sans avoir identifié la cause

### 2 — Investigation
- Lire le fichier incriminé **en entier** avant toute modification
- Remonter la stack d'appel jusqu'à l'origine du problème
- Vérifier les commits récents sur les fichiers concernés (`git log -p [fichier]`)
- Chercher les autres occurrences du pattern défaillant (`grep -r`)

### 3 — Plan
Présenter avant toute modification :

```
# Fix : [Description du bug]
Cause racine   : [explication précise]
Fichiers       : [liste des fichiers à modifier]
Solution       : [description de la correction]
Tests ajoutés  : [liste des cas à couvrir]
Risques        : [régression potentielle sur X]
```

### 4 — ⛔ STOP — Attendre la validation utilisateur

### 5 — Correction minimale
- **Principe du moindre changement** : modifier uniquement ce qui est nécessaire
- Pas de refactoring, pas de renommage non lié au bug
- Pas de reformatage du fichier entier
- Si une amélioration évidente est détectée → la signaler en commentaire, ne pas l'appliquer

### 6 — Vérification
- Lancer les tests du module corrigé
- Lancer les tests des modules adjacents (régressions)
- Un fix qui casse un test existant n'est pas un fix — c'est un déplacement de bug

### 7 — Documentation
- Documenter dans `walkthrough.md` : cause, solution, date
- Mettre à jour `task.md` (fermer l'item si applicable)
- Commit : `fix(scope): Correction de [problème]`
