package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
);

func DrawLine2D(a, b Vector, color color.RGBA, outputSize [2]int, output *image.RGBA) {
    x0, x1, y0, y1 := int(a.X), int(b.X), int(a.Y), int(b.Y)

    dx := AbsInt(x1 - x0)
    dy := AbsInt(y1 - y0)
    err := dx - dy

    sx := 1
    if x0 >= x1 {
        sx = -1
    }

    sy := 1
    if y0 >= y1 {
        sy = -1
    }

    for ; x0 != x1 || y0 != y1 ; {
        if x0 >= 0 && x0 < outputSize[0] && y0 >= 0 && y0 < outputSize[1] {
            output.Set(x0, y0, color)
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

func WorldToCamera(p Vector, camera Camera) Vector {
    var tR Mat3
    right := camera.Up.Cross(camera.Forward).Neg()

    tR[0] = [3]float32{ right.X, camera.Up.X, camera.Forward.X }
    tR[1] = [3]float32{ right.Y, camera.Up.Y, camera.Forward.Y }
    tR[2] = [3]float32{ right.Z, camera.Up.Z, camera.Forward.Z }

    return tR.Mul(p.Sub(camera.Position))
}

func CameraToScreen(scene Scene, pCamera Vector) Vector {
    fov := scene.Camera.Fov * .5 * DEG2RAD
    fovTan := float32(math.Tan(float64(fov)))
    aspectRatio := float32(scene.OutputSize[0]) / float32(scene.OutputSize[1])

    var canvasSize Vector
    canvasSize.X = fovTan * scene.Camera.ZNear * aspectRatio
    canvasSize.Y = fovTan * scene.Camera.ZNear;
    canvasSize = canvasSize.MulScal(2)

    var pScreen Vector
    pScreen.X = scene.Camera.ZNear * pCamera.X / pCamera.Z
    pScreen.Y = scene.Camera.ZNear * pCamera.Y / -pCamera.Z
    pScreen.Z = -pCamera.Z

    var pNDC Vector
    pNDC.X = (2. * pScreen.X) / canvasSize.X
    pNDC.Y = (2. * pScreen.Y) / canvasSize.Y
    pNDC.Z = pScreen.Z

    var pImage Vector
    pImage.X = (pNDC.X + 1.) * .5 * float32(scene.OutputSize[0])
    pImage.Y = (pNDC.Y + 1.) * .5 * float32(scene.OutputSize[1])
    pImage.Z = pScreen.Z

    return pImage
}

func DrawGizmoLine(scene Scene, a, b Vector, color color.RGBA, output *image.RGBA) {
    ca := WorldToCamera(a, scene.Camera)
    cb := WorldToCamera(b, scene.Camera)

    sa := CameraToScreen(scene, ca)
    sb := CameraToScreen(scene, cb)

    DrawLine2D(sa, sb, color, scene.OutputSize, output)
}

func RasterizerDrawBoundingBox(scene Scene, box Box, col color.RGBA, out *image.RGBA) {
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
            DrawGizmoLine(scene, vtx[i + j * 4], vtx[(i + 1) % 4 + j * 4], col, out)
        }
    }

    for i := 0; i < 4; i++ {
        DrawGizmoLine(scene, vtx[i], vtx[i + 4], col, out)
    }
}

func RasterizerDrawDebug(scene Scene, output *image.RGBA) {
    fmt.Printf("drawing debug informations...")

    red := color.RGBA{ 255, 0, 0, 255 }

    box := Box{ Vector{-1, -1, -1}, Vector{ 1, 1, 1 }, 0 }
    RasterizerDrawBoundingBox(scene, box, red, output)

    fmt.Printf("done\n")
}
