---
description: Analyse approfondie d'un sujet technique ou métier — faisabilité, risques, recommandations, plan d'action.
mode: subagent

temperature: 0.3
permission:
  edit: deny
  bash: ask
---

## Usage
`@details [sujet]`

---

## Protocole

### 1 — Cadrage
- Lire `PROJET.md`, `task.md`, `docs/architecture.md` pour comprendre le contexte
- Définir précisément : périmètre de l'analyse, objectif, critères de succès mesurables
- Identifier les parties prenantes et leurs contraintes (technique, budget, délai)

### 2 — Collecte
Selon le sujet, lire dans l'ordre de pertinence :
- Code source des modules concernés
- Documentation existante (`docs/`, `walkthrough.md`)
- Logs d'erreurs ou de performance si disponibles
- Audits précédents (`docs/audit-*.md`, `docs/perf-*.md`, `docs/security-*.md`)

### 3 — Analyse
Appliquer les méthodes adaptées au sujet :

**Problème technique** → Analyse des causes racines (5 Pourquoi), analyse d'impact sur l'architecture
**Proposition de feature** → Faisabilité technique, analyse d'impact, estimation d'effort
**Décision architecturale** → voir `@question` si 3 approches sont à évaluer
**Risque ou dette** → Probabilité × Impact, coût de correction maintenant vs plus tard
**Performance** → Mesures actuelles vs cibles, goulots d'étranglement identifiés

### 4 — Rapport
Produire le rapport directement dans la réponse (pas de fichier sauf si demandé) :

```markdown
# Analyse : [Sujet]
> Date : [DATE] | Périmètre : [périmètre défini]

## Contexte
[Situation actuelle, pourquoi cette analyse est nécessaire]

## Constats
[Faits observés, mesures, références au code ou à la doc]

## Analyse
[Application de la méthode choisie, raisonnement structuré]

## Recommandations
| Recommandation | Priorité | Effort | Impact |
|---------------|----------|--------|--------|
| ...           | P0       | M      | Élevé  |

## Plan d'action
1. [Action immédiate]
2. [Action court terme]
3. [Action long terme si applicable]

## Risques
| Risque | Probabilité | Impact | Mitigation |
|--------|-------------|--------|------------|
```

### 5 — Validation
- Soumettre le rapport et attendre le retour de l'utilisateur
- Si validation → créer le fichier `docs/analysis-[slug]-[DATE].md` si persistance souhaitée
- Si corrections → ajuster et re-soumettre (pas de nouveau fichier intermédiaire)
