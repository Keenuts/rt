package main

import (
    "bufio"
    "github.com/udhos/gwob"
    "os"
);

func ModelFromOBJ(filename string) (model Model) {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    options := gwob.ObjParserOptions{}

    obj, err := gwob.NewObjFromReader("", reader, &options)
    if err != nil {
        panic(err)
    }

    if len(obj.Indices) % 3 != 0 {
        panic("Invalid mesh. Indices count not a multiple of 3.")
    }

    vcount := obj.NumberOfElements()
    for i := 0; i < vcount; i++ {
        x, y, z := obj.VertexCoordinates(i)
        model.Vertex = append(model.Vertex, Vector{ x, y, z })
    }

    for i := 0; i < len(obj.Indices); i += 3 {
        vtx := [3]Vector {
            model.Vertex[obj.Indices[i + 0]],
            model.Vertex[obj.Indices[i + 1]],
            model.Vertex[obj.Indices[i + 2]],
        }

        t := Triangle { vtx }
        model.Triangles = append(model.Triangles, t)
    }

    model.Name = filename
    return
}
