package main

import (
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

    d := math.Sqrt(rad2 - d2)
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
    var tmin, tmax, tymin, tymax, tzmin, tzmax float64

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

    if math.Abs(d) < EPSYLON {
        return false, Intersection{}
    }

    t := a.Sub(r.Origin).Dot(normal) / d
    if t < 0. {
        return false, Intersection{}
    }

    uv := Vector{0, 0, 0}
    var out Intersection
    out.Position = r.Origin.Add(r.Direction.MulScal(t))
    out.Normal = normal
    out.UV = uv
    out.Distance = t
    return true, out
}

func IntersectTri(r Ray, t Triangle) (bool, Intersection) {
    ab := t.Vertex[1].Sub(t.Vertex[0]).Normalize()
    ac := t.Vertex[2].Sub(t.Vertex[0]).Normalize()
    bc := t.Vertex[2].Sub(t.Vertex[1]).Normalize()
    ca := t.Vertex[0].Sub(t.Vertex[2]).Normalize()
    normal := ab.Cross(ac).Normalize()

    if normal.Dot(r.Direction) < 0 && !r.InvertCulling {
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

    info.Normal = Vector{ 0, 0, 0 }
    info.Normal = info.Normal.Add(t.Normals[1].MulScal(bCoord.X))
    info.Normal = info.Normal.Add(t.Normals[2].MulScal(bCoord.Y))
    info.Normal = info.Normal.Add(t.Normals[0].MulScal(bCoord.Z))

    info.UV = Vector{ 0, 0, 0 }
    info.UV = info.UV.Add(t.UV[1].MulScal(bCoord.X))
    info.UV = info.UV.Add(t.UV[2].MulScal(bCoord.Y))
    info.UV = info.UV.Add(t.UV[0].MulScal(bCoord.Z))

    return true, info
}

func IntersectKDTree(ray Ray, tree *KDTree) (touch bool, out Intersection) {
    out.Distance = math.Inf(1)

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

        info.Distance = info.Position.Sub(ray.Origin).Magnitude()

        if hit && info.Distance < out.Distance {
            out = info
        }
    }

    return !math.IsInf(out.Distance, 1), out
}

func Intersect(ray Ray, obj Object) (bool, Intersection) {
    var intersection Intersection
    var hit bool

    if !RaycheckSphere(ray, obj.BoundingSphere) {
        return false, intersection
    }

    if !RaycheckBox(ray, obj.BoundingBox) {
        return false, intersection
    }

    hit, intersection = IntersectKDTree(ray, &obj.Tree)
    if hit {
        intersection.Object = obj
    }

    return hit, intersection
}

func IntersectObjects(ray Ray, objects []Object) (bool, Intersection) {
    var hit bool
    var intersection Intersection
    intersection.Distance = math.Inf(1)

    ray.Direction = ray.Direction.Normalize()
    for _, obj := range objects {
        touch, info := Intersect(ray, obj)

        if !touch || info.Distance > intersection.Distance {
            continue
        }

        hit = true
        intersection = info
    }

    return hit, intersection
}

func TraceRefraction(scene Scene, ray Ray, info Intersection, lDepth float64) Vector {
    ratioOut := info.Object.Material.Refraction
    ratioIn := 1.0 / ratioOut

    posA := info.Position
    norA := info.Normal

    var ray2 Ray
    ray2.Origin = posA.Add(norA.MulScal(-10. * EPSYLON))
    ray2.Direction = Refract(ray.Direction, info.Normal, ratioIn)
    ray2.InvertCulling = true

    hit2, info2 := Intersect(ray2, info.Object)
    ray2.InvertCulling = false

    if hit2 && info.Object.ID == info2.Object.ID {
        posB := info2.Position
        norB := info2.Normal

        ray2.Origin = posB.Add(norB.MulScal(10. * EPSYLON))
        ray2.Direction = Refract(ray2.Direction, info2.Normal.Neg(), ratioOut)
    }

    refracted, _ := TraceRayDepth(scene, ray2, lDepth - 1)
    return refracted
}

func GetDiffuse(info Intersection) Vector {
    mtl := info.Object.Material

    if mtl.DiffuseTex == nil {
        return mtl.Diffuse
    }

    x := int(float64(mtl.DiffuseTex.Width) * info.UV.X)
    y := int(float64(mtl.DiffuseTex.Height) * info.UV.Y)

    color := mtl.DiffuseTex.Pixels.At(x, y)

    return ColorToVector(color)
}

func TraceRayDepth(scene Scene, ray Ray, leftDepth float64) (Vector, float64) {
    if leftDepth <= 0 {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    hit, info := IntersectObjects(ray, scene.Objects)
    if !hit {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    diffuse := GetDiffuse(info)
    specularLevel := info.Object.Material.SpecularLevel

    var reflected Vector
    if specularLevel > 0. {
        var reflRay Ray
        reflRay.Origin = info.Position.Add(info.Normal.MulScal(EPSYLON))
        reflRay.Direction = Reflect(ray.Direction, info.Normal)
        reflected, _ = TraceRayDepth(scene, reflRay, leftDepth - 1)
    }
    reflected = reflected.MulScal(specularLevel * 0.001)

    var fresnel float64
    var refracted Vector
    opacity := info.Object.Material.Opacity
    if opacity < 1. {
        refracted = TraceRefraction(scene, ray, info, leftDepth)
        fresnel = Fresnel(ray.Direction, info.Normal, 1.0 / 1.5161)
    } else {
        refracted = diffuse
        fresnel = 1.
    }
    refracted = refracted.MulScal(1. - fresnel)

    output := diffuse.MulScal(1. - specularLevel * 0.001)
    output = output.Add(reflected).MulScal(fresnel)
    output = output.Add(refracted)

    return Saturate(output), info.Distance
}

func TraceRay(scene Scene, ray Ray) (Vector, float64) {
    return TraceRayDepth(scene, ray, 8)
}
