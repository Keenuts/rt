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

func SceneFileToScene(file SceneFile, models []Model) (out Scene) {
    out.Name = file.Name
    out.OutputSize = file.OutputSize
    out.Camera = file.Camera

    for _, desc := range file.SceneObjects {
        if desc.ObjectID >= len(models) {
            panic("Invalid object ID")
        }

        dst := CreateObject(desc, models[desc.ObjectID])
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

    return SceneFileToScene(file, models)
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

    fmt.Printf("writting '%s'\n", filename)

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

    f.Write(raw)
}

func WriteImageToDisk(config Config, buffer *image.RGBA) {
    if !config.SavePicture {
        return
    }

    f := GetFileHandleToDisk(config, "output-", ".png", config.PictureName)
    defer f.Close()

    png.Encode(f, buffer);
}

