package main

import (
    "encoding/json"
    "io/ioutil"
    "time"
)

type RenderInfo struct {
    SceneName string
    Date string
    OutputFile string
    OutputSize [2]int
    Duration string
    Threads int
};

func WriteRecapToDisk(infos RenderInfo) {
    raw, err := json.MarshalIndent(infos, "", "    ")
    if err != nil {
        panic("unable to serialize rendering informations")
    }

    filename := "rendering-" + infos.Date + ".json"
    err = ioutil.WriteFile(filename, raw, 0644)
    if err != nil {
        panic(err)
    }
}

func main() {
    var infos RenderInfo

    scene := LoadScene("default.json")
    config := LoadConfig("config.json")

    infos.OutputFile = config.OutputPath
    infos.Date = time.Now().Format("2006-01-02--15-04-05")
    infos.Threads = config.MaxThreads
    infos.OutputSize = scene.OutputSize
    infos.SceneName = scene.Name
    start := time.Now()

    output := RenderScene(config, scene)

    infos.Duration = time.Now().Sub(start).String()

    WriteImageToDisk(config.OutputPath, output)
    WriteRecapToDisk(infos)
}
