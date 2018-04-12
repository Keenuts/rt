package main

import (
    "math"
);

func ScreenPointToRay(scene Scene, x, y int) Ray {
    var r Ray

    width := float64(scene.OutputSize[0])
    height := float64(scene.OutputSize[1])
    aspectRatio := width / height
    fov := DEG2RAD * scene.Camera.Fov * 0.5

    right := scene.Camera.Up.Cross(scene.Camera.Forward).Neg().Normalize()

    var normCoords Vector
    normCoords.X = (float64(x) + 0.5) / width * 2. - 1.
    normCoords.Y = (float64(y) + 0.5) / height * 2. - 1.

    zNear := scene.Camera.ZNear / math.Tan(fov)
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
