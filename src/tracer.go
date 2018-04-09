package main

import (
    "image/color"
    "math"
);

func IntersectSphere(r Ray, s Sphere) bool {
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

// Returns true if inside, and closest coord to space between min and max
func GetClosestBound(p, min, max float32) (float32, bool) {
    if p < min {
        return min, false
    }
    if p > max {
        return max, false
    }
    return p, true
}

func IntersectBox(r Ray, b Box) bool {
    ori := []float32 { r.Origin.X, r.Origin.Y, r.Origin.Z }
    dir := []float32 { r.Direction.X, r.Direction.Y, r.Direction.Z }
    min := []float32 { b.Min.X, b.Min.Y, b.Min.Z }
    max := []float32 { b.Max.X, b.Max.Y, b.Max.Z }

    var steps [3][3]float32
    var inside [3]bool

    for i := 0; i < 3; i++ {
        steps[0][i], inside[i] = GetClosestBound(ori[i], min[i], max[i])
    }

    if inside[0] && inside[1] && inside[2] {
        return true
    }

    for i := 0; i < 3; i++ {
        if !inside[i] && !IsZero(dir[i]) {
            steps[1][i] = (steps[0][i] - ori[i]) / dir[i]
        } else {
            steps[1][i] = -1
        }
    }

    t_lim := 0
    for i := 1; i < 3; i++ {
        if steps[1][t_lim] < steps[1][i] {
            t_lim = i
        }
    }

    if steps[1][t_lim] < 0 {
        return false
    }

    for i := 0; i < 3; i++ {
        if i == t_lim {
            steps[2][i] = steps[1][i]
        } else {
            steps[2][i] = ori[i] + dir[i] * steps[1][t_lim]
            if steps[2][i] < min[i] || steps[2][i] > max[i] {
                return false
            }
        }
    }

    return true
}

func IntersectPlane(r Ray, a, normal Vector) (bool, Intersection) {
    d := normal.Dot(r.Direction)

    // Normal and ray are perpendicular
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

func IntersectKDTree(ray Ray, tree *KDTree) (bool, Intersection) {
    var out Intersection
    var touch bool

    if !IntersectBox(ray, tree.BoundingBox) {
        return false, out
    }

    if tree.Left != nil || tree.Right != nil {
        if tree.Left != nil {
            touch, out = IntersectKDTree(ray, tree.Left)
            if touch {
                return touch, out
            }
        }

        if tree.Right != nil {
            return IntersectKDTree(ray, tree.Right)
        }
    }

    depth := float32(math.Inf(1))

    for _, tri := range tree.Triangles {
        hit, info := IntersectTri(ray, tri)

        dist := info.Position.Sub(ray.Origin).Magnitude()

        if hit && dist < depth {
            out = info
            depth = dist
        }
    }

    return !math.IsInf(float64(depth), 1), out
}

func Intersect(ray Ray, obj Object) (bool, Intersection) {
    var intersection Intersection

    if !IntersectSphere(ray, obj.BoundingSphere) {
        return false, intersection
    }

    if !IntersectBox(ray, obj.BoundingBox) {
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
        return color.RGBA{255, 0, 0, 255}
    }

    out := Vector{1., 1., 1.}
    light := Vector{0., -1., 0.02}.Normalize()

    return VectorToRGBA(out.MulScal(Clamp(0, 1, intersection.Normal.Dot(light))))
}
