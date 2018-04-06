package main

import "fmt"

func main() {
    scene := LoadScene("default.json")
    config := LoadConfig("config.json")

    fmt.Println(scene)

    output := RenderScene(config, scene)

    WriteImageToDisk(config.OutputPath, output)
}
