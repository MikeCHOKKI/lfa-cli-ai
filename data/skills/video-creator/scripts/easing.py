"""Easing functions for smooth animation.

Usage:
    t = frame_index / total_frames  # 0.0 to 1.0
    value = ease_out(start, end, t)
"""

import math


def lerp(start, end, t):
    return start + (end - start) * t


def ease_linear(start, end, t):
    return lerp(start, end, t)


def ease_in(start, end, t):
    return lerp(start, end, t * t)


def ease_out(start, end, t):
    return lerp(start, end, 1 - (1 - t) * (1 - t))


def ease_in_out(start, end, t):
    if t < 0.5:
        return lerp(start, end, 2 * t * t)
    return lerp(start, end, 1 - (-2 * t + 2) ** 2 / 2)


def ease_in_cubic(start, end, t):
    return lerp(start, end, t * t * t)


def ease_out_cubic(start, end, t):
    return lerp(start, end, 1 - (1 - t) ** 3)


def ease_in_out_cubic(start, end, t):
    if t < 0.5:
        return lerp(start, end, 4 * t * t * t)
    return lerp(start, end, 1 - (-2 * t + 2) ** 3 / 2)


def bounce_out(start, end, t):
    n1 = 7.5625
    d1 = 2.75
    if t < 1 / d1:
        result = n1 * t * t
    elif t < 2 / d1:
        t -= 1.5 / d1
        result = n1 * t * t + 0.75
    elif t < 2.5 / d1:
        t -= 2.25 / d1
        result = n1 * t * t + 0.9375
    else:
        t -= 2.625 / d1
        result = n1 * t * t + 0.984375
    return lerp(start, end, result)


def elastic_out(start, end, t):
    if t == 0 or t == 1:
        return lerp(start, end, t)
    return lerp(start, end,
                math.pow(2, -10 * t) * math.sin((t * 10 - 0.75) * (2 * math.pi) / 3) + 1)


def back_out(start, end, t):
    c1 = 1.70158
    c3 = c1 + 1
    return lerp(start, end, 1 + c3 * math.pow(t - 1, 3) + c1 * math.pow(t - 1, 2))


FUNCTIONS = {
    "linear": ease_linear,
    "ease_in": ease_in,
    "ease_out": ease_out,
    "ease_in_out": ease_in_out,
    "ease_in_cubic": ease_in_cubic,
    "ease_out_cubic": ease_out_cubic,
    "ease_in_out_cubic": ease_in_out_cubic,
    "bounce_out": bounce_out,
    "elastic_out": elastic_out,
    "back_out": back_out,
}


def interpolate(start, end, t, easing="ease_out"):
    """Interpolate between start and end using named easing function."""
    func = FUNCTIONS.get(easing, ease_out)
    return func(start, end, t)
