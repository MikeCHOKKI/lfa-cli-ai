---
name: svg-art
description: >
  Crée tout type de visuel SVG ou HTML de qualité production : diagrammes (flowchart,
  structurel, illustratif, ERD, séquence), logos et identités visuelles, art abstrait
  génératif, backgrounds et textures, mockups UI, icônes, illustrations, cartes,
  infographies, animations CSS, et widgets interactifs. Déclenche ce skill dès que
  l'utilisateur mentionne "logo", "design", "icône", "illustration", "background",
  "pattern", "art", "abstrait", "bannière", "affiche", "poster", "visuel", "SVG",
  "créer quelque chose de beau" — ou toute demande visuelle qui n'est pas un simple
  diagramme de données. Déclenche aussi pour les specs visuelles sans verbe explicite :
  "logo pour mon projet", "fond géométrique", "icône de paramètres". Ce skill contient
  les règles complètes de design system, coordonnées SVG, dark mode, accessibilité,
  typographie et composition pour produire des outputs irréprochables.
---

# SVG Art & Design Skill

Ce skill couvre **tout ce qui se dessine** : diagrammes, logos, art, backgrounds,
illustrations, mockups, icônes. Lis ce fichier en entier avant d'écrire la première balise.
Consulte ensuite le fichier de référence correspondant au type demandé.

---

## Arbre de décision — quel type de visuel ?

```
L'utilisateur veut...
│
├─ Un processus / des étapes              → Flowchart          [références/diagrams.md]
├─ Ce qui est dans quoi                   → Structural         [références/diagrams.md]
├─ Comment ça fonctionne                  → Illustratif        [références/diagrams.md]
├─ Un schéma de base de données           → ERD / mermaid.js   [références/diagrams.md]
├─ Des données à comparer / visualiser    → Chart              [références/diagrams.md]
├─ Un logo / marque / identité            → Logo               [références/logo-icon.md]
├─ Une icône / pictogramme                → Icon               [références/logo-icon.md]
├─ Un fond / texture / pattern            → Background         [références/art.md]
├─ De l'art abstrait / génératif          → Abstract art       [références/art.md]
├─ Une illustration / scène               → Illustration       [références/art.md]
├─ Une bannière / affiche / poster        → Poster             [références/art.md]
├─ Un formulaire / interface / mockup UI  → Mockup             [références/mockup-ui.md]
└─ Un widget interactif / calculateur     → Interactive        [références/interactive.md]
```

**Lire le fichier de référence du type détecté avant de produire quoi que ce soit.**
Pour les demandes hybrides (ex: logo + background), lire les deux fichiers concernés.

---

## Design system — règles universelles

Ces règles s'appliquent à **tous** les outputs sans exception.

### Typographie

- 14px : labels, titres de nœuds (classe `th` ou `t`)
- 12px : sous-titres, légendes, annotations (classe `ts`)
- Poids : 400 (normal) ou 500 (medium) — jamais 600 ou 700
- Sentence case partout — jamais TOUT MAJUSCULES ni Title Case dans les diagrammes
- Italique interdit dans les diagrammes ; autorisé dans les logos et posters

Pour les logos et l'art : typographie libre — voir `références/logo-icon.md`

### Palette de couleurs

Utilise **uniquement** les variables CSS et les classes de ramp pour tous les outputs
fonctionnels (diagrammes, mockups, widgets).

Pour l'art, logos et backgrounds : les hex codés en dur sont autorisés — voir règles
spécifiques dans `références/art.md` et `références/logo-icon.md`.

**Variables CSS du design system :**
```
--color-background-primary      blanc / fond principal
--color-background-secondary    surfaces
--color-background-tertiary     fond page
--color-text-primary            noir / texte principal
--color-text-secondary          texte atténué
--color-text-tertiary           hints
--color-border-tertiary         bordures légères 0.15α
--color-border-secondary        hover 0.3α
--color-border-primary          0.4α
--font-sans, --font-serif, --font-mono
--border-radius-md (8px), --border-radius-lg (12px), --border-radius-xl (16px)
```

**9 ramps × 7 stops (light/dark gérés automatiquement par les classes) :**

| Classe | 50 | 100 | 200 | 400 | 600 | 800 | 900 |
|---|---|---|---|---|---|---|---|
| `c-purple` | #EEEDFE | #CECBF6 | #AFA9EC | #7F77DD | #534AB7 | #3C3489 | #26215C |
| `c-teal`   | #E1F5EE | #9FE1CB | #5DCAA5 | #1D9E75 | #0F6E56 | #085041 | #04342C |
| `c-coral`  | #FAECE7 | #F5C4B3 | #F0997B | #D85A30 | #993C1D | #712B13 | #4A1B0C |
| `c-pink`   | #FBEAF0 | #F4C0D1 | #ED93B1 | #D4537E | #993556 | #72243E | #4B1528 |
| `c-gray`   | #F1EFE8 | #D3D1C7 | #B4B2A9 | #888780 | #5F5E5A | #444441 | #2C2C2A |
| `c-blue`   | #E6F1FB | #B5D4F4 | #85B7EB | #378ADD | #185FA5 | #0C447C | #042C53 |
| `c-green`  | #EAF3DE | #C0DD97 | #97C459 | #639922 | #3B6D11 | #27500A | #173404 |
| `c-amber`  | #FAEEDA | #FAC775 | #EF9F27 | #BA7517 | #854F0B | #633806 | #412402 |
| `c-red`    | #FCEBEB | #F7C1C1 | #F09595 | #E24B4A | #A32D2D | #791F1F | #501313 |

Mode light : stop 50 fill + 600 stroke + 800 titre / 600 sous-titre
Mode dark  : stop 800 fill + 200 stroke + 100 titre / 200 sous-titre

Ne mets jamais `c-{ramp}` sur un `<path>` — seulement sur `<rect>`, `<circle>`,
`<ellipse>`, `<polygon>`, ou `<g>`.

**Sémantique couleur pour les diagrammes :**
- Gris → structurel / neutre / start / end
- Purple, teal, coral, pink → catégories générales (max 3 ramps par diagramme)
- Blue → informationnel ; Green → succès ; Amber → warning ; Red → erreur

### Effets interdits (sauf exceptions explicites dans les fichiers de référence)

- Pas de drop shadow, blur, glow, neon
- Pas de gradients multiples (max 1 `<linearGradient>` dans les illustratifs physiques)
- Pas de fond coloré sur le conteneur racine SVG (toujours transparent)
- Pas de `position: fixed` en HTML
- Pas de `localStorage` ou `sessionStorage`
- Pas de commentaires HTML/CSS (`<!-- -->` et `/* */`)

---

## Setup SVG universel

### Structure de base
```svg
<svg width="100%" viewBox="0 0 680 H" role="img">
  <title>Description courte et précise</title>
  <desc>Description longue pour screen readers</desc>
  <!-- contenu -->
</svg>
```

- **Largeur 680 est fixe** — ne jamais changer
- Zone de sécurité : x=40 → x=640 (20px marge de chaque côté)
- `H` = y_max_element_bas + 40px de padding — calculer après avoir placé tous les éléments
- Zéro coordonnée négative
- Jamais deux `<svg>` racines dans le même output

Pour les logos et l'art : le viewBox peut être carré (0 0 400 400) ou en format
affiche (0 0 680 960) — voir chaque fichier de référence.

### Marker flèche (pour tous les diagrammes avec connexions)
```svg
<defs>
  <marker id="arrow" viewBox="0 0 10 10" refX="8" refY="5"
          markerWidth="6" markerHeight="6" orient="auto-start-reverse">
    <path d="M2 1L8 5L2 9" fill="none" stroke="context-stroke"
          stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
  </marker>
</defs>
```

### Classes pré-définies disponibles dans le widget SVG
```
t       sans 14px color-text-primary
ts      sans 12px color-text-secondary
th      sans 14px medium (500) color-text-primary
box     rect neutre (bg-secondary fill, border-secondary stroke)
node    groupe cliquable (hover, cursor pointer)
arr     ligne de connexion 1.5px
leader  ligne pointillée guide 0.5px tertiary
```

Chaque `<text>` dans un diagramme **doit** avoir une de ces classes.
Dans les logos et l'art, les `<text>` peuvent avoir leurs propres styles inline.

### Calcul de dimensions

**Texte → boîte :**
- Char moyen 14px ≈ 8px large ; char 12px ≈ 7px
- `width = max_string_chars × 8 + 24px padding`
- Exemple : "Authentication" (14 chars) → `14 × 8 + 24 = 136px`

**Placement vertical :**
Toujours `dominant-baseline="central"` avec y au centre vertical de la zone.
Sans ça, le texte dérive ~4px vers le haut.

**Packing horizontal — vérification obligatoire :**
```
N boîtes × largeur + (N-1) gaps ≤ 600px (entre x=40 et x=640)
Exemple : 4 × 130 + 3 × 20 = 580px → OK
          4 × 160 = 640px → TROP → réduire à 130px ou passer en 2 rangées
```

### Checklist avant de finaliser tout SVG
1. Élément le plus bas : `max(y + height)` + 40 = viewBox height
2. Aucun élément ne dépasse x=640
3. `text-anchor="end"` à x < 60 → risque de débord gauche → utiliser `start`
4. Aucun chevauchement de labels (les formes peuvent se chevaucher dans l'art)
5. Tout `<path>` connecteur a `fill="none"`
6. `role="img"` + `<title>` + `<desc>` présents

---

## Setup HTML interactif

```html
<style>
  /* variables du design system héritées automatiquement */
  /* max 15-20 lignes ici */
</style>

<div style="display:block;width:100%">
  <!-- contenu -->
</div>

<script>
  // s'exécute après le streaming
</script>
```

Règles :
- Pas de `<!DOCTYPE>`, `<html>`, `<head>`, `<body>`
- Librairies uniquement depuis : `cdnjs.cloudflare.com`, `esm.sh`, `cdn.jsdelivr.net`,
  `unpkg.com`, `fonts.googleapis.com`, `fonts.gstatic.com`
- Librairies JS dispo : recharts, d3, plotly, three (r128), papaparse, xlsx,
  mathjs, lodash, chart.js, tone, mammoth, tensorflow

**`sendPrompt(text)`** — fonction globale → envoie un message au chat Claude.
Utiliser pour les nœuds cliquables qui méritent une réponse approfondie.

---

## Accessibilité — non négociable

- SVG : `role="img"` + `<title>` (court) + `<desc>` (long) en premiers enfants
- Icônes décoratives : `aria-hidden="true"`
- Boutons icône seuls : `aria-label="..."`
- Ne jamais différencier uniquement par la couleur — ajouter forme, dash ou texture

---

## Dark mode — vérification mentale avant tout output

Question : "Si le fond était quasi-noir (#1a1a1a), chaque texte resterait-il lisible ?"

- Diagrammes/mockups SVG : textes dans `c-{ramp}` → automatique
- HTML : `var(--color-text-primary)` — jamais `color: #333` ou `color: black`
- Art / logos avec hex : tester la lisibilité sur fond sombre
- Ajouter `@media (prefers-color-scheme: dark)` si nécessaire pour les hex codés

---

## Fichiers de référence — lire avant de produire

| Type de demande | Fichier à lire |
|---|---|
| Flowchart, structurel, illustratif, ERD, séquence | `références/diagrams.md` |
| Logo, icône, emblème, badge, favicon | `références/logo-icon.md` |
| Art abstrait, génératif, background, texture, pattern, poster | `références/art.md` |
| Mockup UI, formulaire, dashboard, composant | `références/mockup-ui.md` |
| Widget interactif, calculateur, stepper, animation | `références/interactive.md` |
