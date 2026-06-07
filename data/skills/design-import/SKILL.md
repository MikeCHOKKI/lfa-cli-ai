---
name: design-import
description: >
  Transforme des designs Figma en code de production via l'API Figma (MCP) ou
  export JSON. Extrait les tokens de design (couleurs, typographie, espacements,
  ombres), génère des composants fidèles au pixel près, avec tous les états,
  responsive, dark mode, accessibilité, animations. Tech-agnostique :
  React, Vue, Svelte, Angular, HTML pur. Déclenche ce skill quand l'utilisateur
  mentionne "figma", "maquette design", "design to code",
  "importer design", "extraire design", "transformer maquette".
---

# Design Import Skill

Transforme des designs Figma en code de production.

Charge les skills suivants selon les besoins :
- `design-system` → tokens, palette, typographie
- `mockup-ui` → patterns UI, composants
- `animation` → micro-interactions, transitions
- `code-quality` → structure, nommage, patterns
- `project-standards` → setup projet, conventions
- `svg-art` → icônes, illustrations, graphics

---

## Pipeline d'import

```
1. ANALYSE      → Lire la source (Figma JSON via MCP / export / capture)
2. TOKENS       → Extraire couleurs, typo, spacing, shadows, radius
3. STRUCTURE    → Décomposer en composants, props, états
4. GÉNÉRATION   → Produire le code dans la techno cible
5. VÉRIFICATION → Comparer rendu vs design original
```

---

## 1. ANALYSE — Formats d'entrée

### Figma — via le MCP server (automatique, recommandé)

Un serveur MCP `figma-mcp-server` est intégré pour récupérer les designs sans export manuel.

**Prérequis :** Token Figma défini dans l'environnement (`FIGMA_TOKEN`) ou passé en paramètre.

**Utilisation directe par l'agent :**
```
$figma_get_file file_key="ABCDEF123" depth=4
```

**Workflow complet :**
1. L'agent appelle `figma_get_file` avec la file_key du design Figma
2. Le MCP fetch le JSON via l'API Figma REST et le filtre (garde seulement les champs utiles)
3. Le résultat filtré est passé à l'étape TOKENS pour extraction des tokens

**Autres outils MCP disponibles :**
| Outil | Description |
|---|---|
| `figma_get_file` | Récupère et filtre tout le fichier (nom, version, arbre) |
| `figma_get_node` | Récupère un nœud spécifique (utile pour cibler un composant) |
| `figma_get_images` | Génère des URLs d'export PNG/SVG/PDF |
| `figma_get_styles` | Liste les styles (couleurs, textes) du fichier |

### Figma — export JSON (manuel)

```json
{
  "document": {
    "children": [
      {
        "type": "FRAME",
        "name": "Button/Primary",
        "children": [ /* ... */ ],
        "backgroundColor": { "r": 0.31, "g": 0.43, "b": 0.97, "a": 1 },
        "cornerRadius": 8,
        "effects": [
          { "type": "DROP_SHADOW", "offset": { "x": 0, "y": 2 },
            "radius": 4, "color": { "r": 0, "g": 0, "b": 0, "a": 0.1 } }
        ]
      }
    ]
  }
}
```

Extraction clé :
| Figma | Code |
|-------|------|
| `backgroundColor` → fill | `background-color` / `bg` |
| `cornerRadius` | `border-radius` |
| `effects[].DROP_SHADOW` | `box-shadow` |
| `style.fontSize`, `style.fontFamily` | `font-size`, `font-family` |
| `style.fontWeight` | `font-weight` |
| `style.lineHeightPx` | `line-height` |
| `style.letterSpacing` | `letter-spacing` |
| `layoutMode`, `itemSpacing` | `display: flex`, `gap` |
| `paddingLeft/Right/Top/Bottom` | `padding` |
| `autoLayout` → contraintes | `flex` layout |

### Design Tokens (format structuré)

Si tu as un système de design tokens (Style Dictionary, W3C Design Tokens, etc.), utilise ce format :

```json
{
  "components": {
    "button_primary": {
      "properties": {
        "label": { "type": "string", "default": "Click me" },
        "variant": { "type": "enum", "options": ["primary", "secondary", "ghost"] },
        "disabled": { "type": "boolean", "default": false },
        "size": { "type": "enum", "options": ["sm", "md", "lg"], "default": "md" }
      },
      "styles": {
        "base": {
          "backgroundColor": "#4F6EF7",
          "borderRadius": "8px",
          "padding": "10px 20px",
          "fontSize": "15px",
          "fontWeight": "500",
          "color": "#FFFFFF"
        },
        "hover": { "backgroundColor": "#3B57D9" },
        "disabled": { "opacity": 0.5, "cursor": "not-allowed" }
      },
      "states": ["default", "hover", "focus", "active", "disabled"],
      "responsive": {
        "sm": { "padding": "8px 16px", "fontSize": "14px" },
        "lg": { "padding": "14px 28px", "fontSize": "17px" }
      }
    }
  },
  "designTokens": {
    "colors": { "primary-500": "#4F6EF7", "primary-600": "#3B57D9" },
    "typography": { "font-sans": "Inter, sans-serif" },
    "spacing": { "sm": "8px", "md": "16px", "lg": "24px" },
    "shadows": { "sm": "0 1px 2px rgba(0,0,0,0.05)" },
    "radius": { "sm": "4px", "md": "8px", "lg": "12px" }
  }
}
```

### Design via capture ou description

Si pas de JSON structuré, extraire manuellement :
```
1. Palette → lister toutes les couleurs utilisées (fonds, textes, bordures)
2. Typographie → styles de texte (font, taille, poids, hauteur)
3. Espacements → padding, margin, gap entre éléments
4. Ombres → offset, blur, spread, couleur
5. Coins → border-radius sur chaque composant
6. Hiérarchie → écran → section → composant → sous-composant
```

---

## 2. TOKENS — Extraction et normalisation

Charger le skill `design-system` pour la structure des tokens.

### Mapping couleur Figma → CSS

```python
# Figma RGB (0-1) → Hex
def figma_rgb_to_hex(r, g, b):
    return f"#{round(r*255):02x}{round(g*255):02x}{round(b*255):02x}"

# Figma opacity → alpha channel
def figma_opacity(a):
    return round(a * 100) / 100  # arrondi à 2 décimales
```

### Détection de la palette

```javascript
// Algorithme : grouper les couleurs par distance < 10%,
// assigner des rôles sémantiques
const colors = extractAllColors(figmaJson)
const palette = {
  primary: dominant(colors.filter(c => isHueInRange(c, 200, 260))),
  secondary: dominant(colors.filter(c => isHueInRange(c, 0, 60))),
  neutral: colors.filter(c => saturation(c) < 10%),
  semantic: {
    success: dominant(colors.filter(c => isHueInRange(c, 120, 150))),
    warning: dominant(colors.filter(c => isHueInRange(c, 35, 55))),
    danger: dominant(colors.filter(c => isHueInRange(c, 0, 15)))
  }
}
```

### Génération des stops

À partir de la couleur primaire détectée :
```javascript
function generateStops(hex, count = 9) {
  // Interpolation HSL entre blanc (50) et noir (900)
  // stop 500 = couleur d'origine
  const stops = {}
  const h = hexToHsl(hex).h
  for (let i = 0; i < count; i++) {
    const lightness = 95 - i * 10  // 95%, 85%, ... 15%
    stops[`${(i+1)*100}`] = hslToHex(h, saturationForStop(i), lightness)
  }
  return stops
}
```

---

## 3. STRUCTURE — Décomposition en composants

### Règles de découpage

```yaml
Écran:
  Page / Route:
    Header:
      Logo (svg-art si vectoriel)
      Navigation
      Avatar
    Main:
      Sidebar (si présent)
      Content:
        Card[]
          CardHeader
          CardBody
          CardFooter
        DataTable
        Form:
          FormField[]
    Footer
```

### États par composant

```
default    → état normal
hover      → survol souris
focus      → focus clavier
active     → clic / pressed
disabled   → désactivé
loading    → chargement (skeleton / spinner)
error      → erreur (shake, bordure rouge)
empty      → aucune donnée (illustration + message)
selected   → sélectionné (tabs, listes)
checked    → coché (checkbox, radio)
```

### Props génériques

```typescript
interface ComponentProps {
  className?: string
  style?: React.CSSProperties
  children?: React.ReactNode
  // États
  disabled?: boolean
  loading?: boolean
  error?: string | null
  // Accessibilité
  ariaLabel?: string
  role?: string
  tabIndex?: number
}
```

---

## 4. GÉNÉRATION — Production du code

### Stack cible

Détecter ou demander la techno. Générique :

```typescript
// React + TypeScript + CSS Modules (stack par défaut)
// Alternative : Vue SFC, Svelte, Angular, Lit, Web Components
```

### Structure de fichier par composant

```
components/
  Button/
    Button.tsx          → composant
    Button.module.css   → styles
    Button.test.tsx     → tests
    Button.stories.tsx  → storybook
    index.ts            → export
```

### Template de composant

```typescript
import { type ButtonHTMLAttributes, type ReactNode } from 'react'
import styles from './Button.module.css'

type ButtonVariant = 'primary' | 'secondary' | 'ghost'
type ButtonSize = 'sm' | 'md' | 'lg'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
  icon?: ReactNode
}

export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  icon,
  children,
  disabled,
  className = '',
  ...props
}: ButtonProps) {
  const classNames = [
    styles.button,
    styles[variant],
    styles[size],
    loading ? styles.loading : '',
    className
  ].filter(Boolean).join(' ')

  return (
    <button
      className={classNames}
      disabled={disabled || loading}
      {...props}
    >
      {loading && <span className={styles.spinner} aria-hidden />}
      {icon && <span className={styles.icon}>{icon}</span>}
      {children && <span className={styles.label}>{children}</span>}
    </button>
  )
}
```

### Styles (CSS Modules)

```css
/* Tokens du design importé */
.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: none;
  cursor: pointer;
  font-family: var(--font-sans);
  font-weight: 500;
  border-radius: var(--radius-md);
  transition: all 0.15s var(--ease-smooth);
  /* Design tokens extraits */
  &.primary {
    background: var(--color-primary-500);
    color: #fff;
    &:hover { background: var(--color-primary-600); }
    &:active { transform: scale(0.98); }
  }
  &.secondary {
    background: transparent;
    color: var(--color-primary-500);
    border: 1px solid var(--color-primary-500);
    &:hover { background: var(--color-primary-50); }
  }
  &.ghost {
    background: transparent;
    color: var(--color-text-secondary);
    &:hover { background: var(--color-background-secondary); }
  }
  &.sm  { padding: 6px 14px; font-size: 13px; }
  &.md  { padding: 10px 20px; font-size: 15px; }
  &.lg  { padding: 14px 28px; font-size: 17px; }
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
  }
}
```

---

## 5. VÉRIFICATION — Fidélité au design

### Checklist pixel-perfect

```
[ ] Couleurs : chaque hex correspond à la maquette (±2%)
[ ] Typographie : font, taille, poids, line-height identiques
[ ] Espacements : padding, margin, gap identiques (±2px)
[ ] Ombres : offset, blur, couleur identiques
[ ] Bordures : épaisseur, couleur, radius identiques
[ ] États : hover, focus, active, disabled présents
[ ] Icônes : taille, stroke, fill, alignement identiques
[ ] Layout : position, taille, alignement respectés
[ ] Responsive : breakpoints de la maquette suivis
[ ] Dark mode : si présent dans la maquette
[ ] Accessibilité : labels, rôles, focus visible, contrastes
```

### Ratios de vérification

- Taille des textes : ±0px (pixel-perfect)
- Espacements : ±2px max
- Couleurs : même hex (tolérance 2% si approximation Figma)
- Ombres : même paramètres

---

## Références

- [Figma REST API](https://www.figma.com/developers/api)
- [Figma Personal Access Tokens](https://help.figma.com/hc/en-us/articles/8085703771159-Manage-personal-access-tokens)
- Serveur MCP Figma : `figma-mcp-server/` (dans ce dossier)
- Design tokens → charger `design-system`
- UI patterns → charger `mockup-ui`
- Animations → charger `animation`
- Qualité code → charger `code-quality`
