package main

import (
    "image/color"
);

func Lerp(a, b, x float32) float32 {
    return a * (1. - x) + b * x
}

func Max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}

func Min(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}

func Clamp(a, b, x float32) float32 {
    return Max(a, Min(b, x))
}

func VectorToRGBA(v Vector) color.RGBA {
    v = v.Normalize()
    return color.RGBA{
        uint8(v.X * 255.),
        uint8(v.Y * 255.),
        uint8(v.Z * 255.),
        255,
    }
}
