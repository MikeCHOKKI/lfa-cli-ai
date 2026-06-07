---
name: mockup-ui
description: >
  Crée des maquettes UI de qualité production en HTML/CSS : pages web, dashboards,
  formulaires, cartes, landing pages, composants isolés. Design responsive,
  accessible, dark mode inclus. Déclenche ce skill quand l'utilisateur demande
  "mockup", "maquette", "interface", "dashboard", "landing page", "page web",
  "composant", "formulaire", "carte UI" ou toute demande de rendu visuel interactif.
---

# Mockup UI Skill

Crée des maquettes UI réalistes, responsives et accessibles.

---

## Stack technique

```html
<style>
  /* Variables CSS système */
  :root { /* design tokens */ }
  @media (prefers-color-scheme: dark) { :root { /* dark */ } }
  /* Layout, composants, responsive */
</style>

<div id="app">
  <!-- Structure de la maquette -->
</div>

<script>
  // Interactivité minimale (si nécessaire)
</script>
```

- Pas de `<!DOCTYPE>`, `<html>`, `<head>`, `<body>`
- Librairies autorisées depuis : cdnjs.cloudflare.com, esm.sh, cdn.jsdelivr.net, unpkg.com
- Frameworks UI autorisés : Alpine.js, Petite-Vue, Shoelace, Lit
- Pas de React, Vue, Svelte (sauf si explicitement demandé)

---

## Design system intégré

### Palette de couleurs (par défaut — override avec le skill design-system)

```css
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --bg-tertiary: #f1f3f5;
  --text-primary: #1a1a1a;
  --text-secondary: #4a4a4a;
  --text-tertiary: #8a8a8a;
  --border: #e2e4e8;
  --border-hover: #c4c8cd;
  --accent: #4f6ef7;
  --accent-hover: #3b57d9;
  --success: #2e8b57;
  --warning: #d4a030;
  --danger: #d32f2f;
  --radius: 8px;
  --radius-lg: 12px;
  --shadow: 0 1px 3px rgba(0,0,0,0.08), 0 1px 2px rgba(0,0,0,0.06);
  --shadow-lg: 0 10px 25px rgba(0,0,0,0.1);
  --font: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  --font-mono: 'JetBrains Mono', 'SF Mono', monospace;
}
```

### Dark mode

```css
@media (prefers-color-scheme: dark) {
  :root {
    --bg-primary: #121212;
    --bg-secondary: #1e1e1e;
    --bg-tertiary: #2a2a2a;
    --text-primary: #f0f0f0;
    --text-secondary: #b0b0b0;
    --text-tertiary: #707070;
    --border: #333;
    --border-hover: #555;
    --shadow: 0 1px 3px rgba(0,0,0,0.3);
    --shadow-lg: 0 10px 25px rgba(0,0,0,0.4);
  }
}
```

---

## Patterns de layout

### Dashboard

```html
<div class="dashboard">
  <aside class="sidebar">
    <nav class="nav">...</nav>
  </aside>
  <main class="main">
    <header class="topbar">...</header>
    <div class="grid">
      <div class="card">...</div>
      <!-- cartes stats, charts, tables -->
    </div>
  </main>
</div>
```

```css
.dashboard { display: grid; grid-template-columns: 240px 1fr; min-height: 100vh; }
.sidebar { background: var(--bg-secondary); border-right: 1px solid var(--border); }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 16px; }
```

### Landing page

```html
<header class="hero">
  <h1>Titre principal</h1>
  <p>Sous-titre descriptif</p>
  <div class="cta-group">...</div>
</header>
<section class="features">
  <div class="feature-card">...</div>
</section>
<section class="testimonials">...</section>
<footer class="footer">...</footer>
```

### Formulaire

```css
.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-group label { font-size: 14px; font-weight: 500; color: var(--text-secondary); }
.form-group input, .form-group select, .form-group textarea {
  padding: 10px 12px; border: 1px solid var(--border); border-radius: var(--radius);
  font-size: 15px; background: var(--bg-primary); color: var(--text-primary);
  transition: border-color .15s;
}
.form-group input:focus { outline: none; border-color: var(--accent); box-shadow: 0 0 0 3px rgba(79,110,247,0.15); }
```

---

## Composants récurrents

### Carte (Card)

```html
<div class="card">
  <div class="card-header">
    <h3>Titre</h3>
    <span class="badge">Statut</span>
  </div>
  <div class="card-body">
    <p>Contenu de la carte avec description et métriques.</p>
  </div>
  <div class="card-footer">
    <button class="btn btn-primary">Action</button>
  </div>
</div>
```

```css
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--radius-lg); 
        padding: 20px; box-shadow: var(--shadow); }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.card-footer { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; padding-top: 12px; 
               border-top: 1px solid var(--border); }
```

### Boutons

```css
.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; 
       border-radius: var(--radius); font-size: 14px; font-weight: 500; 
       border: none; cursor: pointer; transition: all .15s; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-primary:hover { background: var(--accent-hover); }
.btn-secondary { background: var(--bg-secondary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-secondary:hover { border-color: var(--border-hover); }
.btn-ghost { background: transparent; color: var(--text-secondary); }
.btn-ghost:hover { background: var(--bg-secondary); }
.btn-sm { padding: 4px 10px; font-size: 13px; }
.btn-lg { padding: 12px 24px; font-size: 16px; }
```

### Tableau

```css
.table-container { overflow-x: auto; border: 1px solid var(--border); border-radius: var(--radius-lg); }
table { width: 100%; border-collapse: collapse; }
th { text-align: left; padding: 12px 16px; font-size: 12px; font-weight: 600; 
     text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-tertiary); 
     background: var(--bg-secondary); border-bottom: 1px solid var(--border); }
td { padding: 12px 16px; border-bottom: 1px solid var(--border); color: var(--text-primary); font-size: 14px; }
tr:last-child td { border-bottom: none; }
tr:hover td { background: var(--bg-secondary); }
```

### Navigation

```css
.nav { display: flex; flex-direction: column; gap: 2px; padding: 8px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 8px 12px; 
            border-radius: var(--radius); color: var(--text-secondary); 
            text-decoration: none; font-size: 14px; transition: all .15s; }
.nav-item:hover { background: var(--bg-tertiary); color: var(--text-primary); }
.nav-item.active { background: rgba(79,110,247,0.1); color: var(--accent); font-weight: 500; }
```

### Modal / Dialog

```css
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); 
                 display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: var(--bg-primary); border-radius: var(--radius-lg); 
         padding: 24px; max-width: 480px; width: 90%; box-shadow: var(--shadow-lg); }
```

---

## Responsive breakpoints

```css
/* Mobile-first */
/* ≥ 640px  → sm */
/* ≥ 768px  → md */
/* ≥ 1024px → lg */
/* ≥ 1280px → xl */

@media (max-width: 767px) {
  .dashboard { grid-template-columns: 1fr; }
  .sidebar { display: none; }
  .grid { grid-template-columns: 1fr; }
}
```

---

## Règles de qualité

1. Design tokens en haut du CSS, toujours
2. Dark mode systématique
3. Mobile-first responsive
4. Accessibilité : labels, focus visible, `aria-*`, rôles
5. Pas de contenu factice générique — utiliser le contexte utilisateur
6. États : empty, loading, error, success pour chaque composant
7. Transitions douces (150-200ms) sur hover/focus
8. Pas de `!important` sauf utilitaires
9. Espacements cohérents (multiples de 4px)
