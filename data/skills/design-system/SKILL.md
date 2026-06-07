---
name: design-system
description: >
  Génère des systèmes de design complets : palettes de couleurs, échelles typographiques,
  espacements, ombres, design tokens, modes clair/sombre, accessibilité WCAG.
  Déclenche ce skill quand l'utilisateur demande "palette", "couleurs", "typographie",
  "design tokens", "design system", "mode sombre", "accessibilité", "contraste", "font",
  ou toute demande de création ou d'extension d'un système de design cohérent.
---

# Design System Skill

Génération de systèmes de design cohérents, accessibles et réutilisables.

---

## Palette de couleurs

### Génération de palette

Utiliser l'espace HSL pour des harmonies cohérentes :

```
Teinte (H) : 0-360
Saturation (S) : 5-100%
Luminosité (L) : 5-95%
```

### Types d'harmonie

| Type | Règle | Utilisation |
|------|-------|-------------|
| Monochromatique | Même H, varier S/L | UI minimaliste, data viz simple |
| Analogue | H ±30° | Interfaces générales, marques |
| Complémentaire | H +180° | Accents, CTAs, alerts |
| Triadique | H ±120° | Palettes riches, dashboards |
| Split-complémentaire | H +150° +210° | UI équilibrée avec accent |
| Quadratique | H ±90° ±180°| Palettes complexes, charts |

### Structure de palette recommandée

```
c-primary    → Teinte principale de la marque (H₀)
c-secondary  → Teinte secondaire (H₀ ± 30-60)
c-accent     → Teinte d'accent (H₀ + 150-180)
c-neutral    → Gris (H₀, S<10%)
c-success    → Vert (H≈120-140)
c-warning    → Ambre (H≈35-45)
c-danger     → Rouge (H≈0-10)
c-info       → Bleu (H≈200-220)
```

### Stops (9 par couleur)

| Stop | Usage light | Usage dark |
|------|-------------|------------|
| 50   | Fond léger | Non utilisé |
| 100  | Fond hover | Fond foncé |
| 200  | Bordure    | Fond hover |
| 300  | Bordure forte | Accent secondaire |
| 400  | Texte secondaire | Texte principal |
| 500  | Texte principal (MD) | Texte principal (light) |
| 600  | Texte principal | Texte secondaire |
| 700  | Titres | Bordure |
| 800  | Titres forts | Fond |

### Accessibilité des contrastes WCAG 2.1 AA

```javascript
// Ratio minimum : 4.5:1 pour texte normal, 3:1 pour grand texte (18px+ bold ou 24px+)
// Vérifier : L1 + 0.05 / L2 + 0.05 ≥ seuil
// L = 0.2126R + 0.7152G + 0.0722B (linéarisé)
function luminance(hex) {
  const [r, g, b] = hex.match(/../g).map(c => {
    const v = parseInt(c, 16) / 255
    return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4)
  })
  return 0.2126 * r + 0.7152 * g + 0.0722 * b
}
function contrastRatio(hex1, hex2) {
  const l1 = luminance(hex1), l2 = luminance(hex2)
  return (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05)
}
```

### Combinaisons fiables (light mode)
- Texte #1a1a1a sur fond #ffffff → 18.0:1
- Texte #4a4a4a sur fond #ffffff → 8.6:1
- Texte #ffffff sur fond #1a1a1a → 18.0:1
- Texte #ffffff sur primary-600 → toujours vérifier

---

## Typographie

### Échelle typographique modulaire

Base : 16px, Ratio : 1.25 (Major Second) ou 1.333 (Perfect Fourth)

```
Nom           Taille     Ratio 1.25  Ratio 1.333  Poids
caption       12px       ---         ---          400
body-small    14px       0.875       0.875        400
body          16px       1           1            400
body-large    18px       1.125       1.125        400
h6            20px       1.25        1.25         500
h5            24px       1.5         1.333        500
h4            28px       1.75        1.424        600
h3            32px       2           1.602        600
h2            40px       2.5         1.802        600
h1            48px       3           2.027        600
display       64px       4           2.441        700
```

### Font stacks

```css
/* Système (UI) */
--font-sans:  'Inter', 'SF Pro Display', -apple-system, BlinkMacSystemFont, 
              'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif

/* Titres */
--font-heading: 'Plus Jakarta Sans', 'Inter', --font-sans

/* Data, code */
--font-mono:  'JetBrains Mono', 'SF Mono', 'Fira Code', 
              'Cascadia Code', Consolas, monospace

/* Long format */
--font-serif: 'Merriweather', 'Georgia', 'Times New Roman', serif
```

### Line heights
```
caption/body-small   → 1.4
body/body-large      → 1.5
h6/h5/h4             → 1.3
h3/h2/h1             → 1.2
display              → 1.1
```

### Letter spacing
```
caption    → 0.02em
body       → 0
h6 → h4   → -0.01em
h3 → h1   → -0.02em
display    → -0.03em
```

---

## Design tokens

### Format universel

```json
{
  "color": {
    "primary": { "50": "#...", "100": "#...", /* ... */ "800": "#..." },
    "secondary": { /* ... */ },
    "neutral": { /* ... */ },
    "semantic": {
      "success": { "50": "#...", "500": "#..." },
      "warning": { "50": "#...", "500": "#..." },
      "danger": { "50": "#...", "500": "#..." },
      "info": { "50": "#...", "500": "#..." }
    }
  },
  "typography": {
    "fontFamily": { "sans": "...", "heading": "...", "mono": "...", "serif": "..." },
    "fontSize": { "caption": "12px", "body": "16px", /* ... */ "display": "64px" },
    "lineHeight": { "tight": "1.2", "normal": "1.5", "relaxed": "1.75" },
    "fontWeight": { "normal": "400", "medium": "500", "semibold": "600", "bold": "700" }
  },
  "spacing": {
    "xs": "4px", "sm": "8px", "md": "16px", "lg": "24px", 
    "xl": "32px", "2xl": "48px", "3xl": "64px"
  },
  "borderRadius": {
    "sm": "4px", "md": "8px", "lg": "12px", "xl": "16px", "full": "9999px"
  },
  "shadow": {
    "sm": "0 1px 2px rgba(0,0,0,0.05)",
    "md": "0 4px 6px rgba(0,0,0,0.07)",
    "lg": "0 10px 15px rgba(0,0,0,0.1)",
    "xl": "0 20px 25px rgba(0,0,0,0.12)"
  }
}
```

### CSS Custom Properties

```css
:root {
  --color-primary-50: #...;
  --color-primary-100: #...;
  /* ... */
  --font-sans: 'Inter', ...;
  --font-heading: 'Plus Jakarta Sans', ...;
  --font-mono: 'JetBrains Mono', ...;
  --space-xs: 4px;
  --space-sm: 8px;
  /* ... */
  --radius-sm: 4px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
}

@media (prefers-color-scheme: dark) {
  :root {
    /* Inverser les stops : fond utilise 800, texte utilise 100 */
  }
}
```

---

## Mode sombre

Règles de transfert light → dark :

| Token light | Token dark |
|-------------|------------|
| `primary-50` (fond) | `primary-800` (fond) |
| `primary-100` (hover) | `primary-700` (hover) |
| `primary-600` (texte) | `primary-200` (texte) |
| `neutral-50` (fond page) | `neutral-900` (fond page) |
| `neutral-100` (carte) | `neutral-800` (carte) |
| `neutral-800` (titre) | `neutral-100` (titre) |
| shadow: rgba(0,0,0,α) | shadow: rgba(0,0,0,α×2) |

---

## Vérification

1. Palette complète : 9 rampes × 9 stops ?
2. Rapport de contraste ≥ 4.5:1 pour tous les textes
3. Échelle typographique cohérente (ratio)
4. Design tokens couvrent : couleur, typo, spacing, radius, shadow
5. Mode sombre défini pour chaque token
6. Font stacks avec fallbacks
