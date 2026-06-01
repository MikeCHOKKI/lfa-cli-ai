# Référence : Widgets Interactifs & Animations

Couvre : steppers, sliders, calculateurs, animateurs, explorateurs de concepts,
simulations, jeux simples.

---

## Principe

Un widget interactif donne le **contrôle** à l'utilisateur sur un paramètre du système.
Cela vaut infiniment plus qu'une image statique quand le sujet a des dimensions continues
(vitesse, température, fréquence, budget, probabilité).

Choisir le bon contrôle :
- Valeur continue → `<input type="range">` (slider)
- Choix discret → boutons ou `<select>`
- Binaire → toggle switch (voir `mockup-ui.md`)
- Séquence étape par étape → stepper
- Temps → animation avec play/pause

---

## Stepper — contenu paginé

```html
<style>
  .step{display:none}
  .step.active{display:block}
  .step-btn{
    padding:8px 20px;border-radius:8px;border:none;
    font-size:14px;font-weight:500;cursor:pointer;
    background:#534AB7;color:#fff;transition:opacity .15s
  }
  .step-btn:disabled{opacity:.4;cursor:default}
  .step-btn.sec{background:var(--color-background-secondary);
                color:var(--color-text-primary);
                border:1px solid var(--color-border-secondary)}
</style>

<div style="max-width:600px">
  <!-- Indicateur de progression -->
  <div style="display:flex;align-items:center;gap:0;margin-bottom:24px" id="progress">
    <!-- généré par JS -->
  </div>

  <!-- Étapes -->
  <div class="step active" id="step-0">
    <h3 style="font-size:16px;font-weight:500;margin-bottom:8px">Étape 1</h3>
    <p style="font-size:14px;color:var(--color-text-secondary);line-height:1.6">Contenu étape 1.</p>
  </div>
  <div class="step" id="step-1">
    <h3 style="font-size:16px;font-weight:500;margin-bottom:8px">Étape 2</h3>
    <p style="font-size:14px;color:var(--color-text-secondary);line-height:1.6">Contenu étape 2.</p>
  </div>
  <div class="step" id="step-2">
    <h3 style="font-size:16px;font-weight:500;margin-bottom:8px">Étape 3</h3>
    <p style="font-size:14px;color:var(--color-text-secondary);line-height:1.6">Contenu étape 3.</p>
  </div>

  <!-- Navigation -->
  <div style="display:flex;gap:8px;margin-top:24px;align-items:center">
    <button class="step-btn sec" id="prev" onclick="move(-1)">← Précédent</button>
    <button class="step-btn" id="next" onclick="move(1)">Suivant →</button>
    <span id="counter" style="font-size:13px;color:var(--color-text-tertiary);margin-left:8px">1 / 3</span>
  </div>
</div>

<script>
const TOTAL = 3;
let cur = 0;

function move(d) {
  document.getElementById('step-'+cur).classList.remove('active');
  cur = Math.max(0, Math.min(TOTAL-1, cur+d));
  document.getElementById('step-'+cur).classList.add('active');
  document.getElementById('prev').disabled = cur === 0;
  document.getElementById('next').disabled = cur === TOTAL-1;
  document.getElementById('counter').textContent = (cur+1)+' / '+TOTAL;
  renderProgress();
}

function renderProgress() {
  const p = document.getElementById('progress');
  p.innerHTML = Array.from({length:TOTAL}, (_,i) =>
    `<span style="
      width:${i===cur?32:8}px;height:8px;border-radius:4px;
      background:${i<=cur?'#534AB7':'var(--color-border-secondary)'};
      transition:all .3s;margin:0 2px"></span>`
  ).join('');
}

move(0);
</script>
```

---

## Slider interactif avec visualisation live

```html
<div style="max-width:600px;padding:20px">
  <div style="margin-bottom:20px">
    <label style="font-size:14px;font-weight:500;display:flex;justify-content:space-between">
      <span>Paramètre</span>
      <span id="val-display" style="color:#534AB7;font-weight:600">50</span>
    </label>
    <input type="range" id="slider" min="0" max="100" value="50"
      style="width:100%;margin-top:8px;accent-color:#534AB7"
      oninput="updateViz(this.value)"/>
    <div style="display:flex;justify-content:space-between;font-size:12px;
                color:var(--color-text-tertiary);margin-top:2px">
      <span>0</span><span>50</span><span>100</span>
    </div>
  </div>

  <!-- Visualisation SVG réactive -->
  <svg width="100%" viewBox="0 0 560 120" id="viz">
    <rect x="0" y="40" width="560" height="40" rx="6"
          fill="var(--color-background-secondary)"/>
    <rect id="fill-bar" x="0" y="40" width="280" height="40" rx="6" fill="#534AB7"/>
    <circle id="dot" cx="280" cy="60" r="14" fill="white"
            stroke="#534AB7" stroke-width="2"/>
    <text id="dot-label" x="280" y="60" text-anchor="middle"
          dominant-baseline="central" font-size="12" font-weight="500" fill="#534AB7">50</text>
  </svg>
</div>

<script>
function updateViz(v) {
  document.getElementById('val-display').textContent = v;
  const pct = v / 100;
  document.getElementById('fill-bar').setAttribute('width', 560 * pct);
  document.getElementById('dot').setAttribute('cx', 560 * pct);
  document.getElementById('dot-label').setAttribute('x', 560 * pct);
  document.getElementById('dot-label').textContent = v;
}
</script>
```

---

## Calculateur

```html
<div style="max-width:480px">
  <div class="card" style="margin-bottom:16px">
    <h3 style="font-size:15px;font-weight:500;margin-bottom:16px">Calculateur</h3>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:16px">
      <div>
        <label style="font-size:12px;font-weight:500;color:var(--color-text-secondary)">
          Valeur A
        </label>
        <input type="number" id="val-a" value="100"
          oninput="calculate()"
          style="margin-top:4px;padding:8px 12px;width:100%;border-radius:8px;
                 border:1px solid var(--color-border-secondary);
                 background:var(--color-background-primary);
                 color:var(--color-text-primary);font-size:14px"/>
      </div>
      <div>
        <label style="font-size:12px;font-weight:500;color:var(--color-text-secondary)">
          Valeur B
        </label>
        <input type="number" id="val-b" value="20"
          oninput="calculate()"
          style="margin-top:4px;padding:8px 12px;width:100%;border-radius:8px;
                 border:1px solid var(--color-border-secondary);
                 background:var(--color-background-primary);
                 color:var(--color-text-primary);font-size:14px"/>
      </div>
    </div>

    <!-- Résultats -->
    <div style="background:var(--color-background-tertiary);border-radius:8px;padding:16px">
      <div style="display:flex;justify-content:space-between;margin-bottom:8px">
        <span style="font-size:14px;color:var(--color-text-secondary)">Total</span>
        <span id="result-total" style="font-size:14px;font-weight:600;color:var(--color-text-primary)">120</span>
      </div>
      <div style="display:flex;justify-content:space-between">
        <span style="font-size:14px;color:var(--color-text-secondary)">Ratio</span>
        <span id="result-ratio" style="font-size:14px;font-weight:600;color:#534AB7">83%</span>
      </div>
    </div>
  </div>
</div>

<script>
function calculate() {
  const a = parseFloat(document.getElementById('val-a').value) || 0;
  const b = parseFloat(document.getElementById('val-b').value) || 0;
  document.getElementById('result-total').textContent = (a + b).toLocaleString();
  document.getElementById('result-ratio').textContent =
    a+b > 0 ? Math.round(a/(a+b)*100)+'%' : '—';
}
</script>
```

---

## Animation avec contrôle play/pause

```html
<div style="text-align:center;padding:20px">
  <svg width="100%" viewBox="0 0 560 200" id="anim-svg">
    <!-- Objet animé -->
    <circle id="ball" cx="40" cy="100" r="20" fill="#534AB7"/>
    <!-- Chemin guide -->
    <line x1="40" y1="100" x2="520" y2="100"
          stroke="var(--color-border-tertiary)" stroke-width="1" stroke-dasharray="4 4"/>
  </svg>

  <div style="display:flex;gap:12px;justify-content:center;margin-top:16px">
    <button class="btn btn-primary" id="play-btn" onclick="toggleAnim()">▶ Play</button>
    <button class="btn btn-secondary" onclick="resetAnim()">↺ Reset</button>
  </div>
</div>

<script>
let running = false, pos = 40, raf;

function step() {
  pos += 2;
  if (pos > 520) pos = 40;
  document.getElementById('ball').setAttribute('cx', pos);
  if (running) raf = requestAnimationFrame(step);
}

function toggleAnim() {
  running = !running;
  document.getElementById('play-btn').textContent = running ? '⏸ Pause' : '▶ Play';
  if (running) raf = requestAnimationFrame(step);
  else cancelAnimationFrame(raf);
}

function resetAnim() {
  running = false;
  cancelAnimationFrame(raf);
  pos = 40;
  document.getElementById('ball').setAttribute('cx', 40);
  document.getElementById('play-btn').textContent = '▶ Play';
}
</script>
```

---

## Accordion

```html
<div style="max-width:560px">
  <div style="border:1px solid var(--color-border-tertiary);border-radius:var(--border-radius-lg);overflow:hidden">
    <!-- Item accordion -->
    <div style="border-bottom:1px solid var(--color-border-tertiary)">
      <button onclick="toggleAcc(this)"
        style="width:100%;padding:14px 16px;background:none;border:none;
               display:flex;justify-content:space-between;align-items:center;
               cursor:pointer;font-size:14px;font-weight:500;
               color:var(--color-text-primary);text-align:left">
        <span>Question / titre</span>
        <span class="acc-icon" style="transition:transform .2s">▾</span>
      </button>
      <div class="acc-body" style="display:none;padding:0 16px 14px;
           font-size:14px;color:var(--color-text-secondary);line-height:1.6">
        Réponse / contenu développé.
      </div>
    </div>
    <!-- Répéter pour chaque item -->
  </div>
</div>

<script>
function toggleAcc(btn) {
  const body = btn.nextElementSibling;
  const icon = btn.querySelector('.acc-icon');
  const open = body.style.display !== 'none';
  body.style.display = open ? 'none' : 'block';
  icon.style.transform = open ? '' : 'rotate(-180deg)';
}
</script>
```

---

## Progress bar animée

```html
<div style="max-width:560px">
  <div style="display:flex;justify-content:space-between;
              font-size:13px;margin-bottom:6px">
    <span style="font-weight:500;color:var(--color-text-primary)">Progression</span>
    <span id="pct-label" style="color:#534AB7;font-weight:600">0%</span>
  </div>
  <div style="height:8px;background:var(--color-background-secondary);
              border-radius:4px;overflow:hidden">
    <div id="progress-bar" style="height:100%;width:0%;background:#534AB7;
                                   border-radius:4px;transition:width .4s ease"></div>
  </div>
</div>

<script>
let pct = 0;
const interval = setInterval(() => {
  pct = Math.min(100, pct + Math.random() * 8 + 2);
  document.getElementById('progress-bar').style.width = pct + '%';
  document.getElementById('pct-label').textContent = Math.round(pct) + '%';
  if (pct >= 100) clearInterval(interval);
}, 200);
</script>
```

---

## Tooltip

```html
<span style="position:relative;display:inline-block"
      onmouseenter="this.querySelector('.tip').style.opacity=1"
      onmouseleave="this.querySelector('.tip').style.opacity=0">
  <span style="border-bottom:1px dashed var(--color-border-secondary);cursor:help;font-size:14px">
    Terme avec explication
  </span>
  <span class="tip" style="
    opacity:0;pointer-events:none;transition:opacity .15s;
    position:absolute;bottom:calc(100% + 6px);left:50%;transform:translateX(-50%);
    background:#2C2C2A;color:#F1EFE8;
    font-size:12px;padding:6px 10px;border-radius:6px;white-space:nowrap;z-index:10;">
    Explication courte ici
  </span>
</span>
```
