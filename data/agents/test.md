---
description: Génère des tests unitaires, d'intégration et E2E pour un module — couverture complète des cas critiques.
mode: subagent

temperature: 0.3
permission:
  edit: ask
  bash: ask
---

## Usage
`@test [module]`

---

## Protocole

### 1 — Analyse
- Lire tous les fichiers source du module
- Lister les méthodes publiques et les cas déjà testés (lire les tests existants s'ils existent)
- Identifier les dépendances à mocker (BDD, services externes, filesystem, time)
- Identifier les cas critiques non couverts

### 2 — Stratégie de test

| Type | Quand | Ce qu'on mock |
|------|-------|---------------|
| Unitaire | Logique métier isolée | Tout sauf l'unité testée |
| Intégration | Flux service → repository → BDD | Services externes uniquement |
| E2E | Parcours utilisateur critique | Rien (environnement réel ou proche) |

Privilégier les tests unitaires et d'intégration. E2E uniquement pour les flux métier prioritaires (paiement, auth, etc.).

### 3 — Structure AAA
Chaque test suit strictement : **Arrange → Act → Assert**

```typescript
// TypeScript / Jest
describe('UserService', () => {
  describe('createUser', () => {
    it('devrait créer un utilisateur avec un email valide', async () => {
      // Arrange
      const repo = mockUserRepository({ save: jest.fn().mockResolvedValue(fakeUser) })
      const sut = new UserService(repo)
      // Act
      const result = await sut.createUser({ email: 'test@example.com', name: 'Test' })
      // Assert
      expect(result.id).toBeDefined()
      expect(repo.save).toHaveBeenCalledOnce()
    })
  })
})
```

Adapter la syntaxe à la stack du projet (Go `testing`, PHPUnit, JUnit, pytest, etc.).

### 4 — Cas obligatoires par fonction

| Catégorie | Cas à couvrir |
|-----------|---------------|
| Nominal | Entrée valide → résultat attendu |
| Invalide | Entrée incorrecte → erreur/exception attendue |
| Limites | Valeur nulle, vide, 0, max, chaîne vide |
| Erreur externe | Mock qui échoue (BDD down, API timeout) |

**Endpoints HTTP obligatoirement** :
- `200`/`201` : cas nominal
- `400` : validation échouée
- `401`/`403` : non authentifié / non autorisé
- `404` : ressource inexistante
- `500` : erreur interne (mock qui throw)

### 5 — Mocks
- Repositories / BDD : mock avec données réalistes (pas `id: 1, name: "test"`)
- Services externes : mock au niveau du client HTTP, pas du service entier
- Time / Date : toujours injecter ou mocker — jamais `new Date()` dans le code métier testé
- Filesystem : mock ou dossier temporaire dédié aux tests

### 6 — Post-génération
- Lancer les tests : `npm test [module]` ou équivalent
- Afficher la couverture du module (lignes, branches)
- Signaler les cas non testés restants si la couverture < 80%
- Commit : `test([module]): Tests unitaires [module]`
