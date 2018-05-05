package main

import (
    "bufio"
    "fmt"
    "github.com/Keenuts/gwob"
    "os"
);

func triangleGetVertex(obj *gwob.Obj, stride int) Vector {
    x, y, z := obj.VertexCoordinates(stride)
    return Vector{ x, y, z }
}

func triangleGetNormal(obj *gwob.Obj, stride int) Vector {
    x, y, z := obj.NormCoordinates(stride)
    return Vector{ x, y, z }
}

func triangleGetUV(obj *gwob.Obj, stride int) Vector {
    x, y := obj.TextCoordinates(stride)
    return Vector{ x, y, 0 }
}

func ModelFromOBJ(filename string) (model Model) {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    options := gwob.ObjParserOptions{}

    fmt.Printf("reading %s...", filename)
    obj, err := gwob.NewObjFromReader("", reader, &options)
    if err != nil {
        panic(err)
    }

    if len(obj.Indices) % 3 != 0 {
        panic("Invalid mesh. Indices count not a multiple of 3.")
    }
    if obj.NormCoordFound {
        fmt.Printf("found normals...")
    }
    if obj.TextCoordFound {
        fmt.Printf("found UVs...")
    }

    vcount := obj.NumberOfElements()
    for i := 0; i < vcount; i++ {
        x, y, z := obj.VertexCoordinates(i)
        model.Vertex = append(model.Vertex, Vector{ x, y, z })
    }

    for i := 0; i < len(obj.Indices); i += 3 {
        var vtx, nrm, uv [3]Vector

        for j := 0; j < 3; j++ {
            vtx[j] = triangleGetVertex(obj, obj.Indices[i + j])
        }

        if obj.NormCoordFound {
            for j := 0; j < 3; j++ {
                nrm[j] = triangleGetNormal(obj, obj.Indices[i + j])
            }
        } else {
            normal := vtx[2].Sub(vtx[0]).Cross(vtx[1].Sub(vtx[0])).Normalize()
            nrm = [3]Vector{ normal, normal, normal }
        }

        if obj.TextCoordFound {
            for j := 0; j < 3; j++ {
                uv[j] = triangleGetUV(obj, obj.Indices[i + j])
            }
        }

        t := Triangle { vtx, nrm, uv }
        model.Triangles = append(model.Triangles, t)
    }

    model.Name = filename

    fmt.Printf("done\n")
    return
}
