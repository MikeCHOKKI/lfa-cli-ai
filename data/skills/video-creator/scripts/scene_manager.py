"""Scene manager for timeline-based animation."""

import math
from easing import interpolate


class Scene:
    def __init__(self, name, start_sec, end_sec, draw_func):
        self.name = name
        self.start = start_sec
        self.end = end_sec
        self.duration = end_sec - start_sec
        self.draw = draw_func

    def is_active(self, time_sec):
        return self.start <= time_sec <= self.end

    def progress(self, time_sec):
        t = (time_sec - self.start) / max(self.duration, 0.001)
        return max(0, min(1, t))

    def local_time(self, time_sec):
        return time_sec - self.start


class SceneManager:
    def __init__(self, fps):
        self.fps = fps
        self.scenes = []
        self.total_duration = 0

    def add_scene(self, name, start_sec, end_sec, draw_func):
        self.scenes.append(Scene(name, start_sec, end_sec, draw_func))
        self.total_duration = max(self.total_duration, end_sec)

    @property
    def total_frames(self):
        return int(self.total_duration * self.fps)

    def render_frame(self, frame_idx, pixels, fw):
        time_sec = frame_idx / self.fps
        for scene in self.scenes:
            if scene.is_active(time_sec):
                t = scene.progress(time_sec)
                lt = scene.local_time(time_sec)
                scene.draw(pixels, fw, time_sec, t, lt)

    def get_active_scenes(self, time_sec):
        return [s for s in self.scenes if s.is_active(time_sec)]


class AnimValue:
    """Animated value with easing."""
    def __init__(self, start, end, delay=0, duration=1, easing="ease_out"):
        self.start_val = start
        self.end_val = end
        self.delay = delay
        self.duration = duration
        self.easing = easing

    def get(self, t):
        """t is scene-local progress 0..1"""
        local_t = max(0, min(1, (t - self.delay) / max(self.duration, 0.001)))
        return interpolate(self.start_val, self.end_val, local_t, self.easing)


class AnimSequence:
    """Chain of sequential animations."""
    def __init__(self):
        self.animations = []

    def then(self, start, end, delay=0, duration=1, easing="ease_out"):
        self.animations.append(AnimValue(start, end, delay, duration, easing))
        return self

    def get(self, t):
        for anim in reversed(self.animations):
            if t >= anim.delay:
                return anim.get(t)
        return self.animations[0].start_val if self.animations else 0
