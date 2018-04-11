package main

import (
    "math"
);

func ScreenPointToRay(scene Scene, x, y int) Ray {
    var r Ray

    width := float32(scene.OutputSize[0])
    height := float32(scene.OutputSize[1])
    aspectRatio := width / height
    fov := DEG2RAD * float32(scene.Camera.Fov) * 0.5

    right := scene.Camera.Up.Cross(scene.Camera.Forward).Neg().Normalize()

    var normCoords Vector
    normCoords.X = (float32(x) + 0.5) / width * 2. - 1.
    normCoords.Y = (float32(y) + 0.5) / height * 2. - 1.

    zNear := scene.Camera.ZNear / float32(math.Tan(float64(fov)))
    zFar := scene.Camera.ZFar

    var pCamera Vector
    pCamera.X = normCoords.X * (-zFar / zNear) * aspectRatio * -1.
    pCamera.Y = normCoords.Y * (-zFar / zNear)
    pCamera.Z = zFar

    pForward := scene.Camera.Forward.Normalize().MulScal(pCamera.Z)
    pUp := scene.Camera.Up.Normalize().MulScal(pCamera.Y)
    pRight := right.MulScal(pCamera.X)

    pWorld := pForward.Add(pUp).Add(pRight)

    r.Origin = scene.Camera.Position
    r.Direction = pWorld.Sub(r.Origin).Normalize()

    return r
}
