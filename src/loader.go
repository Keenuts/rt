package main

import (
    "encoding/json"
    "fmt"
    "image"
    "image/png"
    "io/ioutil"
    "os"
    "path"
    "time"
);

type SceneFile struct {
    Name string
    OutputSize [2]int

    Camera Camera
    Meshs []string
    MaterialLibs []string

    SceneObjects []SceneObject
}

func CreateDirectory(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, 0755)
        if err != nil {
            return false
        }
    }

    return true
}

func LoadJSON(filename string, storage interface{}) {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    switch (storage).(type) {
        case *SceneFile:
            err = json.Unmarshal(content, storage.(*SceneFile))
        case *Config:
            err = json.Unmarshal(content, storage.(*Config))
        default:
            panic("Unknown type !")
    }
}

func SceneFileToScene(file SceneFile, models []Model, mtlLibs []MaterialLib) (out Scene) {
    out.Name = file.Name
    out.OutputSize = file.OutputSize
    out.Camera = file.Camera
    out.Camera.Forward = out.Camera.Forward.Normalize()
    out.Camera.Up = out.Camera.Up.Normalize()

    for _, desc := range file.SceneObjects {
        if desc.ObjectID >= len(models) {
            panic("Invalid object ID")
        }

        var mtl Material
        if len(mtlLibs) > desc.MaterialLibID {
            mtl = mtlLibs[desc.MaterialLibID][desc.MaterialName]
        } else {
            mtl.Diffuse = Vector{ 1, 1, 1 }
        }
        dst := CreateObject(desc, models[desc.ObjectID], mtl)
        out.Objects = append(out.Objects, dst)
    }

    return
}

func LoadScene(filename string) Scene {
    var file SceneFile
    LoadJSON(filename, &file)

    var models []Model
    for _,path := range file.Meshs {
        models = append(models, ModelFromOBJ(path))
    }

    var materialLibs []MaterialLib
    for _,path := range file.MaterialLibs {
        materialLibs = append(materialLibs, MaterialLibFromMTL(path))
    }

    return SceneFileToScene(file, models, materialLibs)
}

func LoadConfig(filename string) Config {
    var out Config
    LoadJSON(filename, &out)
    return out
}

func GetFileHandleToDisk(config Config, prefix string, ext string, alt string) *os.File {
    if !CreateDirectory(config.OutputDir) {
        panic("unable to create output directory")
    }

    filename := ""
    if config.ForceOutputName {
        filename = path.Join(config.OutputDir, alt)
    } else {
        filename = prefix + time.Now().Format("2006-01-02--15-04-05") + ext
        filename = path.Join(config.OutputDir, filename)
    }


    f, err := os.Create(filename);
    if err != nil {
        panic(err)
    }

    return f
}

func WriteReportToDisk(config Config, infos RenderInfo) {
    if !config.SaveReport {
        return
    }

    raw, err := json.MarshalIndent(infos, "", "    ")
    if err != nil {
        panic("unable to serialize rendering informations")
    }

    f := GetFileHandleToDisk(config, "report-", ".json", config.ReportName)
    defer f.Close()

    fmt.Printf("writting report...\r")
    f.Write(raw)
    fmt.Printf("writting report...done\n")
}

func WriteImageToDisk(config Config, buffer *image.RGBA) {
    if !config.SavePicture {
        return
    }

    f := GetFileHandleToDisk(config, "output-", ".png", config.PictureName)
    defer f.Close()

    fmt.Printf("writting picture...\r")
    png.Encode(f, buffer);
    fmt.Printf("writting picture...done\n")
}

