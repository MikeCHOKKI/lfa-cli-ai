# Référence : Diagrammes

Couvre : flowchart, structurel, illustratif, ERD, séquence, chart de données.

---

## Flowchart — processus séquentiels, pipelines, décisions

### Quand l'utiliser
Flux étape par étape avec branchements conditionnels. Pas pour expliquer
un mécanisme complexe (→ illustratif) ni une hiérarchie (→ structural).

### Nœud une ligne (hauteur 44px)
```svg
<g class="node c-blue">
  <rect x="100" y="20" width="180" height="44" rx="8" stroke-width="0.5"/>
  <text class="th" x="190" y="42" text-anchor="middle" dominant-baseline="central">
    Mon étape
  </text>
</g>
```

### Nœud deux lignes (hauteur 60px)
```svg
<g class="node c-teal">
  <rect x="100" y="20" width="200" height="60" rx="8" stroke-width="0.5"/>
  <text class="th" x="200" y="40" text-anchor="middle" dominant-baseline="central">Titre</text>
  <text class="ts" x="200" y="58" text-anchor="middle" dominant-baseline="central">Sous-titre</text>
</g>
```

### Nœud losange (décision)
```svg
<polygon points="200,20 260,50 200,80 140,50" class="c-amber"
         stroke-width="0.5"/>
<text class="ts" x="200" y="50" text-anchor="middle" dominant-baseline="central">Oui ?</text>
```

### Nœud terminal (start / end) — coins arrondis extrêmes
```svg
<rect x="100" y="20" width="140" height="36" rx="18" class="c-gray" stroke-width="0.5"/>
<text class="th" x="170" y="38" text-anchor="middle" dominant-baseline="central">Début</text>
```

### Connexions
```svg
<!-- Ligne droite verticale -->
<line x1="190" y1="64" x2="190" y2="100" class="arr" marker-end="url(#arrow)"/>

<!-- Détour en L (anti-croisement) -->
<path d="M 190 64 L 190 90 L 340 90 L 340 100"
      fill="none" class="arr" marker-end="url(#arrow)"/>

<!-- Label sur connexion -->
<text class="ts" x="265" y="86" text-anchor="middle">Non</text>
```

Règle anti-croisement : avant chaque flèche, vérifier que son chemin ne traverse
pas une boîte tierce. Si oui, utiliser le détour en L avec y_milieu = entre les deux nœuds.

Espacement : 60px min entre boîtes, 10px entre tête de flèche et bord de boîte.

**Budget de complexité : max 7 nœuds. Au-delà → plusieurs diagrammes avec texte entre eux.**

Cycles : pas de ring fermé. Utiliser une flèche de retour courbée latérale + `↻ label`.

---

## Diagramme structurel — contenances, architectures, hiérarchies

### Quand l'utiliser
Ce qui est dans quoi : systèmes, comptes, serveurs, modules, réseaux.
Pas de flèches de processus — des zones imbriquées.

### Conteneur externe
```svg
<g class="c-gray">
  <rect x="16" y="16" width="648" height="300" rx="16" stroke-width="0.5"/>
  <text class="th" x="340" y="40" text-anchor="middle">Nom du conteneur</text>
  <text class="ts" x="340" y="56" text-anchor="middle">Sous-titre / rôle</text>
</g>
```

### Région interne (dans le conteneur)
```svg
<g class="c-purple">
  <rect x="32" y="72" width="300" height="120" rx="12" stroke-width="0.5"/>
  <text class="th" x="182" y="94" text-anchor="middle">Région</text>
</g>
```

Règles :
- Padding min 20px entre bord du conteneur et son contenu direct
- Chaque niveau d'imbrication = ramp de couleur différente du parent
- Flèches entre régions : `stroke-width="1"` + label court 1-2 mots
- Les connexions inter-régions passent toujours à l'extérieur des boîtes
- Max 3 niveaux d'imbrication

---

## Diagramme illustratif — mécanismes, fonctionnement, intuition

### Quand l'utiliser
Expliquer *comment ça marche* : attention transformer, cycle Krebs, moteur,
réseau de neurones. Le sujet a une géométrie propre, pas une grille.

### Deux saveurs

**Physique** (objet du monde réel) :
- Dessiner une coupe simplifiée ou une vue éclatée
- Proportions approximatives mais reconnaissables
- `<path>` et `<ellipse>` libres — pas limité aux rects
- Labels à l'extérieur avec `class="leader"` pointillée, tous du même côté

**Abstrait** (concept mathématique ou informatique) :
- Inventer une métaphore spatiale qui rend le mécanisme évident
- Ex : attention = projecteurs sur des tokens ; gradient descent = bille qui roule

### Ce qui est autorisé ici (pas ailleurs)
- 1 `<linearGradient>` pour encoder une propriété continue (température, intensité)
- Formes libres : `<path>`, `<ellipse>`, `<polygon>`, `<circle>` combinées
- Chevauchements de formes OK — pas des chevauchements de labels
- Couleur = intensité : amber/coral = chaud/énergie, blue/teal = froid/calme

### Gradient autorisé
```svg
<defs>
  <linearGradient id="heat" x1="0" y1="0" x2="1" y2="0">
    <stop offset="0%" stop-color="#85B7EB"/>
    <stop offset="100%" stop-color="#D85A30"/>
  </linearGradient>
</defs>
<rect x="100" y="200" width="480" height="20" fill="url(#heat)" rx="4"/>
```

### Labels avec leader lines
```svg
<!-- Label à droite avec pointillé -->
<line x1="220" y1="150" x2="380" y2="150" class="leader"/>
<text class="ts" x="385" y="150" dominant-baseline="central">Membrane</text>
```
Réserver 140px de marge côté labels. Choisir un seul côté pour tous.

### Préférer l'interactif
Si le sujet a un paramètre contrôlable (vitesse, température, fréquence) →
encapsuler dans du HTML avec un `<input type="range">` et SVG inline.
Voir `références/interactive.md`.

---

## ERD — schémas de base de données

Utiliser mermaid.js, pas SVG manuel.

```html
<div id="erd" style="width:100%"></div>
<script type="module">
import mermaid from 'https://esm.sh/mermaid@11/dist/mermaid.esm.min.mjs';
const dark = matchMedia('(prefers-color-scheme: dark)').matches;
await document.fonts.ready;
mermaid.initialize({
  startOnLoad: false,
  theme: 'base',
  themeVariables: {
    darkMode: dark,
    fontSize: '13px',
    fontFamily: 'system-ui, sans-serif',
    lineColor: dark ? '#9c9a92' : '#73726c',
    textColor: dark ? '#c2c0b6' : '#3d3d3a',
    primaryColor: dark ? '#3C3489' : '#EEEDFE',
    primaryBorderColor: dark ? '#AFA9EC' : '#534AB7',
  },
});
const { svg } = await mermaid.render('erd-diagram', `erDiagram
  USERS ||--o{ POSTS : writes
  USERS {
    uuid id PK
    string email
    string name
  }
  POSTS {
    uuid id PK
    uuid user_id FK
    string title
    text content
  }`);
document.getElementById('erd').innerHTML = svg;
</script>
```

---

## Diagramme de séquence — interactions entre acteurs

Utiliser mermaid.js également :

```html
<div id="seq" style="width:100%"></div>
<script type="module">
import mermaid from 'https://esm.sh/mermaid@11/dist/mermaid.esm.min.mjs';
const dark = matchMedia('(prefers-color-scheme: dark)').matches;
await document.fonts.ready;
mermaid.initialize({
  startOnLoad: false, theme: 'base',
  themeVariables: { darkMode: dark, fontSize: '13px' },
});
const { svg } = await mermaid.render('seq-diagram', `sequenceDiagram
  participant U as Utilisateur
  participant API as API
  participant DB as Base de données
  U->>API: POST /login
  API->>DB: SELECT user WHERE email=?
  DB-->>API: user record
  API-->>U: JWT token`);
document.getElementById('seq').innerHTML = svg;
</script>
```

---

## Chart de données — comparaisons, évolutions, distributions

Toujours HTML + bibliothèque, jamais SVG manuel pour les données.

### Bar chart (recharts)
```html
<div id="chart" style="width:100%;height:300px"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/recharts/2.8.0/Recharts.min.js"></script>
<script>
const { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } = Recharts;
const data = [
  { name: 'Jan', valeur: 400 },
  { name: 'Fév', valeur: 300 },
  { name: 'Mar', valeur: 600 },
];
const el = React.createElement(ResponsiveContainer, { width: '100%', height: 300 },
  React.createElement(BarChart, { data },
    React.createElement(XAxis, { dataKey: 'name' }),
    React.createElement(YAxis, null),
    React.createElement(Tooltip, null),
    React.createElement(Bar, { dataKey: 'valeur', fill: '#7F77DD', radius: [4,4,0,0] })
  )
);
ReactDOM.render(el, document.getElementById('chart'));
</script>
```

Couleurs recommandées pour les charts (hex direct, pas de classes) :
- Purple #7F77DD, Teal #1D9E75, Coral #D85A30, Blue #378ADD
- Toujours utiliser les stops 400 de la palette pour la lisibilité

### Règles charts
- Toujours un titre en `<p>` ou `<h3>` au-dessus du chart
- Légende si plusieurs séries
- Tooltips activés
- Axes avec unités si applicable
