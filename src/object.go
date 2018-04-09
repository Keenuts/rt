package main

import (
    "bufio"
    "github.com/udhos/gwob"
    "os"
);

// Bounding box not used for now
type Box struct {
    Min, Max Vector
    Volume float32
}

type Sphere struct {
    Center Vector
    Radius, Volume float32
}

type Triangle struct {
    A, B, C Vector
}

type Object struct {
    Name string
    Center Vector

    BoundingSphere Sphere
    BoundingBox Box

    Triangles []Triangle
    Vertex []Vector
}

func ObjectCreateFromOBJ(filename string) (o Object) {
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

    for i := 0; i < len(obj.Coord); i += 8 {
        o.Vertex = append(o.Vertex, Vector{ obj.Coord[i + 0],
                                            obj.Coord[i + 1],
                                            obj.Coord[i + 2]})
    }

    for i := 0; i < len(obj.Indices); i += 3 {
        t := Triangle {
            o.Vertex[obj.Indices[i + 0]],
            o.Vertex[obj.Indices[i + 1]],
            o.Vertex[obj.Indices[i + 2]],
        }
        o.Triangles = append(o.Triangles, t)
    }

    o.Center = ObjectFindCenter(o)
    o = ObjectCreateBounds(o)

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

func ObjectCreateBounds(o Object) Object {
    var min = Vector{0, 0, 0}
    var max = Vector{0, 0, 0}

    for _, vtx := range o.Vertex {
        min = MinVec(min, vtx)
        max = MaxVec(max, vtx)
    }

    o.BoundingBox = Box{ min, max, BoxVolume(min, max) }

    o.BoundingSphere.Center = max.Sub(min).MulScal(.5).Add(min)
    o.BoundingSphere.Radius = Max(Max(max.X - min.X, max.Y - min.Y), max.Z - min.Z) * .5
    o.BoundingSphere.Volume = SphereVolume(o.BoundingSphere.Radius)

    return o
}
