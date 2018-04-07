package main

import (
    "image/color"
    "math"
);

type Intersection struct {
    Position, Normal Vector
}

func IntersectPlane(r Ray, a, normal Vector) (bool, Intersection) {

    r.Direction = r.Direction.Normalize()

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
