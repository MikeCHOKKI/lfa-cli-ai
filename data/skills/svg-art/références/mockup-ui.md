# Référence : Mockups UI & Composants

Couvre : formulaires, dashboards, cartes, tables, navbars, modales, composants.

---

## Principe général

Les mockups UI sont toujours en **HTML**, jamais en SVG.
SVG pour les mockups = pas de hover, pas d'états, pas de scroll.
HTML = design system vivant avec états, interactions, dark mode.

---

## Squelette d'un mockup

```html
<style>
  *{box-sizing:border-box;margin:0;padding:0}
  body{font-family:system-ui,sans-serif;color:var(--color-text-primary)}
  .card{
    background:var(--color-background-secondary);
    border:1px solid var(--color-border-tertiary);
    border-radius:var(--border-radius-lg);
    padding:20px;
  }
  .btn{
    display:inline-flex;align-items:center;gap:6px;
    padding:8px 16px;border-radius:var(--border-radius-md);
    font-size:14px;font-weight:500;cursor:pointer;border:none;
    transition:opacity .15s;
  }
  .btn:hover{opacity:0.85}
  .btn-primary{background:#534AB7;color:#fff}
  .btn-secondary{
    background:transparent;
    border:1px solid var(--color-border-secondary);
    color:var(--color-text-primary);
  }
</style>

<div style="padding:24px;max-width:680px">
  <!-- contenu du mockup -->
</div>
```

---

## Formulaires

### Input standard
```html
<div style="display:flex;flex-direction:column;gap:6px;margin-bottom:16px">
  <label style="font-size:13px;font-weight:500;color:var(--color-text-primary)">
    Label
  </label>
  <input type="text" placeholder="Placeholder"
    style="
      padding:10px 12px;
      border-radius:var(--border-radius-md);
      border:1px solid var(--color-border-secondary);
      background:var(--color-background-primary);
      color:var(--color-text-primary);
      font-size:14px;
      width:100%;
      outline:none;
      transition:border-color .15s;
    "
    onfocus="this.style.borderColor='#534AB7'"
    onblur="this.style.borderColor='var(--color-border-secondary)'"/>
  <span style="font-size:12px;color:var(--color-text-tertiary)">Texte d'aide optionnel</span>
</div>
```

### Select
```html
<select style="
  padding:10px 12px;
  border-radius:var(--border-radius-md);
  border:1px solid var(--color-border-secondary);
  background:var(--color-background-primary);
  color:var(--color-text-primary);
  font-size:14px;width:100%;cursor:pointer;">
  <option>Option A</option>
  <option>Option B</option>
</select>
```

### Toggle switch complet
```html
<label style="display:flex;align-items:center;gap:10px;cursor:pointer">
  <span style="position:relative;display:inline-block;width:40px;height:22px">
    <input type="checkbox" id="tog" style="opacity:0;width:0;height:0;position:absolute"
           onchange="
             var t=this.nextElementSibling;
             t.style.background=this.checked?'#534AB7':'var(--color-border-secondary)';
             t.querySelector('span').style.transform=this.checked?'translateX(18px)':'translateX(0)';
           "/>
    <span style="
      position:absolute;inset:0;border-radius:11px;
      background:var(--color-border-secondary);transition:background .2s;cursor:pointer;">
      <span style="
        position:absolute;top:3px;left:3px;
        width:16px;height:16px;border-radius:50%;
        background:white;transition:transform .2s;pointer-events:none;"></span>
    </span>
  </span>
  <span style="font-size:14px;color:var(--color-text-primary)">Activer la fonctionnalité</span>
</label>
```

### Checkbox
```html
<label style="display:flex;align-items:center;gap:8px;cursor:pointer;font-size:14px">
  <input type="checkbox" style="
    width:16px;height:16px;accent-color:#534AB7;cursor:pointer;"/>
  Label de la case
</label>
```

### Radio group
```html
<div style="display:flex;flex-direction:column;gap:8px">
  <label style="display:flex;align-items:center;gap:8px;cursor:pointer;font-size:14px">
    <input type="radio" name="grp" value="a" style="accent-color:#534AB7"/> Option A
  </label>
  <label style="display:flex;align-items:center;gap:8px;cursor:pointer;font-size:14px">
    <input type="radio" name="grp" value="b" style="accent-color:#534AB7"/> Option B
  </label>
</div>
```

---

## Cartes et layouts

### Card standard
```html
<div class="card">
  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:12px">
    <h3 style="font-size:15px;font-weight:500;color:var(--color-text-primary)">Titre de la carte</h3>
    <span style="font-size:12px;color:var(--color-text-tertiary)">Sous-info</span>
  </div>
  <p style="font-size:14px;color:var(--color-text-secondary);line-height:1.5">
    Contenu de la carte.
  </p>
  <div style="margin-top:16px;padding-top:16px;border-top:1px solid var(--color-border-tertiary);
              display:flex;gap:8px">
    <button class="btn btn-primary">Action principale</button>
    <button class="btn btn-secondary">Annuler</button>
  </div>
</div>
```

### Grid de cartes
```html
<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:16px">
  <!-- répéter des .card ici -->
</div>
```

### Stat card (KPI)
```html
<div class="card" style="text-align:center">
  <div style="font-size:40px;font-weight:600;color:#534AB7;line-height:1">1,284</div>
  <div style="font-size:13px;color:var(--color-text-secondary);margin-top:4px">Utilisateurs actifs</div>
  <div style="font-size:12px;color:#3B6D11;margin-top:8px">↑ 12% ce mois</div>
</div>
```

---

## Tables

```html
<div style="overflow-x:auto;border-radius:var(--border-radius-lg);
            border:1px solid var(--color-border-tertiary)">
  <table style="width:100%;border-collapse:collapse;font-size:14px">
    <thead>
      <tr style="background:var(--color-background-secondary);
                 border-bottom:1px solid var(--color-border-secondary)">
        <th style="padding:10px 16px;text-align:left;font-weight:500;
                   color:var(--color-text-secondary);font-size:12px;text-transform:uppercase;
                   letter-spacing:.5px">Colonne A</th>
        <th style="padding:10px 16px;text-align:left;font-weight:500;
                   color:var(--color-text-secondary);font-size:12px;text-transform:uppercase;
                   letter-spacing:.5px">Colonne B</th>
        <th style="padding:10px 16px;text-align:right;font-weight:500;
                   color:var(--color-text-secondary);font-size:12px;text-transform:uppercase;
                   letter-spacing:.5px">Valeur</th>
      </tr>
    </thead>
    <tbody>
      <tr style="border-bottom:1px solid var(--color-border-tertiary);
                 transition:background .1s"
          onmouseover="this.style.background='var(--color-background-secondary)'"
          onmouseout="this.style.background='transparent'">
        <td style="padding:12px 16px;color:var(--color-text-primary)">Item A</td>
        <td style="padding:12px 16px;color:var(--color-text-secondary)">Détail</td>
        <td style="padding:12px 16px;text-align:right;font-weight:500">42</td>
      </tr>
    </tbody>
  </table>
</div>
```

---

## Navigation

### Navbar
```html
<nav style="
  display:flex;align-items:center;justify-content:space-between;
  padding:12px 24px;
  background:var(--color-background-primary);
  border-bottom:1px solid var(--color-border-tertiary);
  position:sticky;top:0;z-index:10;">
  <div style="font-size:16px;font-weight:600;color:var(--color-text-primary)">
    Marque
  </div>
  <div style="display:flex;gap:4px">
    <a href="#" style="padding:6px 12px;border-radius:var(--border-radius-md);
                       font-size:14px;color:var(--color-text-secondary);
                       text-decoration:none;transition:background .15s"
       onmouseover="this.style.background='var(--color-background-secondary)'"
       onmouseout="this.style.background='transparent'">Accueil</a>
    <a href="#" style="padding:6px 12px;border-radius:var(--border-radius-md);
                       font-size:14px;color:var(--color-text-primary);font-weight:500;
                       text-decoration:none;background:var(--color-background-secondary)">
      Page active</a>
  </div>
</nav>
```

### Tabs
```html
<div style="display:flex;border-bottom:1px solid var(--color-border-tertiary);margin-bottom:20px">
  <button onclick="switchTab(this,'tab1')" id="tab-btn-1"
    style="padding:10px 16px;border:none;background:none;cursor:pointer;
           font-size:14px;font-weight:500;color:#534AB7;
           border-bottom:2px solid #534AB7;margin-bottom:-1px;">
    Onglet 1
  </button>
  <button onclick="switchTab(this,'tab2')" id="tab-btn-2"
    style="padding:10px 16px;border:none;background:none;cursor:pointer;
           font-size:14px;color:var(--color-text-secondary);border-bottom:2px solid transparent;margin-bottom:-1px;">
    Onglet 2
  </button>
</div>
<div id="tab1">Contenu onglet 1</div>
<div id="tab2" style="display:none">Contenu onglet 2</div>
<script>
function switchTab(btn, id) {
  document.querySelectorAll('[id^=tab-btn]').forEach(b => {
    b.style.color = 'var(--color-text-secondary)';
    b.style.borderBottomColor = 'transparent';
  });
  btn.style.color = '#534AB7';
  btn.style.borderBottomColor = '#534AB7';
  document.querySelectorAll('[id^=tab]:not([id^=tab-btn])').forEach(t => t.style.display = 'none');
  document.getElementById(id).style.display = 'block';
}
</script>
```

---

## Badges et tags

```html
<!-- Badge coloré -->
<span style="
  display:inline-flex;align-items:center;gap:4px;
  padding:3px 8px;border-radius:20px;font-size:12px;font-weight:500;
  background:#EEEDFE;color:#534AB7;">
  ● Actif
</span>

<!-- Tag neutre -->
<span style="
  padding:2px 8px;border-radius:4px;font-size:12px;
  background:var(--color-background-secondary);
  color:var(--color-text-secondary);
  border:1px solid var(--color-border-tertiary);">
  tag
</span>
```

---

## Modale (overlay)

```html
<div id="modal-overlay" onclick="closeModal()"
  style="display:none;position:fixed;inset:0;background:rgba(0,0,0,.5);
         z-index:100;align-items:center;justify-content:center;">
  <div onclick="event.stopPropagation()"
    style="background:var(--color-background-primary);border-radius:var(--border-radius-xl);
           padding:28px;max-width:480px;width:90%;box-shadow:0 8px 40px rgba(0,0,0,.15)">
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px">
      <h2 style="font-size:16px;font-weight:500">Titre de la modale</h2>
      <button onclick="closeModal()" style="border:none;background:none;cursor:pointer;
              font-size:18px;color:var(--color-text-tertiary)">×</button>
    </div>
    <p style="font-size:14px;color:var(--color-text-secondary);line-height:1.5;margin-bottom:20px">
      Contenu de la modale.
    </p>
    <div style="display:flex;gap:8px;justify-content:flex-end">
      <button class="btn btn-secondary" onclick="closeModal()">Annuler</button>
      <button class="btn btn-primary">Confirmer</button>
    </div>
  </div>
</div>
<button class="btn btn-primary" onclick="document.getElementById('modal-overlay').style.display='flex'">
  Ouvrir la modale
</button>
<script>
function closeModal(){document.getElementById('modal-overlay').style.display='none'}
</script>
```

---

## Sidebar layout

```html
<div style="display:flex;height:100vh;overflow:hidden">
  <!-- Sidebar -->
  <aside style="width:220px;flex-shrink:0;
                background:var(--color-background-secondary);
                border-right:1px solid var(--color-border-tertiary);
                padding:20px 12px;display:flex;flex-direction:column;gap:4px">
    <div style="font-size:15px;font-weight:600;padding:8px;margin-bottom:12px">
      Menu
    </div>
    <a href="#" style="display:flex;align-items:center;gap:10px;padding:8px 12px;
                       border-radius:var(--border-radius-md);font-size:14px;
                       background:#EEEDFE;color:#534AB7;text-decoration:none;font-weight:500">
      ⊞ Dashboard
    </a>
    <a href="#" style="display:flex;align-items:center;gap:10px;padding:8px 12px;
                       border-radius:var(--border-radius-md);font-size:14px;
                       color:var(--color-text-secondary);text-decoration:none">
      ☰ Projets
    </a>
  </aside>

  <!-- Main content -->
  <main style="flex:1;overflow-y:auto;padding:28px">
    <!-- contenu principal -->
  </main>
</div>
```
