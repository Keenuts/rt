package main

import (
    "bufio"
    "github.com/udhos/gwob"
    "os"
);

// Bounding box not used for now
type Box struct {
    Min, Max Vector
}

type Triangle struct {
    A, B, C Vector
}

type Object struct {
    Name string
    Bounds Box
    Triangles []Triangle
}

func CreateObjectFromOBJ(filename string) (o Object) {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    options := gwob.ObjParserOptions{}

    obj, err := gwob.NewObjFromReader("FIXME", reader, &options)
    if err != nil {
        panic(err)
    }

    if len(obj.Indices) % 3 != 0 {
        panic("Invalid mesh. Indices count not a multiple of 3.")
    }

    for i := 0; i < len(obj.Indices); i += 3 {
        v := make([]Vector, 3)

        for j := 0 ; j < 3; j++ {
            v[j] = Vector{obj.Coord[obj.Indices[i + j] * 8 + 0],
                          obj.Coord[obj.Indices[i + j] * 8 + 1],
                          obj.Coord[obj.Indices[i + j] * 8 + 2]}
        }
        t := Triangle{ v[0], v[1], v[2] }
        o.Triangles = append(o.Triangles, t)
    }

    return
}
