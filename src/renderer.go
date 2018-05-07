package main

import (
    "fmt"
    "image"
    "image/draw"
    "math/rand"
    "sync"
    "time"
);

type Task struct {
    Area image.Rectangle
    Pixels *image.RGBA
    Depth [][]float64
    Width, Height int
};

func blitDepthBuffer(out, in [][]float64, rect image.Rectangle) {
    for y := rect.Min.Y; y < rect.Max.Y; y++ {
        for x := rect.Min.X; x < rect.Max.X; x++ {
            out[y][x] = in[y - rect.Min.Y][x - rect.Min.X]
        }
    }
}

func createDepthBuffer(width, height int) [][]float64 {
    depth := make([][]float64, height)
    for id, _ := range depth {
        depth[id] = make([]float64, width)
    }

    return depth
}

func createRenderTasks(config Config, scene Scene) (taskList []Task) {
    taskList = make([]Task, 0)

    for y := 0; y < scene.OutputSize[1]; y += config.BlockSize {
        for x := 0; x < scene.OutputSize[0]; x += config.BlockSize {

            var task Task

            task.Area.Min.X = x;
            task.Area.Min.Y = y;

            task.Width = MinInt(scene.OutputSize[0] - x, config.BlockSize)
            task.Height = MinInt(scene.OutputSize[1] - y, config.BlockSize)

            task.Area.Max.X = task.Area.Min.X + task.Width
            task.Area.Max.Y = task.Area.Min.Y + task.Height

            taskList = append(taskList, task)
        }
    }

    return
}

func renderWeldBlocks(scene Scene, blockList []Task) (frame Frame) {
    fmt.Printf("welding blocks...\r")
    rect := image.Rect(0, 0, scene.OutputSize[0], scene.OutputSize[1])
    frame.Pixels = image.NewRGBA(rect)
    frame.Depth = createDepthBuffer(scene.OutputSize[0], scene.OutputSize[1])
    frame.Width = scene.OutputSize[0]
    frame.Height = scene.OutputSize[1]

    for _, patch := range blockList {
        draw.Draw(frame.Pixels, patch.Area, patch.Pixels, image.ZP, draw.Src)
        blitDepthBuffer(frame.Depth, patch.Depth, patch.Area)
    }

    fmt.Printf("welding blocks...done\n")
    return frame
}

func renderArea(config Config, scene Scene, task *Task) {
    rect := image.Rect(0, 0, task.Width, task.Height)
    task.Pixels = image.NewRGBA(rect)
    task.Depth = createDepthBuffer(task.Width, task.Height)

    var frame Frame
    frame.Pixels = task.Pixels
    frame.Depth = task.Depth

    for y := 0; y < task.Height; y++ {
        for x := 0; x < task.Width; x++ {

            color, depth := Vector{ 0, 0, 0 }, 0.
            rays := ScreenPointToRaysDOF(config, scene, task.Area.Min.X + x,
                                                        task.Area.Min.Y + y)

            for _, r := range rays {
                sc, sd := RaytracerRender(scene, r)
                color = color.Add(sc)
                depth += sd
            }

            divisor := 1. / float64(len(rays))
            color = color.MulScal(divisor)
            depth = depth * divisor

            frame.Pixels.Set(x, y, VectorToRGBA(color))
            frame.Depth[y][x] = depth
        }
    }
}

func RenderScene(config Config, scene Scene) *image.RGBA {
    taskList := createRenderTasks(config, scene)
    blockList := make([]Task, 0)

    var wg sync.WaitGroup
    var blockCount = len(taskList)
    var mux sync.Mutex

    wg.Add(config.MaxThreads + 1)

    for i := 0; i < config.MaxThreads; i++ {
        go func (scene Scene) {
            defer wg.Done()

            src := rand.NewSource(time.Now().UnixNano())
            gen := rand.New(src)
            scene.Random = gen

            for len(taskList) > 0 {

                mux.Lock()
                if len(taskList) == 0 {
                    mux.Unlock()
                    break
                }
                t := taskList[0]
                taskList = taskList[1:]
                mux.Unlock()

                renderArea(config, scene, &t)

                mux.Lock()
                blockList = append(blockList, t)

                mux.Unlock()
            }
        }(scene)
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

    frame := renderWeldBlocks(scene, blockList)

    if config.ShowDebug {
        RasterizerDrawDebug(scene, frame)
    }

    return frame.Pixels
}
