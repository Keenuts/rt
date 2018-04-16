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
    Pixels *image.RGBA
    Depth [][]float64
};

func BlitDepthBuffer(out, in [][]float64, rect image.Rectangle) {
    for y := rect.Min.Y; y < rect.Max.Y; y++ {
        for x := rect.Min.X; x < rect.Max.X; x++ {
            out[y][x] = in[y - rect.Min.Y][x - rect.Min.X]
        }
    }
}

func CreateDepthBuffer(width, height int) [][]float64 {
    depth := make([][]float64, height)
    for id, _ := range depth {
        depth[id] = make([]float64, width)
    }

    return depth
}

func RenderArea(config Config, scene Scene, task *Task) {
    rect := image.Rect(0, 0, config.BlockSize, config.BlockSize)
    task.Pixels = image.NewRGBA(rect)
    task.Depth = CreateDepthBuffer(config.BlockSize, config.BlockSize)

    var frame Frame
    frame.Pixels = task.Pixels
    frame.Depth = task.Depth

    for y := 0; y < config.BlockSize; y++ {
        for x := 0; x < config.BlockSize; x++ {

            r := ScreenPointToRay(scene, task.Area.Min.X + x, task.Area.Min.Y + y)
            color, depth := TraceRay(config, scene, r)

            frame.Pixels.Set(x, y, VectorToRGBA(color))
            frame.Depth[y][x] = depth
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

func RenderWeldBlocks(scene Scene, blockList []Task) (frame Frame) {
    fmt.Printf("welding blocks...\r")
    rect := image.Rect(0, 0, scene.OutputSize[0], scene.OutputSize[1])
    frame.Pixels = image.NewRGBA(rect)
    frame.Depth = CreateDepthBuffer(scene.OutputSize[0], scene.OutputSize[1])
    frame.Width = scene.OutputSize[0]
    frame.Height = scene.OutputSize[1]

    for _, elt := range blockList {
        draw.Draw(frame.Pixels, elt.Area, elt.Pixels, image.ZP, draw.Src)
        BlitDepthBuffer(frame.Depth, elt.Depth, elt.Area)
    }

    fmt.Printf("welding blocks...done\n")
    return frame
}

func RenderScene(config Config, scene Scene) *image.RGBA {

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

    frame := RenderWeldBlocks(scene, blockList)

    if config.ShowDebug {
        RasterizerDrawDebug(scene, frame)
    }

    return frame.Pixels
}
