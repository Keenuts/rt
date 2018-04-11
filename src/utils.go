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

func AbsInt(a int) int {
    if a >= 0 {
        return a
    }
    return -a
}

func Clamp(a, b, x float32) float32 {
    return Max(a, Min(b, x))
}

func Saturate(v Vector) Vector {
    return Vector{
        Clamp(0., 1., v.X),
        Clamp(0., 1., v.Y),
        Clamp(0., 1., v.Z),
    }
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

func TriangleSurface(a, b, c Vector) float32 {
    return b.Sub(a).Cross(a.Sub(c)).Magnitude()
}

func GetBarycentric(p Vector, tri Triangle) Vector {
    IAabc := 1. / TriangleSurface(tri.Vertex[0], tri.Vertex[1], tri.Vertex[2])
    Aapc := TriangleSurface(tri.Vertex[0], p, tri.Vertex[2])
    Aapb := TriangleSurface(tri.Vertex[0], p, tri.Vertex[1])
    Abpc := TriangleSurface(tri.Vertex[1], p, tri.Vertex[2])

    u := Aapc * IAabc
    v := Aapb * IAabc
    w := Abpc * IAabc

    return Vector{ u, v, w }
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
        v := tri.Vertex[0].Add(tri.Vertex[1]).Add(tri.Vertex[2])
        sum = sum.Add(v.MulScal(divider))
    }

    return sum
}

func MeshFindBounds(triangles []Triangle) (box Box, sphere Sphere) {
    min := triangles[0].Vertex[0]
    max := triangles[0].Vertex[0]

    for _, tri := range triangles {
        min = MinVec(MinVec(MinVec(min, tri.Vertex[0]), tri.Vertex[1]), tri.Vertex[2])
        max = MaxVec(MaxVec(MaxVec(max, tri.Vertex[0]), tri.Vertex[1]), tri.Vertex[2])
    }

    box = Box{ min, max, BoxVolume(min, max) }

    sphere.Center = max.Sub(min).MulScal(.5).Add(min)
    sphere.Radius = Max(min.Sub(sphere.Center).Magnitude(), max.Sub(sphere.Center).Magnitude())
    sphere.Volume = SphereVolume(sphere.Radius)

    return
}

func TriangleFindCenter(t Triangle) Vector {
    return t.Vertex[0].Add(t.Vertex[1].Add(t.Vertex[2])).MulScal(1./3.)
}
