import struct
import os
import math
import shutil
import subprocess
import re


class FrameWriter:
    """Write PPM/PAM frames to disk for FFmpeg encoding.

    Usage:
        with FrameWriter(width=1920, height=1080, fps=30, out_dir="frames") as fw:
            for frame_idx in range(fw.total_frames(duration=5)):
                pixels = fw.create_buffer()
                fw.draw_circle_sdf(pixels, cx, cy, r, color)
                fw.write_frame(pixels)
            fw.encode("output.mp4", output_format="mp4")
    """

    FORMAT_PRESETS = {
        "mp4":  {"codec": "libx264",       "pix_fmt": "yuv420p", "ext": ".mp4"},
        "webm": {"codec": "libvpx-vp9",     "pix_fmt": "yuva420p", "ext": ".webm"},
        "avi":  {"codec": "mpeg4",          "pix_fmt": "yuv420p", "ext": ".avi"},
        "mov":  {"codec": "libx264",        "pix_fmt": "yuv420p", "ext": ".mov"},
        "gif":  {"codec": "gif",            "pix_fmt": "rgb24",   "ext": ".gif"},
        "prores": {"codec": "prores_ks",    "pix_fmt": "yuva444p10le", "ext": ".mov"},
        "h265": {"codec": "libx265",        "pix_fmt": "yuv420p", "ext": ".mp4"},
    }

    def __init__(self, width=1920, height=1080, fps=30, out_dir="frames",
                 bg_color=(0, 0, 0), clear_before=True, use_alpha=False):
        self.width = width
        self.height = height
        self.fps = fps
        self.out_dir = out_dir
        self.bg_color = bg_color
        self.use_alpha = use_alpha
        self.frame_num = 0

        if clear_before and os.path.exists(out_dir):
            shutil.rmtree(out_dir)

    def __enter__(self):
        os.makedirs(self.out_dir, exist_ok=True)
        return self

    def __exit__(self, *args):
        pass

    def set_resolution(self, width, height):
        self.width = width
        self.height = height

    def total_frames(self, duration=None, frames=None):
        if frames is not None:
            return frames
        return int(self.fps * duration)

    def create_buffer(self):
        return [[self.bg_color for _ in range(self.width)] for _ in range(self.height)]

    # --- Buffer operations ---

    def fill(self, pixels, color):
        for y in range(self.height):
            for x in range(self.width):
                pixels[y][x] = color

    def gradient_h(self, pixels, color_top, color_bottom):
        for y in range(self.height):
            t = y / max(self.height - 1, 1)
            r = int(color_top[0] * (1 - t) + color_bottom[0] * t)
            g = int(color_top[1] * (1 - t) + color_bottom[1] * t)
            b = int(color_top[2] * (1 - t) + color_bottom[2] * t)
            for x in range(self.width):
                pixels[y][x] = (r, g, b)

    def gradient_v(self, pixels, color_left, color_right):
        for y in range(self.height):
            for x in range(self.width):
                t = x / max(self.width - 1, 1)
                r = int(color_left[0] * (1 - t) + color_right[0] * t)
                g = int(color_left[1] * (1 - t) + color_right[1] * t)
                b = int(color_left[2] * (1 - t) + color_right[2] * t)
                pixels[y][x] = (r, g, b)

    def radial_gradient(self, pixels, cx, cy, color_center, color_edge, max_r=None):
        if max_r is None:
            max_r = math.sqrt(self.width ** 2 + self.height ** 2)
        for y in range(self.height):
            for x in range(self.width):
                dist = math.sqrt((x - cx) ** 2 + (y - cy) ** 2)
                t = min(1.0, dist / max_r)
                r = int(color_center[0] * (1 - t) + color_edge[0] * t)
                g = int(color_center[1] * (1 - t) + color_edge[1] * t)
                b = int(color_center[2] * (1 - t) + color_edge[2] * t)
                pixels[y][x] = (r, g, b)

    def draw_pixel(self, pixels, x, y, color):
        if 0 <= y < self.height and 0 <= x < self.width:
            pixels[y][x] = color

    def draw_rect(self, pixels, x1, y1, x2, y2, color, fill=True):
        x1, x2 = max(0, int(x1)), min(self.width, int(x2))
        y1, y2 = max(0, int(y1)), min(self.height, int(y2))
        for y in range(y1, y2):
            for x in range(x1, x2):
                pixels[y][x] = color

    # --- Binary (fast) drawing ---

    def draw_circle(self, pixels, cx, cy, r, color):
        r2 = r * r
        for y in range(max(0, cy - r), min(self.height, cy + r + 1)):
            for x in range(max(0, cx - r), min(self.width, cx + r + 1)):
                dx, dy = x - cx, y - cy
                if dx * dx + dy * dy <= r2:
                    pixels[y][x] = color

    def draw_line(self, pixels, x0, y0, x1, y1, color):
        dx, dy = abs(x1 - x0), -abs(y1 - y0)
        sx = 1 if x0 < x1 else -1
        sy = 1 if y0 < y1 else -1
        err = dx + dy
        while True:
            if 0 <= y0 < self.height and 0 <= x0 < self.width:
                pixels[y0][x0] = color
            if x0 == x1 and y0 == y1:
                break
            e2 = 2 * err
            if e2 >= dy:
                err += dy; x0 += sx
            if e2 <= dx:
                err += dx; y0 += sy

    def draw_polygon(self, pixels, points, color):
        if not points:
            return
        ys = [p[1] for p in points]
        min_y, max_y = int(max(0, min(ys))), int(min(self.height - 1, max(ys)))
        for y in range(min_y, max_y + 1):
            intersections = []
            n = len(points)
            for i in range(n):
                x1, y1 = points[i]
                x2, y2 = points[(i + 1) % n]
                if (y1 <= y < y2) or (y2 <= y < y1):
                    t = (y - y1) / (y2 - y1) if y2 != y1 else 0
                    x = x1 + t * (x2 - x1)
                    intersections.append(x)
            intersections.sort()
            for i in range(0, len(intersections) - 1, 2):
                x_start = int(max(0, intersections[i]))
                x_end = int(min(self.width - 1, intersections[i + 1]))
                for x in range(x_start, x_end + 1):
                    pixels[y][x] = color

    def draw_rotated_rect(self, pixels, cx, cy, w, h, angle, color):
        cos_a = math.cos(angle)
        sin_a = math.sin(angle)
        half_w, half_h = w / 2, h / 2
        corners = [
            (cx + (-half_w) * cos_a - (-half_h) * sin_a,
             cy + (-half_w) * sin_a + (-half_h) * cos_a),
            (cx + half_w * cos_a - (-half_h) * sin_a,
             cy + half_w * sin_a + (-half_h) * cos_a),
            (cx + half_w * cos_a - half_h * sin_a,
             cy + half_w * sin_a + half_h * cos_a),
            (cx + (-half_w) * cos_a - half_h * sin_a,
             cy + (-half_w) * sin_a + half_h * cos_a),
        ]
        self.draw_polygon(pixels, corners, color)

    # --- SDF Anti-aliased drawing ---

    def blend(self, bg, fg, alpha):
        a = max(0.0, min(1.0, alpha))
        if a >= 1:
            return fg
        if a <= 0:
            return bg
        return (
            int(bg[0] * (1 - a) + fg[0] * a),
            int(bg[1] * (1 - a) + fg[1] * a),
            int(bg[2] * (1 - a) + fg[2] * a),
        )

    def draw_circle_sdf(self, pixels, cx, cy, r, color):
        cx_f, cy_f = cx + 0.0, cy + 0.0
        x1 = max(0, int(cx - r - 1))
        x2 = min(self.width, int(cx + r + 2))
        y1 = max(0, int(cy - r - 1))
        y2 = min(self.height, int(cy + r + 2))
        for y in range(y1, y2):
            for x in range(x1, x2):
                dist = math.sqrt((x + 0.5 - cx_f) ** 2 + (y + 0.5 - cy_f) ** 2) - r
                alpha = max(0.0, min(1.0, 0.5 - dist))
                if alpha > 0:
                    pixels[y][x] = self.blend(pixels[y][x], color, alpha)

    def draw_line_aa(self, pixels, x0, y0, x1, y1, color):
        """Xiaolin Wu anti-aliased line algorithm."""
        dx = abs(x1 - x0)
        dy = abs(y1 - y0)
        steep = dy > dx

        if steep:
            x0, y0 = y0, x0
            x1, y1 = y1, x1

        if x0 > x1:
            x0, x1 = x1, x0
            y0, y1 = y1, y0

        dx = x1 - x0
        dy = y1 - y0
        gradient = dy / dx if dx != 0 else 1.0

        xend = round(x0)
        yend = y0 + gradient * (xend - x0)
        xgap = 1 - (x0 + 0.5 - int(x0))
        xpxl1 = xend
        ypxl1 = int(yend)
        self._plot_aa(pixels, xpxl1, ypxl1, color, (1 - (yend - int(yend))) * xgap, steep)
        self._plot_aa(pixels, xpxl1, ypxl1 + 1, color, (yend - int(yend)) * xgap, steep)

        intery = yend + gradient

        xend = round(x1)
        yend = y1 + gradient * (xend - x1)
        xgap = x1 + 0.5 - (x1 - int(x1))
        xpxl2 = xend
        ypxl2 = int(yend)
        self._plot_aa(pixels, xpxl2, ypxl2, color, (1 - (yend - int(yend))) * xgap, steep)
        self._plot_aa(pixels, xpxl2, ypxl2 + 1, color, (yend - int(yend)) * xgap, steep)

        for x in range(xpxl1 + 1, xpxl2):
            self._plot_aa(pixels, x, int(intery), color, 1 - (intery - int(intery)), steep)
            self._plot_aa(pixels, x, int(intery) + 1, color, intery - int(intery), steep)
            intery += gradient

    def _plot_aa(self, pixels, x, y, color, alpha, steep):
        if alpha <= 0:
            return
        if steep:
            px, py = y, x
        else:
            px, py = x, y
        if 0 <= py < self.height and 0 <= px < self.width:
            pixels[py][px] = self.blend(pixels[py][px], color, min(1.0, alpha))

    def draw_rounded_rect_sdf(self, pixels, x1, y1, x2, y2, r, color):
        """Anti-aliased rounded rectangle using SDF."""
        x1, y1, x2, y2 = float(x1), float(y1), float(x2), float(y2)
        r = min(r, (x2 - x1) / 2, (y2 - y1) / 2)
        half = 0.5
        x_start = max(0, int(x1 - 1))
        x_end = min(self.width, int(x2 + 2))
        y_start = max(0, int(y1 - 1))
        y_end = min(self.height, int(y2 + 2))

        for y in range(y_start, y_end):
            for x in range(x_start, x_end):
                px, py = x + 0.5, y + 0.5
                dx = 0.0
                if px < x1 + r:
                    dx = px - (x1 + r)
                elif px > x2 - r:
                    dx = px - (x2 - r)
                dy = 0.0
                if py < y1 + r:
                    dy = py - (y1 + r)
                elif py > y2 - r:
                    dy = py - (y2 - r)

                if dx == 0 and dy == 0:
                    pixels[y][x] = color
                else:
                    dist = math.sqrt(dx * dx + dy * dy) - r
                    alpha = max(0.0, min(1.0, 0.5 - dist))
                    if alpha > 0:
                        pixels[y][x] = self.blend(pixels[y][x], color, alpha)

    def draw_star_sdf(self, pixels, cx, cy, outer_r, inner_r, points=5, rotation=0, color=None):
        """Anti-aliased star using SDF."""
        if color is None:
            color = (255, 255, 255)
        pts = []
        for i in range(points * 2):
            angle = rotation + i * math.pi / points - math.pi / 2
            r = outer_r if i % 2 == 0 else inner_r
            pts.append((cx + r * math.cos(angle), cy + r * math.sin(angle)))

        x_start = max(0, int(cx - outer_r - 1))
        x_end = min(self.width, int(cx + outer_r + 2))
        y_start = max(0, int(cy - outer_r - 1))
        y_end = min(self.height, int(cy + outer_r + 2))

        n = len(pts)
        for y in range(y_start, y_end):
            for x in range(x_start, x_end):
                px, py = x + 0.5, y + 0.5
                inside = False
                for i in range(n):
                    x1, y1 = pts[i]
                    x2, y2 = pts[(i + 1) % n]
                    if ((y1 > py) != (y2 > py)) and (px < (x2 - x1) * (py - y1) / (y2 - y1) + x1):
                        inside = not inside

                if inside:
                    min_dist = float('inf')
                    for i in range(n):
                        x1, y1 = pts[i]
                        x2, y2 = pts[(i + 1) % n]
                        ex, ey = x2 - x1, y2 - y1
                        t = max(0, min(1, ((px - x1) * ex + (py - y1) * ey) / (ex * ex + ey * ey)))
                        d = math.sqrt((px - (x1 + t * ex)) ** 2 + (py - (y1 + t * ey)) ** 2)
                        min_dist = min(min_dist, d)
                    alpha = max(0.0, min(1.0, 0.5 - (min_dist - 0.5)))
                    if alpha > 0:
                        pixels[y][x] = self.blend(pixels[y][x], color, alpha)

    # --- Effects ---

    def noise(self, pixels, intensity=0.1, seed=None):
        import random
        if seed is not None:
            random.seed(seed)
        for y in range(self.height):
            for x in range(self.width):
                if random.random() < intensity:
                    n = random.randint(-30, 30)
                    px = pixels[y][x]
                    r = max(0, min(255, px[0] + n))
                    g = max(0, min(255, px[1] + n))
                    b = max(0, min(255, px[2] + n))
                    pixels[y][x] = (r, g, b)

    def glow(self, pixels, cx, cy, radius, color, intensity=0.3):
        """Add a radial glow effect around a point."""
        for y in range(max(0, cy - radius), min(self.height, cy + radius + 1)):
            for x in range(max(0, cx - radius), min(self.width, cx + radius + 1)):
                dist = math.sqrt((x - cx) ** 2 + (y - cy) ** 2)
                if dist <= radius:
                    alpha = intensity * (1 - dist / radius)
                    pixels[y][x] = self.blend(pixels[y][x], color, alpha)

    # --- Image import via FFmpeg ---

    @staticmethod
    def import_image(path, width=None, height=None):
        """Import any image format as pixel buffer using FFmpeg conversion.

        Supports PNG, JPEG, BMP, WebP, GIF, SVG, etc. (anything FFmpeg can read).
        Returns (pixels, img_width, img_height) or None on failure.
        """
        if not os.path.exists(path):
            return None

        ppm_path = path + ".import.tmp.ppm"
        try:
            scale = ""
            if width and height:
                scale = f" -vf scale={width}:{height}"
            subprocess.run(
                f"ffmpeg -y -i \"{path}\"{scale} \"{ppm_path}\"",
                shell=True, capture_output=True, timeout=30
            )
            if not os.path.exists(ppm_path):
                return None
            return FrameWriter.read_ppm(ppm_path)
        finally:
            if os.path.exists(ppm_path):
                os.remove(ppm_path)

    @staticmethod
    def read_ppm(path):
        """Read a PPM file into a pixel buffer. Returns (pixels, w, h) or None."""
        try:
            with open(path, 'rb') as f:
                header = f.readline().strip()
                if header != b'P6':
                    return None
                line = f.readline().strip()
                while line.startswith(b'#'):
                    line = f.readline().strip()
                parts = line.split()
                w, h = int(parts[0]), int(parts[1])
                maxval = int(f.readline().strip())
                data = f.read()
                pixels = [[(0, 0, 0) for _ in range(w)] for _ in range(h)]
                stride = 3
                for y in range(h):
                    for x in range(w):
                        idx = (y * w + x) * stride
                        if idx + 2 < len(data):
                            r = data[idx] * 255 // maxval
                            g = data[idx + 1] * 255 // maxval
                            b = data[idx + 2] * 255 // maxval
                            pixels[y][x] = (r, g, b)
                return (pixels, w, h)
        except Exception:
            return None

    @staticmethod
    def read_pam(path):
        """Read a PAM (RGBA) file. Returns (pixels, w, h) or None."""
        try:
            with open(path, 'rb') as f:
                magic = f.readline().strip()
                if magic != b'P7':
                    return None
                w, h, depth, maxval = 0, 0, 3, 255
                for _ in range(20):
                    line = f.readline().strip()
                    if line == b'ENDHDR':
                        break
                    parts = line.split()
                    if parts[0] == b'WIDTH':
                        w = int(parts[1])
                    elif parts[0] == b'HEIGHT':
                        h = int(parts[1])
                    elif parts[0] == b'DEPTH':
                        depth = int(parts[1])
                    elif parts[0] == b'MAXVAL':
                        maxval = int(parts[1])
                data = f.read()
                pixels = [[(0, 0, 0) for _ in range(w)] for _ in range(h)]
                for y in range(h):
                    for x in range(w):
                        idx = (y * w + x) * depth
                        if idx + 2 < len(data):
                            r = data[idx] * 255 // maxval
                            g = data[idx + 1] * 255 // maxval
                            b = data[idx + 2] * 255 // maxval
                            pixels[y][x] = (r, g, b)
                return (pixels, w, h)
        except Exception:
            return None

    @staticmethod
    def convert_image(input_path, output_path=None, width=None, height=None):
        """Convert any image to PPM via FFmpeg. Returns path to PPM or None."""
        if output_path is None:
            output_path = input_path + ".converted.ppm"
        scale = ""
        if width and height:
            scale = f" -vf scale={width}:{height}"
        try:
            r = subprocess.run(
                f"ffmpeg -y -i \"{input_path}\"{scale} \"{output_path}\"",
                shell=True, capture_output=True, timeout=30
            )
            if r.returncode == 0 and os.path.exists(output_path):
                return output_path
            return None
        except Exception:
            return None

    # --- Frame writing ---

    def write_frame(self, pixels):
        if self.use_alpha:
            self.write_frame_pam(pixels)
        else:
            self.write_frame_ppm(pixels)

    def write_frame_ppm(self, pixels):
        path = os.path.join(self.out_dir, f"frame_{self.frame_num + 1:04d}.ppm")
        self.frame_num += 1
        with open(path, 'wb') as f:
            f.write(f'P6\n{self.width} {self.height}\n255\n'.encode())
            for y in range(self.height):
                for x in range(self.width):
                    r, g, b = pixels[y][x][:3]
                    f.write(struct.pack('BBB',
                           max(0, min(255, int(r))),
                           max(0, min(255, int(g))),
                           max(0, min(255, int(b)))))

    def write_frame_pam(self, pixels):
        """Write frame as PAM with alpha (4 channels)."""
        path = os.path.join(self.out_dir, f"frame_{self.frame_num + 1:04d}.pam")
        self.frame_num += 1
        with open(path, 'wb') as f:
            f.write(f'P7\nWIDTH {self.width}\nHEIGHT {self.height}\nDEPTH 4\nMAXVAL 255\nTUPLTYPE RGB_ALPHA\nENDHDR\n'.encode())
            for y in range(self.height):
                for x in range(self.width):
                    px = pixels[y][x]
                    r, g, b = px[0], px[1], px[2]
                    a = px[3] if len(px) > 3 else 255
                    f.write(struct.pack('BBBB',
                           max(0, min(255, int(r))),
                           max(0, min(255, int(g))),
                           max(0, min(255, int(b))),
                           max(0, min(255, int(a)))))

    def cleanup(self):
        if os.path.exists(self.out_dir):
            shutil.rmtree(self.out_dir)

    # --- Encoding ---

    def encode(self, output="output.mp4", output_format="mp4", crf=18, preset="medium",
               framerate=None, cleanup=True, extra_filters="", scale=None,
               drawtext=None, upscale=None):
        """Encode frames to video.

        Args:
            output: Output file path (or None to auto-generate from format).
            output_format: One of FORMAT_PRESETS keys ('mp4','webm','avi','mov','gif','prores','h265').
            crf: Quality (lower = better, 0=lossless for h264).
            preset: FFmpeg preset (ultrafast, fast, medium, slow, veryslow).
            framerate: Override FPS.
            cleanup: Remove frame directory after encoding.
            extra_filters: Additional FFmpeg video filters.
            scale: Resize output (e.g. "1920:1080").
            drawtext: Dict with text overlay params (see _build_drawtext).
            upscale: Upscale resolution (e.g. "1920:1080" uses lanczos).
        """
        fmt = self.FORMAT_PRESETS.get(output_format)
        if fmt is None:
            fmt = self.FORMAT_PRESETS["mp4"]

        if output is None:
            output = f"output{fmt['ext']}"

        codec = fmt["codec"]
        pix_fmt = fmt["pix_fmt"]
        fps = framerate or self.fps

        ext = ".pam" if self.use_alpha else ".ppm"
        pattern = f"{self.out_dir}/frame_%04d{ext}"
        input_type = "pam" if self.use_alpha else "ppm"

        filters = []
        if scale:
            filters.append(f"scale={scale}:flags=lanczos")
        if upscale:
            filters.append(f"scale={upscale}:flags=lanczos")
        if drawtext:
            dt = self._build_drawtext(drawtext)
            if dt:
                filters.append(dt)
        if extra_filters:
            filters.append(extra_filters)

        vf = ""
        if filters:
            vf = " -vf \"" + ",".join(filters) + "\""

        cmd = (
            f"ffmpeg -y -framerate {fps} -f image2 -pattern_type sequence"
            f" -i \"{pattern}\""
            f" -c:v {codec} -pix_fmt {pix_fmt} -preset {preset} -crf {crf}"
            f"{vf} \"{output}\""
        )

        ret = os.system(cmd)
        if cleanup:
            self.cleanup()
        return ret, output

    def _build_drawtext(self, dt):
        """Build FFmpeg drawtext filter from dict.

        dt = {
            'text': 'Hello',
            'font': '/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf',
            'size': 48,
            'color': 'white',
            'x': '(w-text_w)/2',
            'y': '(h-text_h)/2',
            'box': 1,
            'box_color': 'black@0.5',
            'box_border': 2,
            'time': None,       # tuple (start_sec, end_sec) or None
        }
        """
        if not dt or not dt.get('text'):
            return None

        font = dt.get('font', '/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf')
        size = dt.get('size', 48)
        color = dt.get('color', 'white')
        x = dt.get('x', '(w-text_w)/2')
        y = dt.get('y', '(h-text_h)/2')

        parts = [
            f"drawtext=text='{dt['text']}'",
            f"fontfile={font}",
            f"fontsize={size}",
            f"fontcolor={color}",
            f"x={x}",
            f"y={y}",
        ]

        if dt.get('box'):
            parts.append(f"box={dt['box']}")
            parts.append(f"boxcolor={dt.get('box_color', 'black@0.5')}")
            if dt.get('box_border'):
                parts.append(f"boxborderw={dt['box_border']}")

        if dt.get('enable'):
            parts.append(f"enable='{dt['enable']}'")

        return ':'.join(parts)

    def add_audio(self, video_path, audio_path, output_path=None, volume=1.0):
        """Add/replace audio track in video."""
        if output_path is None:
            base, ext = os.path.splitext(video_path)
            output_path = f"{base}_with_audio{ext}"
        cmd = (
            f"ffmpeg -y -i \"{video_path}\" -i \"{audio_path}\""
            f" -c:v copy -c:a aac -map 0:v:0 -map 1:a:0"
            f" -af volume={volume}"
            f" -shortest \"{output_path}\""
        )
        ret = os.system(cmd)
        return ret, output_path

    def generate_audio(self, output_path="audio.aac", freq=440, duration=5, type="sine"):
        """Generate a test tone audio file using FFmpeg."""
        if type == "sine":
            cmd = f"ffmpeg -y -f lavfi -i \"sine=frequency={freq}:duration={duration}\" -c:a aac \"{output_path}\""
        elif type == "noise":
            cmd = f"ffmpeg -y -f lavfi -i \"anoisesrc=d={duration}\" -c:a aac \"{output_path}\""
        else:
            return 1, None
        ret = os.system(cmd)
        return ret, output_path
