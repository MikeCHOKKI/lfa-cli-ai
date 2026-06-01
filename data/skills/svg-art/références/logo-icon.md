# Référence : Logos & Icônes

Couvre : logotypes, monogrammes, emblèmes, badges, icônes, favicons, pictogrammes.

---

## Philosophie du logo SVG

Un bon logo SVG est :
- **Lisible à toutes tailles** : de 16px (favicon) à 1000px (affiche)
- **Monochromatique possible** : fonctionne en noir seul
- **Vectoriel pur** : zéro rastérisation, zéro bitmap intégré
- **Mémorable** : forme simple, reconnaissable en silhouette

---

## ViewBox pour les logos

Contrairement aux diagrammes, les logos ont leur propre viewBox.

| Format     | ViewBox       | Usage                             |
|------------|---------------|-----------------------------------|
| Carré      | `0 0 400 400` | App icon, favicon, monogramme     |
| Horizontal | `0 0 600 200` | Logotype (marque + nom)           |
| Vertical   | `0 0 400 500` | Empilé (symbole au-dessus du nom) |
| Badge      | `0 0 300 300` | Emblème circulaire                |

Toujours `width="100%"` et conserver le viewBox — le logo s'adapte à son conteneur.

---

## Anatomie d'un logo

```svg
<svg width="100%" viewBox="0 0 600 200" role="img">
  <title>Logo NomMarca</title>

  <!-- 1. Symbole / marque graphique -->
  <g transform="translate(20, 20)">
    <!-- formes géométriques du symbole -->
  </g>

  <!-- 2. Logotype (nom de la marque) -->
  <text font-family="system-ui, sans-serif"
        font-size="48" font-weight="600"
        x="200" y="115" dominant-baseline="central">
    NomMarca
  </text>

  <!-- 3. Tagline (optionnel) -->
  <text font-family="system-ui, sans-serif"
        font-size="18" font-weight="400" letter-spacing="3"
        x="200" y="155" dominant-baseline="central">
    TAGLINE ICI
  </text>
</svg>
```

---

## Symboles graphiques — techniques de construction

### Monogramme (initiales stylisées)
Superposer les lettres avec rotation, opacité, ou découpe.
```svg
<g transform="translate(60, 60)">
  <text font-family="Georgia, serif" font-size="120" font-weight="700"
        fill="#534AB7" opacity="0.15" x="0" y="100">A</text>
  <text font-family="Georgia, serif" font-size="120" font-weight="700"
        fill="#534AB7" x="20" y="100">A</text>
</g>
```

### Forme géométrique simple
```svg
<!-- Hexagone (technologie, sécurité) -->
<polygon points="100,20 172,60 172,140 100,180 28,140 28,60"
         fill="#EEEDFE" stroke="#534AB7" stroke-width="3"/>

<!-- Cercle avec cutout (modernité) -->
<circle cx="100" cy="100" r="80" fill="#534AB7"/>
<circle cx="100" cy="100" r="50" fill="white"/>

<!-- Losange (dynamisme) -->
<polygon points="100,10 190,100 100,190 10,100"
         fill="none" stroke="#1D9E75" stroke-width="4"/>
```

### Construction modulaire — grille 8pt
Tout logo doit s'aligner sur une grille de 8px.
Unités recommandées : 8, 16, 24, 32, 40, 48, 56, 64, 80, 96.

```svg
<!-- Grille de conception (à retirer du résultat final) -->
<!-- Guide : 8px = 1 unité. Symbole min = 8×8u = 64px -->
```

### Path Bezier pour les formes organiques
```svg
<!-- Feuille / nature -->
<path d="M 100 180 C 140 120, 180 60, 100 20 C 20 60, 60 120, 100 180 Z"
      fill="#3B6D11"/>

<!-- Vague -->
<path d="M 0 100 Q 50 60, 100 100 Q 150 140, 200 100 Q 250 60, 300 100"
      fill="none" stroke="#378ADD" stroke-width="3"/>

<!-- Flèche stylisée -->
<path d="M 20 100 L 160 100 M 120 60 L 160 100 L 120 140"
      fill="none" stroke="#D85A30" stroke-width="5"
      stroke-linecap="round" stroke-linejoin="round"/>
```

---

## Couleurs pour les logos

Les logos **peuvent** utiliser des hex en dur. Choisir une palette primaire + secondaire :

**Palette recommandée : 1 couleur principale + 1 neutre**
```
Primaire : un stop 600 de n'importe quelle ramp
Neutre : #2C2C2A (sombre) ou #F1EFE8 (clair)
```

**Palettes prédéfinies pour logos :**
```
Tech moderne   : #534AB7 (purple 600) + #2C2C2A
Environnement  : #3B6D11 (green 600) + #085041 (teal 800)
Finance/confiance : #185FA5 (blue 600) + #2C2C2A
Énergie/impact : #D85A30 (coral 400) + #412402 (amber 900)
Santé/biotech  : #1D9E75 (teal 400) + #26215C (purple 900)
Luxe/premium   : #2C2C2A + #888780 (gray 400)
```

**Règle or/argent pour les logos premium :**
```svg
<defs>
  <linearGradient id="gold" x1="0" y1="0" x2="0" y2="1">
    <stop offset="0%" stop-color="#F5D675"/>
    <stop offset="50%" stop-color="#C9A227"/>
    <stop offset="100%" stop-color="#F5D675"/>
  </linearGradient>
</defs>
```

### Dark mode pour les logos
Toujours fournir les deux versions ou utiliser `currentColor` + CSS :
```svg
<style>
  @media (prefers-color-scheme: dark) {
    .logo-primary { fill: #AFA9EC; }
    .logo-text    { fill: #CECBF6; }
  }
  .logo-primary { fill: #534AB7; }
  .logo-text    { fill: #26215C; }
</style>
```

---

## Icônes — pictogrammes et interfaces

### ViewBox carré standard
```svg
<svg width="24" height="24" viewBox="0 0 24 24" role="img" aria-hidden="true">
  <!-- tracé de l'icône -->
</svg>
```
Tailles standard : 16, 20, 24, 32, 48px. Toujours conserver le ratio 1:1.

### Style de tracé (outline — recommandé)
```svg
<path d="M..." fill="none" stroke="currentColor"
      stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
```
`currentColor` hérite la couleur du texte parent — dark mode automatique.
`stroke-width="1.5"` à 24px. Ajuster proportionnellement : `1.0` à 16px, `2.0` à 32px.

### Icônes remplies (filled — pour states actifs)
```svg
<path d="M..." fill="currentColor"/>
```

### Grille de construction pour les icônes
Toujours travailler sur une grille 24×24 avec :
- Zone de contenu : 2px de marge → zone active 20×20
- Traits principaux : alignés sur la grille 2px (0, 2, 4, 6...)
- Coins arrondis : `stroke-linecap="round" stroke-linejoin="round"`
- Épaisseur optique : les diagonales à 45° semblent plus épaisses → réduire de 0.1px

### Icônes animées
```svg
<style>
  @media (prefers-reduced-motion: no-preference) {
    .icon-spin { animation: spin 1s linear infinite; }
    @keyframes spin { to { transform: rotate(360deg); } }
  }
</style>
<g class="icon-spin" transform-origin="12 12">
  <!-- loader arc -->
  <path d="M 12 2 A 10 10 0 0 1 22 12"
        fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
</g>
```

### Bibliothèque d'icônes communes (paths 24×24)

**Paramètres / engrenage :**
```svg
<path d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z M9.5 2.1l-.5 1.3a7 7 0 0 0-1.8 1l-1.3-.4L4.1 6l.9 1.1a7 7 0 0 0 0 2.1L4.1 10l1.8 2 1.3-.4a7 7 0 0 0 1.8 1l.5 1.3H14l.5-1.3a7 7 0 0 0 1.8-1l1.3.4L19.5 10l-.9-1.1a7 7 0 0 0 0-2.1l.9-1.1-1.8-2-1.3.4a7 7 0 0 0-1.8-1L14 2.1H9.5Z"/>
```

**Flèche droite :**
```svg
<path d="M5 12h14M13 6l6 6-6 6"/>
```

**Check :**
```svg
<path d="M20 6L9 17l-5-5"/>
```

**Croix :**
```svg
<path d="M18 6L6 18M6 6l12 12"/>
```

**Plus :**
```svg
<path d="M12 5v14M5 12h14"/>
```

**Cadenas :**
```svg
<rect x="3" y="11" width="18" height="11" rx="2"/>
<path d="M7 11V7a5 5 0 0 1 10 0v4"/>
```

**Utilisateur :**
```svg
<circle cx="12" cy="8" r="4"/>
<path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
```

---

## Badge / emblème

Format circulaire ou en bouclier, pour les certifications, rangs, récompenses.

```svg
<svg width="100%" viewBox="0 0 200 200" role="img">
  <title>Badge Expert</title>

  <!-- Fond circulaire -->
  <circle cx="100" cy="100" r="95" fill="#EEEDFE" stroke="#534AB7" stroke-width="3"/>

  <!-- Cercle intérieur décoratif -->
  <circle cx="100" cy="100" r="78" fill="none" stroke="#534AB7" stroke-width="1" stroke-dasharray="4 3"/>

  <!-- Symbole central -->
  <text font-size="64" text-anchor="middle" dominant-baseline="central"
        x="100" y="95">⬡</text>

  <!-- Texte circulaire en arc -->
  <path id="arc-top" d="M 25 100 A 75 75 0 0 1 175 100" fill="none"/>
  <text font-size="13" font-weight="500" fill="#534AB7" letter-spacing="2">
    <textPath href="#arc-top" startOffset="50%" text-anchor="middle">CERTIFIÉ EXPERT</textPath>
  </text>

  <!-- Étoiles décoratives -->
  <text font-size="10" x="60" y="155" text-anchor="middle">★</text>
  <text font-size="10" x="100" y="162" text-anchor="middle">★</text>
  <text font-size="10" x="140" y="155" text-anchor="middle">★</text>
</svg>
```

---

## Favicon SVG

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <!-- Simple, une seule forme ou lettre, lisible à 16×16 -->
  <circle cx="16" cy="16" r="14" fill="#534AB7"/>
  <text font-family="system-ui" font-size="18" font-weight="700"
        fill="white" x="16" y="16"
        text-anchor="middle" dominant-baseline="central">L</text>
</svg>
```

Règles favicon :
- Max 2 couleurs
- Aucun détail fin (disparu à 16px)
- Toujours tester mentalement à 16×16
