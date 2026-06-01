---
description: Checklist de vérification pré-déploiement — bloque sur tout item critique avant staging ou production.
mode: subagent

temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@deploy staging` · `@deploy production`

---

## Protocole

### 1 — Collecte
- Lire `PROJET.md` → stack, variables d'environnement requises, infrastructure cible
- Lire `docs/deploy-*.md` précédents s'ils existent → problèmes connus
- Lire `task.md` → items en cours ou bloquants non résolus

### 2 — Checklist

#### Code & Build
- [ ] Tous les tests passent (`npm test` / `make test` / équivalent)
- [ ] Build sans erreur ni warning bloquant
- [ ] Pas de `console.log`, `var_dump`, `fmt.Println` hors tests
- [ ] Pas de branche `develop` ou feature non mergée incluse dans le déploiement

#### Configuration
- [ ] `.env.deploy` (ou équivalent) renseigné pour l'environnement cible
- [ ] Aucune valeur `localhost` hardcodée hors configuration locale
- [ ] Aucun secret commité dans le repo (`git log --all -S "password"`)
- [ ] Variables d'environnement requises (depuis `.env.example`) toutes définies

#### Base de données
- [ ] Migrations non jouées identifiées et prêtes
- [ ] Procédure de rollback migration documentée
- [ ] Backup de la BDD planifié avant migration (production uniquement)
- [ ] Seeds de production distincts des seeds de développement

#### Sécurité
- [ ] HTTPS configuré et certificat valide (non expiré)
- [ ] Headers de sécurité actifs : `CSP`, `HSTS`, `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`
- [ ] Rate limiting actif sur les endpoints d'authentification
- [ ] `X-Powered-By` supprimé

#### Infrastructure
- [ ] Endpoint `/health` (ou équivalent) opérationnel
- [ ] Rollback plan défini : version précédente identifiée, procédure documentée
- [ ] Logs applicatifs accessibles et configurés au bon niveau (`info` en prod, pas `debug`)
- [ ] Alertes (uptime, erreurs 5xx) configurées

### 3 — Rapport
Générer `docs/deploy-[env]-[DATE].md` :

```markdown
# Pré-déploiement [ENV] — [DATE]

## Résultat : ✅ AUTORISÉ / ❌ BLOQUÉ

| Domaine | Statut | Items en échec |
|---------|--------|----------------|
| Code & Build | ✅/❌ | — |
| Configuration | ✅/❌ | — |
| Base de données | ✅/❌ | — |
| Sécurité | ✅/❌ | — |
| Infrastructure | ✅/❌ | — |

## Items bloquants
[Liste des items critiques non validés]

## Actions requises avant déploiement
1. ...
```

### 4 — Décision
- **Un seul item critique en échec** → `❌ DÉPLOIEMENT BLOQUÉ` — lister les actions requises
- **Tous items critiques validés** → `✅ AUTORISÉ` — noter les warnings non bloquants
- Production : demander confirmation explicite de l'utilisateur même si tout est vert
