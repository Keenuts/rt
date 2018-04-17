package main

import (
    "bufio"
    "fmt"
    "github.com/Keenuts/gwob"
    "os"
);

func MaterialLibFromMTL(filename string) (lib MaterialLib) {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    options := gwob.ObjParserOptions{}

    fmt.Printf("reading %s...", filename)
    mtl, err := gwob.ReadMaterialLibFromReader(reader, &options)
    if err != nil {
        panic(err)
    }

    lib = make(MaterialLib)
    for key, value := range mtl.Lib {
        var m Material

        m.Diffuse = Vector{ value.Kd[0], value.Kd[1], value.Kd[2] }
        m.Specular = Vector{ value.Ks[0], value.Ks[1], value.Ks[2] }
        m.Emission = Vector{ value.Ke[0], value.Ke[1], value.Ke[2] }
        m.Opacity = value.D
        m.Refraction = value.Refr
        m.SpecularLevel = value.Ns
        lib[key] = m
    }

    fmt.Printf("done\n")
    fmt.Println(lib)
    return
}
