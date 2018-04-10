package main

import (
    "image/color"
    "math"
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

func SphereVolume(radius float32) float32 {
    return float32((4. / 3.) * math.Pi * math.Pow(float64(radius), 3))
}

func BoxVolume(min, max Vector) float32 {
    x := max.X - min.X
    y := max.Y - min.Y
    z := max.Z - min.Z

    return x * z * y
}

func IsZero(f float32) bool {
    if f < 0 {
        return f > -EPSYLON
    }
    return f < EPSYLON
}

func MeshFindCenter(triangles []Triangle) Vector {
    var sum Vector
    divider := 1. / float32(len(triangles) * 3)

    for _, tri := range triangles {
        v := tri.A.Add(tri.B).Add(tri.C)
        sum = sum.Add(v.MulScal(divider))
    }

    return sum
}

func MeshFindBounds(triangles []Triangle) (box Box, sphere Sphere) {
    var min = Vector{0, 0, 0}
    var max = Vector{0, 0, 0}

    for _, tri := range triangles {
        min = MinVec(MinVec(MinVec(min, tri.A), tri.B), tri.C)
        max = MaxVec(MaxVec(MaxVec(max, tri.A), tri.B), tri.C)
    }

    box = Box{ min, max, BoxVolume(min, max) }

    sphere.Center = max.Sub(min).MulScal(.5).Add(min)
    sphere.Radius = Max(min.Sub(sphere.Center).Magnitude(), max.Sub(sphere.Center).Magnitude())
    sphere.Volume = SphereVolume(sphere.Radius)

    return
}

func TriangleFindCenter(t Triangle) Vector {
    return t.A.Add(t.B.Add(t.C)).MulScal(1./3.)
}
