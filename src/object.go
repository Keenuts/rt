package main

import "fmt"

func CreateObject(desc SceneObject, model Model) (out Object) {
    fmt.Printf("preprocessing %s...", model.Name)

    out.Name = model.Name

    out.Triangles = make([]Triangle, len(model.Triangles))
    copy(out.Triangles, model.Triangles)

    ObjectTransform(out, desc)

    out.Center = MeshFindCenter(out.Triangles)
    out.BoundingBox, out.BoundingSphere = MeshFindBounds(out.Triangles)
    out.Tree = TreeCreate(out.Triangles)

    fmt.Printf("done\n")
    return
}

func ObjectTransform(obj Object, desc SceneObject) {

    for i, tri := range obj.Triangles {

        for j := 0; j < 3; j++ {
            tri.Vertex[j] = tri.Vertex[j].Scale(desc.Scale).RotateDeg(desc.Rotation).Add(desc.Position)
        }

        obj.Triangles[i] = tri
    }
}
