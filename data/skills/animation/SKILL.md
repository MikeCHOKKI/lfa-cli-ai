---
name: animation
description: >
  Crée des animations CSS et SVG de qualité production : micro-interactions,
  transitions de page, chargements, scroll animations, keyframes complexes,
  morphing SVG, animations déclenchées par Intersection Observer.
  Déclenche ce skill quand l'utilisateur demande "animation", "transition",
  "keyframe", "micro-interaction", "loader", "spinner", "effet", "mouvement",
  "scroll animation", "parallax", "morphing", "animé" ou tout rendu dynamique.
---

# Animation Skill

Animations CSS et SVG fluides, performantes et accessibles.

---

## Principes fondamentaux

### Timing
- UI générale : 150-200ms (hover, focus, toggle)
- Transitions de page : 300-400ms
- Micro-interactions : 100-150ms
- Loaders : 600-1200ms (loop)
- Entrée en scène : 400-800ms

### Easing — courbes de Bézier

```css
:root {
  --ease-out: cubic-bezier(0.16, 1, 0.3, 1);
  --ease-in-out: cubic-bezier(0.65, 0, 0.35, 1);
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);
  --ease-smooth: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-bounce: cubic-bezier(0.18, 0.89, 0.32, 1.28);
}
```

| Usage | Easing |
|-------|--------|
| Entrée | `--ease-out` |
| Sortie | `--ease-in-out` |
| Élastique | `--ease-spring` |
| Standard UI | `--ease-smooth` |
| Rebond | `--ease-bounce` |

---

## Patterns CSS

### Fade in / fade out

```css
@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fade-out {
  from { opacity: 1; }
  to { opacity: 0; }
}

.fade-in { animation: fade-in 0.3s var(--ease-out); }
.fade-out { animation: fade-out 0.2s var(--ease-in-out) forwards; }
```

### Slide

```css
@keyframes slide-up {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slide-down {
  from { opacity: 1; transform: translateY(0); }
  to { opacity: 0; transform: translateY(16px); }
}

@keyframes slide-in-right {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}
```

### Scale

```css
@keyframes scale-in {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

@keyframes scale-bounce {
  0% { transform: scale(0); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}
```

### Loader spinner

```css
.spinner { width: 24px; height: 24px; border: 3px solid var(--border); 
           border-top-color: var(--accent); border-radius: 50%; 
           animation: spin 0.8s linear infinite; }

@keyframes spin { to { transform: rotate(360deg); } }
```

### Loader pulse

```css
.pulse { width: 8px; height: 8px; border-radius: 50%; background: var(--accent); 
         animation: pulse 1.5s ease-in-out infinite; }

@keyframes pulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.5); }
}
```

### Loader bar

```css
.progress { height: 4px; background: var(--border); border-radius: 2px; overflow: hidden; }
.progress-bar { height: 100%; background: var(--accent); border-radius: 2px; 
                animation: progress 2s ease-in-out infinite; width: 40%; }

@keyframes progress {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(350%); }
}
```

---

## Transitions de composants

### Hover card

```css
.card { transition: transform 0.2s var(--ease-out), box-shadow 0.2s var(--ease-out); }
.card:hover { transform: translateY(-2px); box-shadow: var(--shadow-lg); }
```

### Button feedback

```css
.btn:active { transform: scale(0.97); }
```

### Accordion

```css
.accordion-content { 
  display: grid; 
  grid-template-rows: 0fr; 
  transition: grid-template-rows 0.3s var(--ease-out); 
}
.accordion-content.open { grid-template-rows: 1fr; }
.accordion-content > div { overflow: hidden; }
```

---

## Scroll animations (Intersection Observer)

```html
<style>
  .reveal { opacity: 0; transform: translateY(30px); transition: all 0.6s var(--ease-out); }
  .reveal.visible { opacity: 1; transform: translateY(0); }
  .reveal-delay-1 { transition-delay: 0.1s; }
  .reveal-delay-2 { transition-delay: 0.2s; }
  .reveal-delay-3 { transition-delay: 0.3s; }
</style>

<script>
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('visible') })
  }, { threshold: 0.1 })
  document.querySelectorAll('.reveal').forEach(el => observer.observe(el))
</script>
```

---

## Animations SVG

### Stroke dash (dessin progressif)

```css
.draw-path { stroke-dasharray: 1000; stroke-dashoffset: 1000; 
             animation: draw 2s var(--ease-out) forwards; }

@keyframes draw { to { stroke-dashoffset: 0; } }
```

### Morphing SVG

```css
.morph { transition: d 0.4s var(--ease-smooth); }
```

### Rotation continue

```css
@keyframes rotate { to { transform: rotate(360deg); } }
.gear { animation: rotate 4s linear infinite; transform-origin: center; }
```

---

## Accessibilité

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

Toujours inclure la media query `prefers-reduced-motion` dans tout output qui contient des animations.

---

## Performance

- `transform` et `opacity` seulement pour les animations (composite)
- Pas d'animation sur `width`, `height`, `top`, `left`, `margin`, `padding`
- Utiliser `will-change: transform` avec parcimonie
- `contain: layout style` sur les conteneurs animés si possible
- Pas plus de 20 animations simultanées
- Préférer `transition` à `@keyframes` pour les micro-interactions

---

## Checklist

1. `prefers-reduced-motion` présent ?
2. Easing adapté au type d'animation ?
3. Durée cohérente (150-200ms UI, 300-400ms page) ?
4. Pas d'animation sur des propriétés coûteuses (layout/paint) ?
5. Accessibilité : mouvement réductible ?
6. Animation boucle : pause on blur (optionnel) ?
