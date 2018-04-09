package main

import (
    "fmt"
    "image"
    "image/draw"
    "sync"
    "time"
);

type Task struct {
    Area image.Rectangle
    Buffer *image.RGBA
};

func RenderArea(config Config, scene Scene, task *Task) {
    rect := image.Rect(0, 0, config.BlockSize, config.BlockSize)
    task.Buffer = image.NewRGBA(rect)

    for y := 0; y < config.BlockSize; y++ {
        for x := 0; x < config.BlockSize; x++ {

            r := ScreenPointToRay(scene, task.Area.Min.X + x, task.Area.Min.Y + y)
            px := TraceRay(config, scene, r)
            task.Buffer.Set(x, y, px)
        }
    }
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
    fmt.Printf("welding blocks...\r")
    rect := image.Rect(0, 0, scene.OutputSize[0], scene.OutputSize[1])
    output := image.NewRGBA(rect)

    for _, elt := range blockList {
        draw.Draw(output, elt.Area, elt.Buffer, image.ZP, draw.Src)
    }

    fmt.Printf("welding blocks...done\n")
    return output
}

func RenderScene(config Config, scene Scene) (output *image.RGBA) {

    taskList := CreateRenderTasks(config, scene)
    blockList := make([]Task, 0)

    var wg sync.WaitGroup
    var blockCount = len(taskList)
    var mux sync.Mutex

    wg.Add(config.MaxThreads + 1)

    for i := 0; i < config.MaxThreads; i++ {
        go func () {
            defer wg.Done()

            for len(taskList) > 0 {

                mux.Lock()
                if len(taskList) == 0 {
                    mux.Unlock()
                    break
                }
                t := taskList[0]
                taskList = taskList[1:]
                mux.Unlock()

                RenderArea(config, scene, &t)

                mux.Lock()
                blockList = append(blockList, t)

                mux.Unlock()
            }
        }()
    }

    go func() {
        defer wg.Done()

        for len(blockList) < blockCount {
            fmt.Printf("rendering... %d/%d\r", len(blockList), blockCount)
            time.Sleep(2e8)
        }
        fmt.Printf("rendering...done          \n")
    }()

    wg.Wait()

    return RenderWeldBlocks(scene, blockList)
}
