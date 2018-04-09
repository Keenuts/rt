package main

import "fmt"

func CreateObject(desc SceneObject, model Model) (out Object) {
    fmt.Println("Loading ", model.Name)

    out.Name = model.Name

    out.Triangles = make([]Triangle, len(model.Triangles))
    copy(out.Triangles, model.Triangles)
    out.Vertex = make([]Vector, len(model.Vertex))
    copy(out.Vertex, model.Vertex)

    out.Center = ObjectFindCenter(out)
    out.BoundingBox, out.BoundingSphere = ObjectFindBounds(out)

    return
}

func ObjectFindCenter(o Object) Vector {
    var sum Vector
    var count float32

    for _, vtx := range o.Vertex {
        sum = sum.Add(vtx)
        count += 1.
    }

    return sum.MulScal(1. / count)
}

func ObjectFindBounds(o Object) (box Box, sphere Sphere) {
    var min = Vector{0, 0, 0}
    var max = Vector{0, 0, 0}

    for _, vtx := range o.Vertex {
        min = MinVec(min, vtx)
        max = MaxVec(max, vtx)
    }

    box = Box{ min, max, BoxVolume(min, max) }

    sphere.Center = max.Sub(min).MulScal(.5).Add(min)
    sphere.Radius = Max(Max(max.X - min.X, max.Y - min.Y), max.Z - min.Z) * .5
    sphere.Volume = SphereVolume(o.BoundingSphere.Radius)

    return
}
