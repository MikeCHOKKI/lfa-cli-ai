---
name: delegation
description: >
  Règles de délégation et d'orchestration entre agents. Utilisé par l'orchestrateur
  et les subagents pour savoir quand déléguer une sous-tâche, comment formuler
  la requête, et quoi vérifier avant de rendre la main.
---

# Delegation Skill

Règles de délégation entre agents.

## Quand déléguer

| Situation | Action |
|-----------|--------|
| Requête hors scope | Router vers @agent compétent |
| Tâche complexe (3+ étapes) | Planifier, puis déléguer chaque étape |
| Compétence spécialisée | Charger le skill via `task` |
| Ambiguïté bloquante | @question avant d'agir |
| Erreur répétée | @fix avec le contexte |

## Skills disponibles

| Skill | Domaine |
|-------|---------|
| `code-quality` | Clean code, SOLID, patterns |
| `project-standards` | Structure, CI/CD, Docker |
| `design-system` | Palettes, typo, tokens |
| `mockup-ui` | Maquettes HTML/CSS |
| `animation` | CSS/SVG animations |
| `svg-art` | Diagrammes, logos, art |
| `design-import` | Figma/Stitch → code |

## Pattern de délégation

```
Analyse → découpage → [sous-tâche → @agent → vérification] × N → synthèse
```

Toujours vérifier que la sous-tâche est atomique (une seule responsabilité).
