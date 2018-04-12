package main

import (
    "bufio"
    "fmt"
    "github.com/udhos/gwob"
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
        m.Opacity = value.D
        lib[key] = m
    }

    fmt.Printf("done\n")
    fmt.Println(lib)
    return
}
