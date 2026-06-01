# Référence : Art, Backgrounds, Illustrations, Posters

Couvre : art abstrait, génératif, backgrounds, textures, patterns, illustrations,
affiches, bannières, scènes.

---

## Philosophie de l'art SVG

L'art SVG n'est pas un diagramme. Les règles changent :
- Hex codés en dur **autorisés** — liberté de palette totale
- Gradients **autorisés** — autant qu'il en faut
- Formes qui se chevauchent **autorisé** — c'est souvent l'effet voulu
- Labels de texte **optionnels** — parfois zero texte
- Pas de viewBox fixe 680 — adapter au format de la pièce

Mais les interdits restent :
- Pas de `filter: blur()` ni `drop-shadow`
- Pas de glow neon
- Pas de bitmaps intégrés (pas de `<image href="data:..."/>`)

---

## Formats et viewBox pour l'art

| Format           | ViewBox         | Usage                          |
|------------------|-----------------|--------------------------------|
| Carré            | `0 0 600 600`   | Avatar, vignette, pattern tile |
| Bannière web     | `0 0 1200 400`  | Header, hero background        |
| Affiche portrait | `0 0 600 900`   | Poster, flyer                  |
| Affiche paysage  | `0 0 900 600`   | Slide, cover                   |
| Wallpaper        | `0 0 1920 1080` | Background écran               |
| Carte            | `0 0 680 400`   | Intégré dans une page          |

Toujours `width="100%"` — le SVG s'adapte à son conteneur.

---

## Backgrounds et patterns

### Pattern géométrique (tuile répétée)

```svg
<svg width="100%" viewBox="0 0 600 600" role="img">
  <title>Background géométrique</title>
  <defs>
    <!-- Tuile 60×60 px -->
    <pattern id="tile" x="0" y="0" width="60" height="60" patternUnits="userSpaceOnUse">
      <!-- Fond de tuile -->
      <rect width="60" height="60" fill="#EEEDFE"/>
      <!-- Élément de la tuile -->
      <path d="M0 30 L30 0 L60 30 L30 60 Z"
            fill="none" stroke="#AFA9EC" stroke-width="0.8"/>
      <circle cx="30" cy="30" r="3" fill="#534AB7"/>
    </pattern>
  </defs>
  <!-- Appliquer le pattern -->
  <rect width="600" height="600" fill="url(#tile)"/>
</svg>
```

### Types de patterns géométriques

**Grille de points :**
```svg
<pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse">
  <circle cx="10" cy="10" r="1.5" fill="#888780"/>
</pattern>
```

**Grille de lignes :**
```svg
<pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse">
  <path d="M 40 0 L 0 0 0 40" fill="none" stroke="#D3D1C7" stroke-width="0.5"/>
</pattern>
```

**Hexagones :**
```svg
<pattern id="hex" width="52" height="30" patternUnits="userSpaceOnUse">
  <polygon points="26,0 52,15 52,15 26,30 0,15 0,15"
           fill="none" stroke="#AFA9EC" stroke-width="0.8"/>
  <polygon points="26,30 52,45 52,45 26,60 0,45 0,45"
           fill="none" stroke="#AFA9EC" stroke-width="0.8"/>
</pattern>
```

**Chevrons / zigzag :**
```svg
<pattern id="chevron" width="40" height="20" patternUnits="userSpaceOnUse">
  <polyline points="0,10 20,0 40,10" fill="none" stroke="#5DCAA5" stroke-width="1.5"/>
  <polyline points="0,20 20,10 40,20" fill="none" stroke="#5DCAA5" stroke-width="1.5"/>
</pattern>
```

**Diagonales hachurées :**
```svg
<pattern id="hatch" width="8" height="8" patternTransform="rotate(45)" patternUnits="userSpaceOnUse">
  <line x1="0" y1="0" x2="0" y2="8" stroke="#B4B2A9" stroke-width="1.5"/>
</pattern>
```

---

## Art abstrait génératif

### Composition en couches

Construire en 3-5 couches du fond vers l'avant :
1. **Fond** — couleur unie ou dégradé doux
2. **Formes larges** — grandes surfaces de couleur à faible opacité (0.2–0.5)
3. **Formes moyennes** — éléments structurants à opacité 0.6–0.8
4. **Accents** — petites formes vives, opacité 1
5. **Texture** — pattern en overlay à opacité 0.05–0.15

```svg
<svg width="100%" viewBox="0 0 600 600">
  <!-- Couche 1 : fond -->
  <rect width="600" height="600" fill="#1a1a2e"/>

  <!-- Couche 2 : grandes formes -->
  <ellipse cx="200" cy="200" rx="250" ry="200" fill="#534AB7" opacity="0.3"/>
  <ellipse cx="450" cy="400" rx="200" ry="250" fill="#1D9E75" opacity="0.25"/>

  <!-- Couche 3 : formes moyennes -->
  <circle cx="350" cy="250" r="120" fill="#D85A30" opacity="0.6"/>

  <!-- Couche 4 : accents -->
  <circle cx="180" cy="420" r="40" fill="#FAC775"/>
  <circle cx="520" cy="150" r="25" fill="#5DCAA5"/>

  <!-- Couche 5 : texture -->
  <rect width="600" height="600" fill="url(#noise-pattern)" opacity="0.08"/>
</svg>
```

### Formes organiques expressives

```svg
<!-- Blob (forme organique) -->
<path d="M 300 50
         C 450 30, 560 120, 550 250
         C 540 380, 430 470, 300 480
         C 170 490, 60 400, 50 270
         C 40 140, 150 70, 300 50 Z"
      fill="#534AB7" opacity="0.8"/>

<!-- Spirale -->
<path d="M 300 300
         Q 340 260, 360 300
         Q 380 340, 340 380
         Q 300 420, 240 380
         Q 180 340, 180 280
         Q 180 220, 240 200"
      fill="none" stroke="#1D9E75" stroke-width="3" stroke-linecap="round"/>
```

### Lignes de flux (flow lines)

```svg
<g fill="none" stroke="#AFA9EC" stroke-width="0.5" opacity="0.6">
  <path d="M 0 100 Q 150 80, 300 120 Q 450 160, 600 140"/>
  <path d="M 0 130 Q 150 110, 300 150 Q 450 190, 600 170"/>
  <path d="M 0 160 Q 150 140, 300 180 Q 450 220, 600 200"/>
  <path d="M 0 190 Q 150 170, 300 210 Q 450 250, 600 230"/>
</g>
```

### Scatter de formes

```svg
<!-- Générer manuellement des positions pseudo-aléatoires cohérentes -->
<!-- Utiliser une disposition en constellation, pas vraiment aléatoire -->
<g fill="#FAC775">
  <circle cx="80"  cy="120" r="3"/>
  <circle cx="145" cy="85"  r="2"/>
  <circle cx="210" cy="160" r="4"/>
  <circle cx="290" cy="95"  r="2"/>
  <circle cx="380" cy="140" r="3"/>
  <circle cx="450" cy="75"  r="5"/>
  <circle cx="520" cy="120" r="2"/>
</g>
```

---

## Dégradés SVG (autorisés dans l'art)

### Dégradé radial (ambiance lumière)
```svg
<defs>
  <radialGradient id="glow" cx="50%" cy="50%" r="50%">
    <stop offset="0%"   stop-color="#7F77DD" stop-opacity="0.8"/>
    <stop offset="100%" stop-color="#26215C" stop-opacity="0"/>
  </radialGradient>
</defs>
<ellipse cx="300" cy="300" rx="300" ry="300" fill="url(#glow)"/>
```

### Dégradé linéaire (coucher de soleil, horizon)
```svg
<defs>
  <linearGradient id="sunset" x1="0" y1="0" x2="0" y2="1">
    <stop offset="0%"   stop-color="#26215C"/>
    <stop offset="40%"  stop-color="#993C1D"/>
    <stop offset="70%"  stop-color="#D85A30"/>
    <stop offset="100%" stop-color="#FAC775"/>
  </linearGradient>
</defs>
<rect width="600" height="400" fill="url(#sunset)"/>
```

### Dégradé mesh (simulation)
Superposer plusieurs radialGradient à des positions différentes :
```svg
<defs>
  <radialGradient id="m1" cx="30%" cy="30%">
    <stop offset="0%" stop-color="#534AB7" stop-opacity="0.7"/>
    <stop offset="100%" stop-color="#534AB7" stop-opacity="0"/>
  </radialGradient>
  <radialGradient id="m2" cx="70%" cy="70%">
    <stop offset="0%" stop-color="#1D9E75" stop-opacity="0.7"/>
    <stop offset="100%" stop-color="#1D9E75" stop-opacity="0"/>
  </radialGradient>
</defs>
<rect width="600" height="600" fill="#0f0f23"/>
<rect width="600" height="600" fill="url(#m1)"/>
<rect width="600" height="600" fill="url(#m2)"/>
```

---

## Illustrations

### Illustration de scène

Construire par plans (comme un diorama) :
1. Fond de ciel / espace
2. Plans lointains (montagnes, horizon)
3. Plans moyens (bâtiments, arbres)
4. Premier plan (personnages, objets détaillés)
5. Overlay (brume, grain, vignette)

```svg
<!-- Paysage simplifié -->
<svg width="100%" viewBox="0 0 680 400">
  <!-- Ciel -->
  <defs>
    <linearGradient id="sky" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#042C53"/>
      <stop offset="100%" stop-color="#185FA5"/>
    </linearGradient>
  </defs>
  <rect width="680" height="400" fill="url(#sky)"/>

  <!-- Montagnes lointaines -->
  <polygon points="0,280 100,180 200,260 320,150 440,240 560,170 680,250 680,400 0,400"
           fill="#0C447C" opacity="0.7"/>

  <!-- Collines proches -->
  <ellipse cx="150" cy="380" rx="200" ry="60" fill="#0F6E56"/>
  <ellipse cx="530" cy="390" rx="220" ry="55" fill="#085041"/>

  <!-- Soleil / lune -->
  <circle cx="340" cy="120" r="40" fill="#FAC775" opacity="0.9"/>
</svg>
```

### Illustration de personnage simplifié (style flat)

```svg
<!-- Personnage stylisé -->
<g transform="translate(300, 150)">
  <!-- Tête -->
  <circle cx="0" cy="0" r="30" fill="#F5C4B3"/>
  <!-- Corps -->
  <rect x="-25" y="32" width="50" height="60" rx="8" fill="#534AB7"/>
  <!-- Bras gauche -->
  <rect x="-45" y="36" width="20" height="40" rx="8"
        fill="#534AB7" transform="rotate(-15, -35, 40)"/>
  <!-- Bras droit -->
  <rect x="25" y="36" width="20" height="40" rx="8"
        fill="#534AB7" transform="rotate(15, 35, 40)"/>
  <!-- Jambes -->
  <rect x="-22" y="90" width="18" height="45" rx="6" fill="#26215C"/>
  <rect x="4" y="90" width="18" height="45" rx="6" fill="#26215C"/>
</g>
```

---

## Affiches et posters

### Structure d'une affiche typique

```svg
<svg width="100%" viewBox="0 0 600 900">
  <!-- 1. Background (toute la surface) -->
  <rect width="600" height="900" fill="#0f0f23"/>

  <!-- 2. Artwork central (zone 200-700 verticalement) -->
  <!-- [art principal] -->

  <!-- 3. Zone de titre (haut, y=60-180) -->
  <text font-family="system-ui" font-size="64" font-weight="700"
        fill="#EEEDFE" x="300" y="120" text-anchor="middle">TITRE</text>
  <text font-family="system-ui" font-size="22" font-weight="400"
        fill="#AFA9EC" x="300" y="165" text-anchor="middle" letter-spacing="6">
    SOUS-TITRE
  </text>

  <!-- 4. Zone d'info (bas, y=750-870) -->
  <line x1="60" y1="780" x2="540" y2="780" stroke="#534AB7" stroke-width="0.5"/>
  <text font-family="system-ui" font-size="15" fill="#888780"
        x="300" y="820" text-anchor="middle">Date · Lieu · Détails</text>

  <!-- 5. Décorations : filets, points, numérotation -->
</svg>
```

### Typographie d'affiche

Contrairement aux diagrammes, les affiches ont une liberté typographique totale :
- Mélanger les tailles : titre 64-80px, sous-titre 20-28px, corps 14-18px
- `letter-spacing` positif sur les textes en majuscules : `letter-spacing="4"` à `letter-spacing="8"`
- `font-style="italic"` pour les titres expressifs
- Contraste maximum entre les zones de texte et leur fond

### Filets et ornements décoratifs
```svg
<!-- Filet double -->
<line x1="60" y1="200" x2="540" y2="200" stroke="#534AB7" stroke-width="2"/>
<line x1="60" y1="204" x2="540" y2="204" stroke="#534AB7" stroke-width="0.5"/>

<!-- Ornement de coin -->
<path d="M 40 40 L 40 80 M 40 40 L 80 40" stroke="#FAC775" stroke-width="2"
      fill="none" stroke-linecap="round"/>

<!-- Numérotation stylisée -->
<text font-family="Georgia, serif" font-size="120" fill="#534AB7" opacity="0.08"
      x="300" y="500" text-anchor="middle">01</text>
```

---

## Animations CSS dans les SVG/HTML

Uniquement `transform` et `opacity`. Toujours dans le media query.

```svg
<style>
  @media (prefers-reduced-motion: no-preference) {

    .float {
      animation: float 3s ease-in-out infinite;
      transform-origin: center;
    }
    @keyframes float {
      0%, 100% { transform: translateY(0); }
      50%       { transform: translateY(-12px); }
    }

    .pulse-opacity {
      animation: pulse 2s ease-in-out infinite;
    }
    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50%       { opacity: 0.4; }
    }

    .spin-slow {
      animation: spin 8s linear infinite;
      transform-origin: 300px 300px; /* centre du SVG */
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .draw {
      stroke-dasharray: 1000;
      stroke-dashoffset: 1000;
      animation: draw 2s ease forwards;
    }
    @keyframes draw {
      to { stroke-dashoffset: 0; }
    }

  }
</style>
```

---

## Textures et effets de surface (sans filter)

### Grain (simulation sans blur)
```svg
<defs>
  <pattern id="grain" width="4" height="4" patternUnits="userSpaceOnUse">
    <rect width="1" height="1" x="0" y="0" fill="#fff" opacity="0.03"/>
    <rect width="1" height="1" x="2" y="2" fill="#fff" opacity="0.03"/>
    <rect width="1" height="1" x="1" y="3" fill="#fff" opacity="0.02"/>
  </pattern>
</defs>
<rect width="600" height="600" fill="url(#grain)"/>
```

### Vignette (sans filter: blur)
```svg
<defs>
  <radialGradient id="vignette" cx="50%" cy="50%" r="70%">
    <stop offset="0%"   stop-color="black" stop-opacity="0"/>
    <stop offset="100%" stop-color="black" stop-opacity="0.5"/>
  </radialGradient>
</defs>
<rect width="600" height="600" fill="url(#vignette)"/>
```

### Scanlines (CRT / rétro)
```svg
<defs>
  <pattern id="scan" width="2" height="2" patternUnits="userSpaceOnUse">
    <rect width="2" height="1" y="0" fill="#000" opacity="0.05"/>
  </pattern>
</defs>
<rect width="600" height="600" fill="url(#scan)"/>
```
