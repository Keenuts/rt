package main

import (
    "time"
)

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
