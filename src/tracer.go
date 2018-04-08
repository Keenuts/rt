package main

import (
    "image/color"
    "math"
);

type Intersection struct {
    Position, Normal Vector
}

func IntersectPlane(r Ray, a, normal Vector) (bool, Intersection) {
    d := normal.Dot(r.Direction)

    if math.Abs(float64(d)) < EPSYLON {
        return false, Intersection{}
    }

    t := a.Sub(r.Origin).Dot(normal) / d
    if t < 0. {
        return false, Intersection{}
    }

    return true, Intersection{
            r.Origin.Add(r.Direction.MulScal(t)),
            normal,
        }

}

func IntersectSphere(r Ray, o Vector, rad float32) bool {
    e0 := o.Sub(r.Origin)

    v := e0.Dot(r.Direction)
    d2 := e0.Dot(e0) - v * v
    rad2 := rad * rad

    if d2 > rad2 {
        return false
    }

    d := float32(math.Sqrt(float64(rad2 - d2)))
    t0 := v - d
    t1 := v + d

    if t0 > t1 {
        t1, t0 = t0, t1
    }

    if t0 < 0 {
        return t1 > 0
    }

    return t0 > 0
}

func IntersectTri(r Ray, t Triangle) (bool, Intersection) {
    ab := t.B.Sub(t.A).Normalize()
    ac := t.C.Sub(t.A).Normalize()
    normal := ab.Cross(ac).Normalize()

    if normal.Dot(r.Direction) < 0 {
        return false, Intersection{}
    }

    hit, info := IntersectPlane(r, t.A, normal)
    if !hit {
        return false, Intersection{}
    }

    tmp := ab.Cross(info.Position.Sub(t.A))
    if info.Normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    tmp = t.C.Sub(t.B).Cross(info.Position.Sub(t.B))
    if normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    tmp = t.A.Sub(t.C).Cross(info.Position.Sub(t.C))
    if info.Normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    return true, info
}

func Intersect(obj Object, ray Ray) (bool, Intersection) {
    depth := float32(math.Inf(1))
    var intersection Intersection

    if !IntersectSphere(ray, obj.Center, obj.BoundsRadius) {
        return false, intersection
    }

    for _, tri := range obj.Triangles {
        hit, info := IntersectTri(ray, tri)

        dist := info.Position.Sub(ray.Origin).Magnitude()

        if hit && dist < depth {
            intersection = info
            depth = dist
        }
    }

    return !math.IsInf(float64(depth), 1), intersection
}

func TraceRay(config Config, scene Scene, ray Ray) color.Color {

    depth := float32(math.Inf(1));
    var intersection Intersection

    ray.Direction = ray.Direction.Normalize()
    for _, obj := range scene.Objects {
        hit, info := Intersect(obj, ray)

        if !hit {
            continue
        }

        intersection = info
        depth = ray.Origin.Sub(info.Position).Magnitude()
    }

    if math.IsInf(float64(depth), 1) {
        return color.RGBA{255, 0, 0, 255}
    }

    out := Vector{1., 1., 1.}
    light := Vector{0., -1., 0.02}.Normalize()

    return VectorToRGBA(out.MulScal(Clamp(0, 1, intersection.Normal.Dot(light))))
}
