package main

import "math"

func VecFromInt(x, y, z int) Vector {
    return Vector{float32(x), float32(y), float32(z)}
}

func (a Vector) Add(b Vector) Vector {
    return Vector{ a.X + b.X, a.Y + b.Y, a.Z + b.Z }
}

func (a Vector) Sub(b Vector) Vector {
    return Vector{ a.X - b.X, a.Y - b.Y, a.Z - b.Z }
}

func (a Vector) MulScal(m float32) Vector {
    return Vector{ a.X * m, a.Y * m, a.Z * m }
}

func (a Vector) AddScal(s float32) Vector {
    return Vector{ a.X + s, a.Y + s, a.Z + s }
}

func (a Vector) Neg() Vector {
    return Vector{ -a.X, -a.Y, -a.Z }
}

func (v Vector) Magnitude() float32 {
    return float32(math.Sqrt(float64(v.X * v.X + v.Y * v.Y + v.Z * v.Z)))
}

func (v Vector) Normalize() Vector {
    norm := v.Magnitude()
    if norm < EPSYLON && norm > -EPSYLON {
        return v
    }

    return Vector{ v.X / norm, v.Y / norm, v.Z / norm}
}

func (a Vector) Cross(b Vector) Vector {
    return Vector{
        b.Y * a.Z - b.Z * a.Y,
        b.Z * a.X - b.X * a.Z,
        b.X * a.Y - b.Y * a.X,
    }
}

func (a Vector) Dot(b Vector) float32 {
    return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}

func MinVec(a, b Vector) Vector {
    return Vector{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}

func MaxVec(a, b Vector) Vector {
    return Vector{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}
