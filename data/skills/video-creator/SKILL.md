---
name: video-creator
description: "Create videos from generative animations, image sequences, slideshows, or frame-by-frame animations. Use this skill whenever the user asks to create, generate, or edit video files (.mp4, .webm, .avi, .mov, .gif). Handles the full pipeline: generating frames with Python (stdlib only, no pip required) and encoding with FFmpeg. Supports procedural animation (shapes, particles, gradients), image sequence assembly, compositing, text overlays, transitions, audio tracks, and video filters."
license: Complete terms in LICENSE.txt
---

# Video Creator

A self-contained skill for creating videos using Python (stdlib) to generate frames as PPM/PAM images, and FFmpeg to encode them into video files. No external Python libraries are needed.

## Prerequisites

- `ffmpeg` must be available on the system
- `python3` must be available
- No pip or external Python libraries required

## Quick Start

```python
import sys, os
sys.path.insert(0, os.path.expanduser("~/.config/opencode/skills/video-creator/scripts"))
from frame_writer import FrameWriter

# Set your resolution and output format
RENDER_W, RENDER_H = 960, 540   # internal render resolution
OUTPUT_W, OUTPUT_H = 1920, 1080  # final output (upscaled)

with FrameWriter(RENDER_W, RENDER_H, fps=30, out_dir="/tmp/frames") as fw:
    for i in range(fw.total_frames(duration=5)):
        pixels = fw.create_buffer()
        fw.gradient_h(pixels, (10, 10, 30), (40, 20, 80))
        fw.draw_circle_sdf(pixels, RENDER_W//2, RENDER_H//2, 80, (200, 100, 255))
        fw.draw_line_aa(pixels, 0, 0, RENDER_W, RENDER_H, (255, 200, 100))
        fw.write_frame(pixels)

    fw.encode("output.mp4", output_format="mp4", crf=18,
              upscale=f"{OUTPUT_W}:{OUTPUT_H}")
```

## Workflow Overview

1. **Configure resolution & format** — set render resolution, output resolution, format (mp4/webm/mov/etc.)
2. **Generate frames** — draw into pixel buffers using bundled scripts
3. **Encode** — `fw.encode()` with upscale, filters, drawtext, audio
4. **Import images** (optional) — use `FrameWriter.import_image()` to load any image via FFmpeg

## Bundled Scripts

### `scripts/frame_writer.py`
Core frame generator with PPM/PAM output, SDF anti-aliasing, image import, audio.

```python
from frame_writer import FrameWriter

with FrameWriter(width=960, height=540, fps=30, out_dir="/tmp/frames",
                 bg_color=(0,0,0), use_alpha=False) as fw:
    for i in range(fw.total_frames(duration=5)):
        pixels = fw.create_buffer()
        # --- Full buffer ops ---
        fw.fill(pixels, (0, 0, 0))
        fw.gradient_h(pixels, (10,10,30), (40,20,80))
        fw.gradient_v(pixels, (100,0,0), (0,0,100))
        fw.radial_gradient(pixels, cx, cy, (255,255,255), (0,0,0))

        # --- Binary (fast) drawing ---
        fw.draw_rect(pixels, x1, y1, x2, y2, color)
        fw.draw_circle(pixels, cx, cy, r, color)
        fw.draw_line(pixels, x0, y0, x1, y1, color)
        fw.draw_polygon(pixels, points, color)
        fw.draw_rotated_rect(pixels, cx, cy, w, h, angle, color)

        # --- Anti-aliased (SDF) drawing ---
        fw.draw_circle_sdf(pixels, cx, cy, r, color)
        fw.draw_line_aa(pixels, x0, y0, x1, y1, color)
        fw.draw_rounded_rect_sdf(pixels, x1, y1, x2, y2, r, color)
        fw.draw_star_sdf(pixels, cx, cy, outer_r, inner_r, points=5, rotation=0, color=color)

        # --- Effects ---
        fw.noise(pixels, intensity=0.1, seed=42)
        fw.glow(pixels, cx, cy, radius, color, intensity=0.3)

        fw.write_frame(pixels)

    # --- Encoding with settings ---
    fw.encode(output="video.mp4",
              output_format="mp4",  # mp4, webm, avi, mov, gif, prores, h265
              crf=18,               # lower = better quality (0=lossless)
              preset="medium",      # ultrafast, fast, medium, slow, veryslow
              scale="1280:720",     # resize to this (optional)
              upscale="1920:1080",  # upscale to this with Lanczos (optional)
              drawtext={            # FFmpeg drawtext overlay (vector text)
                  'text': 'Hello',
                  'size': 48,
                  'color': 'white',
                  'font': '/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf',
                  'x': '(w-text_w)/2',
                  'y': '(h-text_h)/2',
              },
              extra_filters="hue=s=0")  # additional FFmpeg filters
```

### Image import (screenshots, logos, external images)
```python
from frame_writer import FrameWriter

# Method 1: Import any image directly (uses FFmpeg behind the scenes)
result = FrameWriter.import_image("screenshot.png")
if result:
    pixels, w, h = result
    # Now use `pixels` as a frame buffer or composite it

# Method 2: Convert any image to PPM (for batch processing)
path = FrameWriter.convert_image("logo.png", "/tmp/logo.ppm")

# Method 3: Read converted PPM files
result = FrameWriter.read_ppm("/tmp/logo.ppm")

# Method 4: Read PAM (RGBA)
result = FrameWriter.read_pam("/tmp/frame.pam")

# Composite imported image onto frame:
result = FrameWriter.import_image("overlay.png")
if result:
    overlay, ow, oh = result
    for y in range(min(oh, len(pixels))):
        for x in range(min(ow, len(pixels[0]))):
            if overlay[y][x] != (0, 0, 0):  # simple chroma key
                pixels[y][x] = overlay[y][x]
```

### Audio support
```python
# Generate a test tone
fw.generate_audio("audio.aac", freq=440, duration=5)

# Add audio to existing video
fw.add_audio("video.mp4", "audio.aac", "video_with_audio.mp4")
```

### `scripts/bitmap_font.py`
5×7 bitmap font for text in frames (A-Z, 0-9, basic punctuation).

```python
from bitmap_font import draw_text, draw_text_centered, text_width

draw_text(pixels, x, y, "Hello", color=(255, 255, 255), scale=2)
draw_text_centered(pixels, cx, y, "Centered", color=(200, 200, 255), scale=3)
w = text_width("Hello") * 2  # pre-calculate width at scale 2
```

### `scripts/easing.py`
Easing functions for smooth animation: linear, ease_in, ease_out, ease_in_out, ease_in_cubic, ease_out_cubic, ease_in_out_cubic, bounce_out, elastic_out, back_out.

```python
from easing import interpolate
value = interpolate(start=0, end=100, t=0.5, easing="bounce_out")
```

### `scripts/particles.py`
Particle system for explosions, sparkles, smoke.

```python
from particles import ParticleEmitter, make_star_points
emitter = ParticleEmitter(cx, cy, count=50, spread=100, speed=2, lifetime=30, gravity=0.05)
emitter.burst(30)
emitter.update()
for p in emitter.particles:
    alpha = p['life'] / p['max_life']
    fw.draw_circle_sdf(pixels, int(p['x']), int(p['y']), 3, (255, int(255*alpha), 0))
```

### `scripts/ui_elements.py`
Reusable UI element renderers for app mockups.

```python
from ui_elements import (
    draw_rounded_rect, draw_phone_frame, draw_card,
    draw_progress_bar, draw_button, draw_circle_icon
)

# Phone mockup
screen_bounds = draw_phone_frame(pixels, cx, cy, scale=1.0)

# Cards with shadows
draw_card(pixels, x1, y1, x2, y2, color=(255,255,255), shadow=True)

# Buttons with text
draw_button(pixels, x1, y1, x2, y2, (100,80,200), text="Click", text_color=(255,255,255))

# Progress bars
draw_progress_bar(pixels, x1, y1, x2, y2, progress=0.78, fill_color=(100,150,255))

# Circle icons
draw_circle_icon(pixels, cx, cy, r, (100,80,200), icon_type="check", icon_color=(255,255,255))
```

### `scripts/scene_manager.py`
Scene-based timeline animation.

```python
from scene_manager import SceneManager, AnimValue, Scene

sm = SceneManager(fps=30)
sm.add_scene("intro", 0, 3, lambda px, fw, t, prog, lt: ...)
sm.add_scene("main", 3, 8, lambda px, fw, t, prog, lt: ...)

# Animated values with easing
AnimValue(start=0, end=100, delay=0.2, duration=0.8, easing="ease_out").get(progress)

# In your render loop:
for i in range(sm.total_frames):
    pixels = fw.create_buffer()
    sm.render_frame(i, pixels, fw)
    fw.write_frame(pixels)
```

## Encoder Presets (output_format)

| Format | Codec | Pix_Fmt | Use Case |
|--------|-------|---------|----------|
| `mp4` | libx264 | yuv420p | Best compatibility (default) |
| `webm` | libvpx-vp9 | yuva420p | Open web format, smaller files |
| `avi` | mpeg4 | yuv420p | Legacy compatibility |
| `mov` | libx264 | yuv420p | Apple ecosystem |
| `gif` | gif | rgb24 | Animated GIF |
| `prores` | prores_ks | yuva444p10le | Pro video editing |
| `h265` | libx265 | yuv420p | Best compression, modern devices |

## Common Patterns

### Generate at lower res, upscale to HD/4K
```python
fw.encode("output.mp4", upscale="1920:1080")   # 960x540 → 1920x1080 via lanczos
fw.encode("output.mp4", upscale="3840:2160")   # 960x540 → 4K
```

### Vector text overlay (much sharper than bitmap font)
```python
fw.encode("output.mp4", drawtext={
    'text': 'Hello World',
    'font': '/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf',
    'size': 64,
    'color': 'white',
    'x': '(w-text_w)/2',
    'y': '(h-text_h)/2',
    'box': 1,
    'box_color': 'black@0.5',
    'box_border': 8,
})
```

### Import and composite a logo
```python
result = FrameWriter.import_image("logo.png")
if result:
    logo, lw, lh = result
    for frame in frames:
        pixels = fw.create_buffer()
        # draw scene ...
        # overlay logo at top-left
        for y in range(min(lh, H)):
            for x in range(min(lw, W)):
                r, g, b = logo[y][x]
                if (r, g, b) != (0, 0, 0):  # skip black background
                    pixels[y][x] = logo[y][x]
        fw.write_frame(pixels)
```

## Advanced FFmpeg Filters

### Add to `extra_filters` parameter:
```python
fw.encode("output.mp4", extra_filters="hue=s=1.5")           # saturation
fw.encode("output.mp4", extra_filters="eq=brightness=0.05:contrast=1.1")  # color grade
fw.encode("output.mp4", extra_filters="unsharp=3:3:0.5")     # sharpen
fw.encode("output.mp4", extra_filters="fade=t=in:st=0:d=1")  # fade in
```

## Performance Tips

- Render at 960×540 or lower, upscale to 1080p/4K with `upscale=` parameter
- Use binary drawing (`draw_circle`, `draw_line`) for large simple shapes
- Use SDF drawing (`draw_circle_sdf`, etc.) for small detailed elements
- Batch similar FFmpeg filters in `extra_filters` (single FFmpeg pass)
- Set `crf=28` for previews, `crf=18` for final, `crf=0` for lossless

## Guidelines

- Always clean up temp frames (`cleanup=True` by default)
- Use `-pix_fmt yuv420p` for maximum compatibility
- For fast previews: lower resolution and CRF 28
- For final renders: full resolution and CRF 18 or lower
- Put temp directories in `/tmp/` unless user asks otherwise
- Verify output with ffprobe

## Philosophy

This skill provides:
- **Knowledge**: Video creation pipeline, FFmpeg, SDF anti-aliasing
- **Utilities**: FrameWriter, easing, particles, UI elements, scene manager
- **Self-containment**: Python stdlib + FFmpeg only

It does NOT provide:
- Pre-rendered templates or stock footage
- GUI video editing
- AI upscaling / frame interpolation (unless user provides external tools)
