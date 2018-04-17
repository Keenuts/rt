package main

import "math"

func VecFromInt(x, y, z int) Vector {
    return Vector{float64(x), float64(y), float64(z)}
}

func (a Vector) Add(b Vector) Vector {
    return Vector{ a.X + b.X, a.Y + b.Y, a.Z + b.Z }
}

func (a Vector) Sub(b Vector) Vector {
    return Vector{ a.X - b.X, a.Y - b.Y, a.Z - b.Z }
}

func (a Vector) Scale(b Vector) Vector {
    return Vector{ a.X * b.X, a.Y * b.Y, a.Z * b.Z }
}

func (a Vector) RotateDeg(b Vector) Vector {
    return a.RotateRad(b.MulScal(DEG2RAD))
}

func (a Vector) RotateRad(b Vector) Vector {
    mX := MatrixCreateRotateX(b.X)
    mY := MatrixCreateRotateY(b.Y)
    mZ := MatrixCreateRotateZ(b.Z)

    return mZ.Mul(mY.Mul(mX.Mul(a)))
}

func (a Vector) MulScal(m float64) Vector {
    return Vector{ a.X * m, a.Y * m, a.Z * m }
}

func (a Vector) AddScal(s float64) Vector {
    return Vector{ a.X + s, a.Y + s, a.Z + s }
}

func (a Vector) Neg() Vector {
    return Vector{ -a.X, -a.Y, -a.Z }
}

func (v Vector) Magnitude() float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y + v.Z * v.Z)
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

func (a Vector) Dot(b Vector) float64 {
    return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}

func MinVec(a, b Vector) Vector {
    return Vector{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}

func MaxVec(a, b Vector) Vector {
    return Vector{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}

func Reflect(in, n Vector) (out Vector) {
    in = in.Normalize()
    n = n.Normalize()

    return in.Sub(n.MulScal(2. * n.Dot(in)))
}

func Refract(i, n Vector, eta float64) (out Vector) {
    n = n.Normalize()
    i = i.Normalize()

    cosi := i.Neg().Dot(n)
    cost2 := 1. - eta * eta * (1. - cosi * cosi)

    x := eta * cosi - math.Sqrt(math.Abs(cost2))
    t := i.MulScal(eta).Add(n.MulScal(x))

    if cost2 > 0 {
        return t
    }
    return i
}
