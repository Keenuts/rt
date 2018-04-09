package main

import (
    "math"
);

func ScreenPointToRay(scene Scene, x, y int) Ray {
    var r Ray

    width := float64(scene.OutputSize[0])
    height := float64(scene.OutputSize[1])
    aspectRatio := width / height
    fov := DEG2RAD * float64(scene.Camera.Fov) * 0.5

    right := scene.Camera.Up.Cross(scene.Camera.Forward).Neg()

    spX := ((float64(x) + 0.5) / width) * 2. - 1.
    spY := ((float64(y) + 0.5) / height) * 2. - 1.
    spY *= -1; // Y coordinate is flipped between img and world

    spX *= math.Tan(fov) * aspectRatio
    spY *= math.Tan(fov)

    middle := scene.Camera.Position.Add(scene.Camera.Forward)
    middle = middle.Add(scene.Camera.Up.MulScal(float32(spY)))
    middle = middle.Add(right.MulScal(float32(spX)))

    r.Origin = scene.Camera.Position
    r.Direction = middle.Sub(r.Origin).Normalize()

    return r
}
