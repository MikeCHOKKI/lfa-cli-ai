"""UI element renderers for app mockups in PPM frames."""

import math


def draw_rounded_rect(pixels, x1, y1, x2, y2, r, color, fill=True, border_w=0, border_color=None):
    """Draw a rounded rectangle using SDF approach."""
    h, w = len(pixels), len(pixels[0])
    x1, x2 = int(x1), int(x2)
    y1, y2 = int(y1), int(y2)
    r = min(r, (x2 - x1) // 2, (y2 - y1) // 2)
    r2 = r * r

    for y in range(max(0, y1), min(h, y2 + 1)):
        for x in range(max(0, x1), min(w, x2 + 1)):
            # Determine distance to rectangle border
            dx = 0
            if x < x1 + r:
                dx = x - (x1 + r)
            elif x > x2 - r:
                dx = x - (x2 - r)

            dy = 0
            if y < y1 + r:
                dy = y - (y1 + r)
            elif y > y2 - r:
                dy = y - (y2 - r)

            inside_corner = dx * dx + dy * dy <= r2
            in_corner = dx != 0 or dy != 0

            if not in_corner or (in_corner and inside_corner):
                if fill:
                    pixels[y][x] = color
            elif border_w > 0 and border_color:
                # Check if pixel is within border width of edge
                dist_to_edge = min(
                    y - y1, y2 - y, x - x1, x2 - x,
                    math.sqrt((x - x1) ** 2 + (y - y1) ** 2) if y - y1 < r and x - x1 < r else 999,
                    math.sqrt((x - x2) ** 2 + (y - y1) ** 2) if y - y1 < r and x2 - x < r else 999,
                    math.sqrt((x - x1) ** 2 + (y - y2) ** 2) if y2 - y < r and x - x1 < r else 999,
                    math.sqrt((x - x2) ** 2 + (y - y2) ** 2) if y2 - y < r and x2 - x < r else 999,
                )
                if dist_to_edge <= border_w and dist_to_edge >= 0:
                    pixels[y][x] = border_color


def draw_phone_frame(pixels, cx, cy, scale=1.0, color=(30, 30, 40), screen_color=(240, 242, 248),
                     border_color=(60, 60, 80)):
    """Draw a phone frame centered at (cx, cy)."""
    phone_w = int(180 * scale)
    phone_h = int(360 * scale)
    r = int(16 * scale)

    x1 = cx - phone_w // 2
    y1 = cy - phone_h // 2

    # Phone body
    draw_rounded_rect(pixels, x1, y1, x1 + phone_w, y1 + phone_h, r, color, fill=True)

    # Screen
    margin = int(8 * scale)
    draw_rounded_rect(pixels, x1 + margin, y1 + margin,
                     x1 + phone_w - margin, y1 + phone_h - margin,
                     int(8 * scale), screen_color, fill=True)

    # Notch / Dynamic island
    notch_w = int(60 * scale)
    notch_h = int(8 * scale)
    draw_rounded_rect(pixels, cx - notch_w // 2, y1 + margin + int(3 * scale),
                     cx + notch_w // 2, y1 + margin + notch_h + int(3 * scale),
                     int(4 * scale), color, fill=True)

    # Bottom bar indicator
    bar_w = int(40 * scale)
    bar_h = int(3 * scale)
    bar_y = y1 + phone_h - margin - int(6 * scale)
    draw_rounded_rect(pixels, cx - bar_w // 2, bar_y, cx + bar_w // 2, bar_y + bar_h,
                     int(2 * scale), border_color, fill=True)

    return {
        'x1': x1 + margin, 'y1': y1 + margin,
        'x2': x1 + phone_w - margin, 'y2': y1 + phone_h - margin
    }


def draw_card(pixels, x1, y1, x2, y2, color=(255, 255, 255), shadow=True, r=8):
    """Draw a card with optional shadow."""
    if shadow:
        shadow_color = (200, 200, 210)
        draw_rounded_rect(pixels, x1 + 2, y1 + 2, x2 + 2, y2 + 2, r, shadow_color, fill=True)
    draw_rounded_rect(pixels, x1, y1, x2, y2, r, color, fill=True)
    return (x1, y1, x2, y2)


def draw_progress_bar(pixels, x1, y1, x2, y2, progress, bg_color=(220, 220, 230),
                      fill_color=(100, 150, 255), r=6):
    """Draw a progress bar."""
    draw_rounded_rect(pixels, x1, y1, x2, y2, r, bg_color, fill=True)
    fill_w = int((x2 - x1) * max(0, min(1, progress)))
    if fill_w > 0:
        draw_rounded_rect(pixels, x1, y1, x1 + fill_w, y2, r, fill_color, fill=True)


def draw_circle_icon(pixels, cx, cy, r, color, icon_type="check", icon_color=(255, 255, 255)):
    """Draw a circle icon with a simple symbol inside."""
    # Circle background
    draw_rounded_rect(pixels, cx - r, cy - r, cx + r, cy + r, r, color, fill=True)

    if icon_type == "check":
        # Simple checkmark using Bresenham line
        pts = [(cx - r // 3, cy), (cx - r // 6, cy + r // 3), (cx + r // 2, cy - r // 3)]
        for i in range(len(pts) - 1):
            x0, y0 = pts[i]
            x1, y1 = pts[i + 1]
            dx, dy = abs(x1 - x0), -abs(y1 - y0)
            sx = 1 if x0 < x1 else -1
            sy = 1 if y0 < y1 else -1
            err = dx + dy
            px, py = int(x0), int(y0)
            while True:
                if 0 <= py < len(pixels) and 0 <= px < len(pixels[0]):
                    pixels[py][px] = icon_color
                if px == int(x1) and py == int(y1):
                    break
                e2 = 2 * err
                if e2 >= dy:
                    err += dy; px += sx
                if e2 <= dx:
                    err += dx; py += sy

    elif icon_type == "heart":
        sr = r * 0.4
        for a in range(0, 360, 5):
            rad = math.radians(a)
            hx = cx + sr * 16 * (math.sin(rad) ** 3)
            hy = cy - sr * (13 * math.cos(rad) - 5 * math.cos(2 * rad) - 2 * math.cos(3 * rad) - math.cos(4 * rad))
            x0, y0 = int(hx), int(hy)
            if 0 <= y0 < len(pixels) and 0 <= x0 < len(pixels[0]):
                pixels[y0][x0] = icon_color


def draw_bar_chart(pixels, x1, y1, x2, y2, values, bar_color=(100, 180, 255),
                   bg_color=(230, 235, 245)):
    """Draw a simple bar chart."""
    h, w = y2 - y1, x2 - x1
    draw_rounded_rect(pixels, x1, y1, x2, y2, 4, bg_color, fill=True)

    n = len(values)
    max_val = max(values) if values else 1
    bar_w = max(1, w // n - 4)
    gap = (w - bar_w * n) // (n + 1)

    for i, v in enumerate(v):
        bar_h = int((v / max_val) * (h - 8))
        bx = x1 + gap + i * (bar_w + gap)
        by = y2 - 4 - bar_h
        draw_rounded_rect(pixels, bx, by, bx + bar_w, y2 - 4, 2, bar_color, fill=True)


def draw_button(pixels, x1, y1, x2, y2, color, text="", text_color=(255, 255, 255),
                text_scale=1, r=6):
    """Draw a button with optional text."""
    draw_rounded_rect(pixels, x1, y1, x2, y2, r, color, fill=True)
    if text:
        from bitmap_font import draw_text_centered
        cy = (y1 + y2) // 2 - (4 * text_scale) // 2
        draw_text_centered(pixels, (x1 + x2) // 2, cy, text, text_color, scale=text_scale)
