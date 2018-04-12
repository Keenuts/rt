package main

import (
    "math"
)

// Warning: Matrix are written like this: [ [col1] [col2] [col3] ]
// Don't be fooled by the layout
func MatrixCreateRotateX(angle float64) (m Mat3) {
    return Mat3{
        { 1, 0,                              0 },
        { 0, math.Cos(angle),  math.Sin(angle) },
        { 0, -math.Sin(angle), math.Cos(angle) },
    }
}

func MatrixCreateRotateY(angle float64) (m Mat3) {
    return Mat3{
        { math.Cos(angle),    0,  -math.Sin(angle) },
        { 0,                            1,  0 },
        { math.Sin(angle),    0,  math.Cos(angle) },
    }
}

func MatrixCreateRotateZ(angle float64) (m Mat3) {
    return Mat3{
        { math.Cos(angle),    math.Sin(angle),  0 },
        { -math.Sin(angle),   math.Cos(angle),  0 },
        { 0,                            0,                          1 },
    }
}

func (m Mat3) Mul(v Vector) (out Vector) {
    out.X = m[0][0] * v.X + m[1][0] * v.Y + m[2][0] * v.Z
    out.Y = m[0][1] * v.X + m[1][1] * v.Y + m[2][1] * v.Z
    out.Z = m[0][2] * v.X + m[1][2] * v.Y + m[2][2] * v.Z
    return
}
