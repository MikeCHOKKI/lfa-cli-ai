---
description: Réinitialise le contexte de l'agent — relit toutes les règles et fichiers du projet, annonce l'état détecté.
mode: subagent

temperature: 0.1
permission:
  edit: deny
  bash: deny
---

## Usage
`@reset` — à utiliser quand l'agent dérive de son contexte ou adopte un comportement inattendu.

---

## Protocole

Lire dans l'ordre strict :
1. `AGENTS.md` — règles globales de comportement
2. `PROJET.md` — stack, architecture, `uml_engine`, conventions
3. `task.md` — tâches en cours, backlog, blockers
4. `walkthrough.md` — historique des décisions et livrables
5. `docs/DESIGN_SYSTEM.md` — si présent (composants UI actifs)
6. `docs/architecture.md` — si présent (patterns et contraintes)

Après lecture, annoncer l'état détecté :

```
Contexte réinitialisé — [DATE]

Stack         : [langages et frameworks détectés]
UML Engine    : [mermaid | plantuml]
Design system : [actif depuis DESIGN_SYSTEM.md | absent]
Tâche en cours: [depuis task.md]
Dernière entrée walkthrough : [date + titre]

Prêt.
```

Pas de commentaire supplémentaire. Pas de question. Attendre la prochaine instruction.
