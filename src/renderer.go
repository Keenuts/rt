package main

import (
    "image";
    "fmt";
    "sync";
    "image/draw";
    "image/color";
);

type Task struct {
    Area image.Rectangle
    Buffer *image.RGBA
};

func RenderArea(config Config, scene Scene, task *Task) {
    rect := image.Rect(0, 0, config.BlockSize, config.BlockSize)
    task.Buffer = image.NewRGBA(rect)


    r := (float32(task.Area.Min.X) / 512.) * 255.
    g := (float32(task.Area.Min.Y) / 512.) * 255.

    c := color.RGBA{uint8(r), uint8(g), 128, 255}

    draw.Draw(task.Buffer, rect, &image.Uniform{c}, image.ZP, draw.Src)
}

func CreateRenderTasks(config Config, scene Scene) (taskList []Task) {
    taskList = make([]Task, 0)

    for y := 0; y < scene.OutputSize[1]; y += config.BlockSize {
        for x := 0; x < scene.OutputSize[0]; x += config.BlockSize {

            var task Task
            task.Area = image.Rect(x, y, x + config.BlockSize, y + config.BlockSize)
            taskList = append(taskList, task)
        }
    }

    return
}

func RenderWeldBlocks(scene Scene, blockList []Task) *image.RGBA {
    rect := image.Rect(0, 0, scene.OutputSize[0], scene.OutputSize[1])
    output := image.NewRGBA(rect)

    for _, elt := range blockList {
        draw.Draw(output, elt.Area, elt.Buffer, image.ZP, draw.Src)
    }

    return output
}

func RenderScene(config Config, scene Scene) (output *image.RGBA) {

    taskList := CreateRenderTasks(config, scene)
    blockList := make([]Task, 0)

    var wg sync.WaitGroup
    var blockCount = len(taskList)

    wg.Add(1)

    go func () {
        defer wg.Done()

        for len(taskList) > 0 {

            t := taskList[0]
            taskList = taskList[1:]

            RenderArea(config, scene, &t)
            blockList = append(blockList, t)

            fmt.Printf("done: %d/%d\r", len(blockList), blockCount)
        }
    }()

    wg.Wait()
    fmt.Printf("\noutputing now.\n")

    return RenderWeldBlocks(scene, blockList)
}
