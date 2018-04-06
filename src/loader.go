package main

import (
    "encoding/json";
    "image";
    "image/png";
    "io/ioutil";
    "os";
);

type SceneObject struct {
    ObjectID int
    Position Vector
    Rotation Vector
    Scale Vector
}

type Camera struct {
    Position, Forward, Up Vector
    Fov float32
}

type Scene struct {
    Name string
    OutputSize [2]int

    Camera Camera
    Models []string
    Objects []Object
    Scene []SceneObject

}

type Config struct {
    MaxThreads int
    BlockSize int
    OutputPath string
}

func LoadJSON(filename string, storage interface{}) {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    switch (storage).(type) {
        case *Scene:
            err = json.Unmarshal(content, storage.(*Scene))
        case *Config:
            err = json.Unmarshal(content, storage.(*Config))
        default:
            panic("Unknown type !")
    }
}

func LoadScene(filename string) Scene {
    var out Scene
    LoadJSON(filename, &out)

    for _,m := range out.Models {
        out.Objects = append(out.Objects, CreateObjectFromOBJ(m))
    }

    return out
}

func LoadConfig(filename string) Config {
    var out Config
    LoadJSON(filename, &out)
    return out
}

func WriteImageToDisk(path string, buffer *image.RGBA) {
    f, err := os.Create("/tmp/output.png");
    if err != nil {
        panic(err)
    }
    defer f.Close()

    png.Encode(f, buffer);
}
