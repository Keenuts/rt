package main

import (
    "image/color"
    "math"
);

func RaycheckSphere(r Ray, s Sphere) bool {
    e0 := s.Center.Sub(r.Origin)

    v := e0.Dot(r.Direction)
    d2 := e0.Dot(e0) - v * v
    rad2 := s.Radius * s.Radius

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

func RaycheckBox(r Ray, b Box) bool {
    var tmin, tmax, tymin, tymax, tzmin, tzmax float32

    bounds := [2]Vector{ b.Min, b.Max }

    divx := 1. / r.Direction.X
    if divx >= 0 {
        tmin = (bounds[0].X - r.Origin.X) * divx
        tmax = (bounds[1].X - r.Origin.X) * divx
    } else {
        tmin = (bounds[1].X - r.Origin.X) * divx
        tmax = (bounds[0].X - r.Origin.X) * divx
    }

    divy := 1. / r.Direction.Y
    if divy >= 0 {
        tymin = (bounds[0].Y - r.Origin.Y) * divy
        tymax = (bounds[1].Y - r.Origin.Y) * divy
    } else {
        tymin = (bounds[1].Y - r.Origin.Y) * divy
        tymax = (bounds[0].Y - r.Origin.Y) * divy
    }

    if tmin > tymax || tymin > tmax {
        return false
    }

    tmin = Max(tmin, tymin)
    tmax = Min(tmax, tymax)

    divz := 1. / r.Direction.Z
    if divz >= 0 {
        tzmin = (bounds[0].Z - r.Origin.Z) * divz
        tzmax = (bounds[1].Z - r.Origin.Z) * divz
    } else {
        tzmin = (bounds[1].Z - r.Origin.Z) * divz
        tzmax = (bounds[0].Z - r.Origin.Z) * divz
    }

    return tmin <= tzmax && tzmin <= tmax
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

    out := Intersection{ r.Origin.Add(r.Direction.MulScal(t)), normal, t }
    return true, out
}

func IntersectTri(r Ray, t Triangle) (bool, Intersection) {
    ab := t.Vertex[1].Sub(t.Vertex[0]).Normalize()
    ac := t.Vertex[2].Sub(t.Vertex[0]).Normalize()
    bc := t.Vertex[2].Sub(t.Vertex[1]).Normalize()
    ca := t.Vertex[0].Sub(t.Vertex[2]).Normalize()
    normal := ab.Cross(ac).Normalize()

    if normal.Dot(r.Direction) < 0 {
        return false, Intersection{}
    }

    hit, info := IntersectPlane(r, t.Vertex[0], normal)
    if !hit {
        return false, Intersection{}
    }

    tmp := ab.Cross(info.Position.Sub(t.Vertex[0]))
    if info.Normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    tmp = bc.Cross(info.Position.Sub(t.Vertex[1]))
    if normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    tmp = ca.Cross(info.Position.Sub(t.Vertex[2]))
    if info.Normal.Dot(tmp) < 0 {
        return false, Intersection{}
    }

    bCoord := GetBarycentric(info.Position, t)

    info.Normal = t.Normals[0].MulScal(bCoord.X)
    info.Normal = info.Normal.Add(t.Normals[1].MulScal(bCoord.Y))
    info.Normal = info.Normal.Add(t.Normals[2].MulScal(bCoord.Z))

    return true, info
}

func IntersectKDTree(ray Ray, tree *KDTree) (touch bool, out Intersection) {
    out.Distance = float32(math.Inf(1))

    if tree == nil {
        return false, out
    }

    if !RaycheckBox(ray, tree.BoundingBox) {
        return false, out
    }


    if tree.Left != nil || tree.Right != nil {
        touchR, outR := IntersectKDTree(ray, tree.Left)
        touchL, outL := IntersectKDTree(ray, tree.Right)

        if outL.Distance < outR.Distance {
            out = outL
        } else {
            out = outR
        }

        if touchL || touchR {
            return true, out
        }
    }

    for _, tri := range tree.Triangles {
        hit, info := IntersectTri(ray, tri)

        distance := info.Position.Sub(ray.Origin).Magnitude()

        if hit && distance < out.Distance {
            out = info
        }
    }

    return !math.IsInf(float64(out.Distance), 1), out
}

func Intersect(ray Ray, obj Object) (bool, Intersection) {
    var intersection Intersection

    if !RaycheckSphere(ray, obj.BoundingSphere) {
        return false, intersection
    }

    if !RaycheckBox(ray, obj.BoundingBox) {
        return false, intersection
    }

    return IntersectKDTree(ray, &obj.Tree)
}

func TraceRay(config Config, scene Scene, ray Ray) color.Color {

    depth := float32(math.Inf(1));
    var intersection Intersection

    ray.Direction = ray.Direction.Normalize()
    for _, obj := range scene.Objects {
        hit, info := Intersect(ray, obj)

        if !hit {
            continue
        }

        intersection = info
        depth = ray.Origin.Sub(info.Position).Magnitude()
    }

    if math.IsInf(float64(depth), 1) {
        return color.RGBA{0, 0, 0, 255}
    }

    out := Vector{1., 1., 1.}
    light := Vector{0., -1., 0.02}.Normalize().Neg()

    out = intersection.Normal.AddScal(1.).MulScal(.5)
    _ = light

    return VectorToRGBA(out)
}
