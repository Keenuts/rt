package main

import "fmt"

func CreateObject(desc SceneObject, model Model) (out Object) {
    fmt.Println("Loading ", model.Name)

    out.Name = model.Name

    out.Triangles = make([]Triangle, len(model.Triangles))
    copy(out.Triangles, model.Triangles)

    ObjectTransform(out, desc)

    out.Center = MeshFindCenter(out.Triangles)
    out.BoundingBox, out.BoundingSphere = MeshFindBounds(out.Triangles)
    out.Tree = TreeCreate(out.Triangles)

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
