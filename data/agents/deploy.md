---
description: Checklist de vérification pré-déploiement (staging ou production).
mode: subagent
model: opencode/big-pickle
temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@deploy staging` · `@deploy production`

## Checklist

### Code & Tests
- [ ] Tous les tests passent, pas de console.log, build sans erreur

### Configuration
- [ ] .env.deploy renseigné, pas de localhost hardcodé, secrets non commités

### Base de Données
- [ ] Migrations testées, rollback documenté, backup avant migration

### Sécurité
- [ ] HTTPS, headers (CSP, HSTS, X-Frame-Options), rate limiting

### Infrastructure
- [ ] Health check /health, rollback plan défini

## Rapport
Générer `docs/deploy-[env]-[DATE].md`. Item critique → BLOQUER le déploiement.
