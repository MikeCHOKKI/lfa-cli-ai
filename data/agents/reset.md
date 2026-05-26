---
description: Réinitialise le contexte de l'agent en relisant toutes les règles et fichiers du projet.
mode: subagent
model: opencode/big-pickle
temperature: 0.1
permission:
  edit: deny
  bash: deny
---

## Usage
`@reset` — quand dérive de contexte ou comportement inattendu.

## Protocole
1. Relire AGENTS.md du projet
2. Relire PROJET.md
3. Relire task.md
4. Relire walkthrough.md
5. Lire uml_engine dans PROJET.md
6. Annoncer :
```
Contexte réinitialisé
- Stack : [détectée]
- UML Engine : [mermaid|plantuml]
- Design actif : [depuis DESIGN_SYSTEM.md]
- Tâche en cours : [depuis task.md]
- Prêt à continuer.
```
