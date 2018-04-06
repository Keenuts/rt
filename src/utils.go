package main

func lerp(a, b, x float32) float32 {
    return a * (1. - x) + b * x
}
