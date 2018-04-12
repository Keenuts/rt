package main

import "fmt"

func CreateObject(desc SceneObject, model Model, mtl Material) (out Object) {
    fmt.Printf("preprocessing %s...", model.Name)

    out.Name = model.Name

    out.Triangles = make([]Triangle, len(model.Triangles))
    copy(out.Triangles, model.Triangles)

    ObjectTransform(out, desc)

    out.Center = MeshFindCenter(out.Triangles)
    out.BoundingBox, out.BoundingSphere = MeshFindBounds(out.Triangles)
    out.Tree = TreeCreate(out.Triangles)
    out.Material = mtl

    fmt.Printf("done\n")
    return
}

func ObjectTransform(obj Object, desc SceneObject) {

    for i, tri := range obj.Triangles {

        for j := 0; j < 3; j++ {
            v := tri.Vertex[j].Scale(desc.Scale)
            v = v.RotateDeg(desc.Rotation)
            v = v.Add(desc.Position)

            tri.Vertex[j] = v
        }

        for j := 0; j < 3; j++ {
            v := tri.Normals[j].RotateDeg(desc.Rotation)

            tri.Normals[j] = v.Normalize()
        }

        obj.Triangles[i] = tri
    }
}
