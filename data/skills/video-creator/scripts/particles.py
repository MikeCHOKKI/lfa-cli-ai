"""Simple particle system for video animations.

Usage:
    emitter = ParticleEmitter(cx, cy, count=100)
    for frame in range(total_frames):
        pixels = fw.create_buffer()
        emitter.update()
        for p in emitter.particles:
            alpha = p['life'] / p['max_life']
            color = (255, int(255 * alpha), int(255 * alpha * 0.5))
            fw.draw_pixel(pixels, int(p['x']), int(p['y']), color)
"""

import math
import random


class ParticleEmitter:
    def __init__(self, cx, cy, count=50, spread=100, speed=2, lifetime=30,
                 gravity=0.05, size=2):
        self.cx = cx
        self.cy = cy
        self.count = count
        self.spread = spread
        self.speed = speed
        self.lifetime = lifetime
        self.gravity = gravity
        self.size = size
        self.particles = []

    def burst(self, count=None):
        n = count or self.count
        for _ in range(n):
            angle = random.uniform(0, 2 * math.pi)
            vel = random.uniform(0.5, 1.5) * self.speed
            self.particles.append({
                'x': self.cx + random.uniform(-self.spread * 0.1, self.spread * 0.1),
                'y': self.cy + random.uniform(-self.spread * 0.1, self.spread * 0.1),
                'vx': math.cos(angle) * vel * random.uniform(0.5, 1.5),
                'vy': math.sin(angle) * vel * random.uniform(0.5, 1.5),
                'life': self.lifetime,
                'max_life': self.lifetime,
            })

    def update(self):
        self.particles = [p for p in self.particles if p['life'] > 0]
        for p in self.particles:
            p['x'] += p['vx']
            p['y'] += p['vy']
            p['vy'] += self.gravity
            p['vx'] *= 0.99  # drag
            p['life'] -= 1

    def emit_continuous(self):
        self.burst(max(1, self.count // self.lifetime))

    def clear(self):
        self.particles = []


def make_star_points(cx, cy, outer_r, inner_r, points=5, rotation=0):
    """Generate points for a star polygon."""
    pts = []
    for i in range(points * 2):
        angle = rotation + i * math.pi / points - math.pi / 2
        r = outer_r if i % 2 == 0 else inner_r
        pts.append((cx + r * math.cos(angle), cy + r * math.sin(angle)))
    return pts
