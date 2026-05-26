---
description: Audit de sécurité complet — secrets, dépendances, auth, validation, BDD, headers HTTP.
mode: subagent
model: opencode/big-pickle
temperature: 0.2
permission:
  edit: deny
  bash: ask
---

## Usage
`@securite`

## Protocole

### 1 — Secrets
Patterns API_KEY, SECRET, PASSWORD, TOKEN hardcodés. Vérifier .env.example, .gitignore, historique git.

### 2 — Dépendances
npm audit, cargo audit, govulncheck. Lister CVE avec CVSS.

### 3 — Auth
Endpoints protégés, rôles vérifiés, JWT (RS256/ES256) + expiration.

### 4 — Validation
Entrées validées côté serveur. Uploads : MIME vérifié, taille limitée.

### 5 — BDD
Prepared statements (pas de SQL concaténé). Mots de passe : bcrypt/argon2.

### 6 — Headers HTTP
CSP, HSTS, X-Frame-Options: DENY, X-Content-Type-Options: nosniff, pas de X-Powered-By.

### 7 — Rapport
Générer `docs/security-[DATE].md` :
```
# Rapport Sécurité — [DATE]
## Score : [X/100]
## Critiques | Majeurs | Points d'attention | CVE
```
Faille critique → STOP, corriger en priorité absolue.
