package main

import "fmt"

func CreateObject(desc SceneObject, model Model) (out Object) {
    fmt.Println("Loading ", model.Name)

    out.Name = model.Name

    out.Triangles = make([]Triangle, len(model.Triangles))
    copy(out.Triangles, model.Triangles)

    ObjectTransform(out, desc)

    out.Center = ObjectFindCenter(out)
    out.BoundingBox, out.BoundingSphere = ObjectFindBounds(out)

    return
}

func ObjectTransform(obj Object, desc SceneObject) {

    for i, tri := range obj.Triangles {
        tri.A = tri.A.Scale(desc.Scale).RotateDeg(desc.Rotation).Add(desc.Position)
        tri.B = tri.B.Scale(desc.Scale).RotateDeg(desc.Rotation).Add(desc.Position)
        tri.C = tri.C.Scale(desc.Scale).RotateDeg(desc.Rotation).Add(desc.Position)

        obj.Triangles[i] = tri
    }
}

func ObjectFindCenter(o Object) Vector {
    var sum Vector
    var count float32

    for _, tri := range o.Triangles {
        sum = sum.Add(tri.A).Add(tri.B).Add(tri.C)
        count += 3.
    }

    return sum.MulScal(1. / count)
}

func ObjectFindBounds(o Object) (box Box, sphere Sphere) {
    var min = Vector{0, 0, 0}
    var max = Vector{0, 0, 0}

    for _, tri := range o.Triangles {
        min = MinVec(MinVec(MinVec(min, tri.A), tri.B), tri.C)
        max = MaxVec(MaxVec(MaxVec(max, tri.A), tri.B), tri.C)
    }

    box = Box{ min, max, BoxVolume(min, max) }

    sphere.Center = max.Sub(min).MulScal(.5).Add(min)
    sphere.Radius = Max(Max(max.X - min.X, max.Y - min.Y), max.Z - min.Z) * .5
    sphere.Volume = SphereVolume(o.BoundingSphere.Radius)

    return
}
