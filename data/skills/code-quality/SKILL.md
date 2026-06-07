---
name: code-quality
description: >
  Applique les standards de qualité de code dans tout langage : clean code,
  principes SOLID, design patterns, code review, nommage, structure de fichiers,
  anti-patterns, refactoring. Déclenche ce skill quand l'utilisateur demande
  "qualité", "clean code", "refactor", "pattern", "SOLID", "DRY", "KISS",
  "code review", "lint", "anti-pattern", "bonnes pratiques", "architecture code",
  ou toute demande d'amélioration structurée du code existant.
---

# Code Quality Skill

Standards de qualité projets et langages.

---

## Principes universels (tout langage)

### SOLID
| Principe | Règle | Violation typique |
|----------|-------|-------------------|
| SRP | Une classe/fonction = une responsabilité | God class, fonction de 200+ lignes |
| OCP | Extensible sans modification | Switch/cas sur type |
| LSP | Sous-type substitue le type | Héritage qui brise le contrat |
| ISP | Interfaces fines et spécifiques | Interface fourre-tout |
| DIP | Modules haut niveau ne dépendent pas des bas niveau | New() en dur dans une classe métier |

### GRASP
- Information Expert : placer la responsabilité là où est l'info
- Controller : couche intermédiaire entre UI et domaine
- Low Coupling / High Cohesion : modules indépendants, cohérents en interne
- Polymorphism : remplacer les conditions par du polymorphisme

### Nommage
```
/-- Mauvais --/                /-- Bon --/
a, b, x                       index, total, userCount
data, info, stuff, manager    invoice, sessionStore, paymentHandler
getData()                     fetchUserProfile()
handleClick()                 submitLoginForm()
processStuff()                normalizePhoneNumber()
isStatus                      isValid, hasAccess, containsDuplicates
```

### Règles de fonction
- < 20 lignes idéal, 40 max
- 0-2 paramètres idéal, 3 acceptable, 4+ → objet paramètre
- Un seul niveau d'abstraction par fonction
- Pas de `else` si possible (early return)
- Pas de flag booléen en paramètre → diviser en 2 fonctions

### Anti-patterns courants
| Anti-pattern | Solution |
|---|---|
| God Class | Découper par responsabilité |
| Spaghetti Code | Architecture en couches |
| Golden Hammer | Choisir l'outil adapté au problème |
| Premature Optimization | Écrire clair d'abord, profiler ensuite |
| Copy-Pasta | Factoriser, paramétrer |
| Shotgun Surgery | Regrouper les responsabilités qui changent ensemble |

---

## Par langage

### JavaScript / TypeScript
```typescript
// Préférer const > let > var
const users = await fetchUsers()

// Destructuring
const { name, email } = user

// Optional chaining
const city = user?.address?.city

// Default params
function createUser({ name = 'Anonymous', email = '' } = {}) {}

// Pas de any en TS → unknown ou type explicite
type User = { id: string; name: string; email: string }

// Pas de magic strings
const STATUS = { ACTIVE: 'active', INACTIVE: 'inactive' } as const

// Async/await > .then()
async function loadData() {
  try {
    const data = await api.fetch()
    return data
  } catch (error) {
    logger.error('Fetch failed', error)
    throw error
  }
}
```

### Python
```python
# Type hints obligatoires
def calculate_total(items: list[Item], discount: float = 0.0) -> Decimal: ...

# Context managers
with open('file.txt') as f:
    content = f.read()

# Enum au lieu de strings
class Status(Enum):
    ACTIVE = auto()
    INACTIVE = auto()

# Pas de except nu
try:
    result = risky_operation()
except ValueError as e:
    logger.warning(f"Invalid input: {e}")
except DatabaseError as e:
    logger.error(f"DB error: {e}")
    raise
```

### Go
```go
// Erreurs toujours gérées
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething: %w", err)
}

// Interface minimales
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Constructeur vs initialisation directe
func NewUser(name string) *User {
    return &User{Name: name, createdAt: time.Now()}
}
```

---

## Code Review — checklist

```
[ ] Nommage clair et cohérent (casse, préfixes, verbes)
[ ] Fonctions < 20 lignes (40 max)
[ ] Pas de duplication détectable
[ ] Tests présents et significatifs
[ ] Gestion d'erreurs explicite
[ ] Sécurité : pas de secrets, pas de injection XSS/SQL
[ ] Pas de commentaires (le code s'explique)
[ ] Typage explicite (TypeScript/Python/Go)
[ ] Side-effects documentés
[ ] Performance : pas de N+1, pas de boucle dans boucle inutile
[ ] Accessibilité (pour le frontend)
[ ] Internationalisation (tous les textes utilisateur)
```

---

## Tests

### Structure AAA
```typescript
// Arrange
const user = new User({ name: 'Test' })

// Act
const result = user.isValid()

// Assert
expect(result).toBe(true)
```

### Pyramid de tests
```
E2E        ⬆ 10%
Integration ⬇ 20%
Unitaires   ⬇ 70%
```

### Règles
- Un test = un comportement
- Nommer comme une phrase : `should_return_true_when_user_is_valid`
- Tester les cas limites (empty, null, max, edge)
- Mocker les appels externes (API, DB, filesystem)
- Pas de logique conditionnelle dans les tests

---

## Git & Commits

### Conventional Commits
```
feat(scope): Ajout du module de paiement
fix(api): Correction timeout sur /products
refactor(auth): Extraction du middleware JWT
test(cart): Ajout tests panier vide
docs(api): Mise à jour endpoints README
chore(deps): Mise à jour lodash
```

### Branches
```
main / master     → production
develop           → intégration
feat/description  → nouvelle fonctionnalité
fix/description   → correction
refactor/desc     → refactoring
docs/desc         → documentation
```

---

## Architecture

### Layers (tout langage)
```
📁 presentation/   → UI, API routes, controllers
📁 application/    → use cases, orchestration (ports)
📁 domain/         → entités, règles métier (coeurs)
📁 infrastructure/ → DB, API externes, filesystem (adapters)
```

### Dependency Rule
```
presentation → application ← infrastructure
                    ↓
                domain (indépendant)
```

L'infrastructure implémente les ports définis dans application.
Le domaine ne dépend de rien.
