package main

import (
    "image/color"
    "math"
    "math/rand"
);

func Lerp(a, b, x float64) float64 {
    return a * (1. - x) + b * x
}

func Max(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

func Min(a, b float64) float64 {
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

func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func MinInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func Clamp(a, b, x float64) float64 {
    return Max(a, Min(b, x))
}

func Saturate(v Vector) Vector {
    return Vector{
        Clamp(0., 1., v.X),
        Clamp(0., 1., v.Y),
        Clamp(0., 1., v.Z),
    }
}

func ColorToVector(c color.Color) Vector {
    values := c.(color.RGBA)

    return Vector{
        float64(values.R) / 255.,
        float64(values.G) / 255.,
        float64(values.B) / 255.,
    }
}

func VectorToRGBA(v Vector) color.RGBA {
    return color.RGBA{
        uint8(v.X * 255.),
        uint8(v.Y * 255.),
        uint8(v.Z * 255.),
        255,
    }
}

func TriangleSurface(a, b, c Vector) float64 {
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

func SphereVolume(radius float64) float64 {
    return (4. / 3.) * math.Pi * math.Pow(radius, 3)
}

func BoxVolume(min, max Vector) float64 {
    x := max.X - min.X
    y := max.Y - min.Y
    z := max.Z - min.Z

    return x * z * y
}

func IsZero(f float64) bool {
    if f < 0 {
        return f > -EPSYLON
    }
    return f < EPSYLON
}

func MeshFindCenter(triangles []Triangle) Vector {
    var sum Vector
    divider := 1. / float64(len(triangles) * 3)

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

func CheckerGetColor(v Vector) Vector {
    v.X *= 5
    v.Y *= 5
    if (int(v.X) + int(v.Y)) % 2 == 0 {
        return Vector{0, 0, 0 }
    }
    return Vector{1, 1, 1}
}

func Fresnel(i, n Vector, eta float64) float64 {
    facing := Clamp(0., 1., 1.0 - Max(i.Neg().Dot(n), 0.))
    return Clamp(0., 1., eta + (1. - eta) * math.Pow(facing, 5))
}

func RandomUnitVector() Vector {
    var v Vector

    for v.X * v.X + v.Y * v.Y + v.Z * v.Z < EPSYLON {
        v = Vector{ rand.Float64(), rand.Float64(), rand.Float64() }
        v = v.AddScal(-.5).MulScal(2.)
    }

    return v.Normalize()
}

func RandomHemisphereVector(normal Vector) Vector {
    v := RandomUnitVector()

    if v.Dot(normal) < 0 {
        return v.Neg()
    }
    return v
}
