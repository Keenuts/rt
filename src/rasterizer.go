package main

import (
    "fmt"
    "image/color"
    "math"
);

func drawLine2D(line Line, frame Frame) {
    x0, x1, y0, y1 := int(line.A.X), int(line.B.X), int(line.A.Y), int(line.B.Y)
    direction := line.A.Sub(line.B).Normalize()

    dx, dy := AbsInt(x1 - x0), AbsInt(y1 - y0)
    err := dx - dy

    sx, sy := 1, 1
    if x0 >= x1 {
        sx = -1
    }
    if y0 >= y1 {
        sy = -1
    }

    for ; x0 != x1 || y0 != y1 ; {
        if x0 >= 0 && x0 < frame.Width && y0 >= 0 && y0 < frame.Height {
            dist := VecFromInt(x0, y0, 0).Sub(VecFromInt(x1, y1, 0)).Magnitude()
            pt := line.B.Add(direction.MulScal(dist))
            if frame.Depth[y0][x0] > pt.Z {
                frame.Pixels.Set(x0, y0, line.Color)
            }
        }

        ed := 2 * err
        if ed > -dy {
            err -= dy
            x0 += sx
        }

        if ed < dx {
            err += dx
            y0 += sy
        }
    }
}

func worldToCamera(p Vector, camera Camera) Vector {
    var tR Mat3
    right := camera.Up.Cross(camera.Forward).Neg()

    tR[0] = [3]float64{ right.X, camera.Up.X, camera.Forward.X }
    tR[1] = [3]float64{ right.Y, camera.Up.Y, camera.Forward.Y }
    tR[2] = [3]float64{ right.Z, camera.Up.Z, camera.Forward.Z }

    return tR.Mul(p.Sub(camera.Position))
}

func cameraToScreen(scene Scene, pCamera Vector) Vector {
    fov := scene.Camera.Fov * .5 * DEG2RAD
    fovTan := math.Tan(fov)
    aspectRatio := float64(scene.OutputSize[0]) / float64(scene.OutputSize[1])
    flippedRatio := false

    if aspectRatio > 1. {
        aspectRatio = float64(scene.OutputSize[1]) / float64(scene.OutputSize[0])
        flippedRatio = true
    }

    var canvasSize Vector
    canvasSize.X = fovTan * scene.Camera.ZNear
    canvasSize.Y = fovTan * scene.Camera.ZNear
    canvasSize = canvasSize.MulScal(2)

    if flippedRatio {
        canvasSize.Y *= aspectRatio
    } else {
        canvasSize.X *= aspectRatio
    }


    pCamera.Z = Max(scene.Camera.ZNear, Min(scene.Camera.ZFar, pCamera.Z))
    var pScreen Vector
    pScreen.X = scene.Camera.ZNear * pCamera.X / pCamera.Z
    pScreen.Y = scene.Camera.ZNear * pCamera.Y / -pCamera.Z
    pScreen.Z = pCamera.Z

    var pNDC Vector
    pNDC.X = (2. * pScreen.X) / canvasSize.X
    pNDC.Y = (2. * pScreen.Y) / canvasSize.Y
    pNDC.Z = pScreen.Z

    var pImage Vector
    pImage.X = (pNDC.X + 1.) * .5 * float64(scene.OutputSize[0])
    pImage.Y = (pNDC.Y + 1.) * .5 * float64(scene.OutputSize[1])
    pImage.Z = pScreen.Z

    return pImage
}

func drawGizmoLine(scene Scene, line Line, frame Frame) {
    ca := worldToCamera(line.A, scene.Camera)
    cb := worldToCamera(line.B, scene.Camera)

    sa := cameraToScreen(scene, ca)
    sb := cameraToScreen(scene, cb)

    var line2D Line
    line2D.A = sa
    line2D.B = sb
    line2D.Color = line.Color
    drawLine2D(line2D, frame)
}

func rasterizerDrawBoundingBox(scene Scene, box Box, col color.RGBA, frame Frame) {
    var vtx [8]Vector

    vtx[0] = box.Min.Scale(Vector{1, 1, 1}).Add(box.Max.Scale(Vector{0, 0, 0}))
    vtx[1] = box.Min.Scale(Vector{0, 1, 1}).Add(box.Max.Scale(Vector{1, 0, 0}))
    vtx[2] = box.Min.Scale(Vector{0, 0, 1}).Add(box.Max.Scale(Vector{1, 1, 0}))
    vtx[3] = box.Min.Scale(Vector{1, 0, 1}).Add(box.Max.Scale(Vector{0, 1, 0}))

    vtx[4] = box.Min.Scale(Vector{1, 1, 0}).Add(box.Max.Scale(Vector{0, 0, 1}))
    vtx[5] = box.Min.Scale(Vector{0, 1, 0}).Add(box.Max.Scale(Vector{1, 0, 1}))
    vtx[6] = box.Min.Scale(Vector{0, 0, 0}).Add(box.Max.Scale(Vector{1, 1, 1}))
    vtx[7] = box.Min.Scale(Vector{1, 0, 0}).Add(box.Max.Scale(Vector{0, 1, 1}))

    for j := 0; j < 2; j++ {
        for i := 0; i < 4; i++ {
            line := Line{ vtx[i + j * 4], vtx[(i + 1) % 4 + j * 4], col }
            drawGizmoLine(scene, line, frame)
        }
    }

    for i := 0; i < 4; i++ {
        line := Line{ vtx[i], vtx[i + 4], col }
        drawGizmoLine(scene, line, frame)
    }
}

func rasterizerDrawTree(scene Scene, root KDTree, col color.RGBA, frame Frame) {
    queue := []*KDTree { &root }

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]

        rasterizerDrawBoundingBox(scene, node.BoundingBox, col, frame)

        if node.Left != nil {
            queue = append(queue, node.Left)
        }

        if node.Right != nil {
            queue = append(queue, node.Right)
        }
    }
}

func rasterizerDrawPoint(scene Scene, point Point, frame Frame) {
    const SCALE = 0.02

    dist := scene.Camera.Position.Sub(point.Position).Magnitude()
    size := math.Tan(scene.Camera.Fov * .5) * dist

    min := point.Position.AddScal(-size * .5 * SCALE)
    max := point.Position.AddScal(size * .5 * SCALE)
    rasterizerDrawBoundingBox(scene, Box{ min, max, 0 }, point.Color, frame)
}

// This function used built-in rasterizer to draw debug gizmos.
// To be drawn, the showDebug boolean must be set to True in the config file.
func RasterizerDrawDebug(scene Scene, frame Frame) {
    fmt.Printf("drawing debug informations...")

    green := color.RGBA{ 0, 255, 0, 255 }
    for _, o := range scene.Objects {
        rasterizerDrawTree(scene, o.Tree, green, frame)
    }

    red := color.RGBA{ 255, 0, 0, 255 }
    for _, o := range scene.Objects {
        rasterizerDrawBoundingBox(scene, o.BoundingBox, red, frame)
    }

    fmt.Printf("done\n")
}
