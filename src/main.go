package main

import (
    "time"
)

type RenderInfo struct {
    SceneName string
    Date string
    OutputSize [2]int
    Duration string
    Threads int
    Config Config
};

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
    OutputDir string
    SavePicture bool
    SaveReport bool
    ForceOutputName bool
    PictureName string
    ReportName string
}

func main() {
    var infos RenderInfo

    scene := LoadScene("default.json")
    config := LoadConfig("config.json")

    infos.Config = config
    infos.Date = time.Now().Format("2006-01-02--15-04-05")
    infos.OutputSize = scene.OutputSize
    infos.SceneName = scene.Name
    start := time.Now()

    output := RenderScene(config, scene)

    infos.Duration = time.Now().Sub(start).String()

    WriteImageToDisk(config, output)
    WriteReportToDisk(config, infos)
}
