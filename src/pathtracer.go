package main

import (
    "math"
);

const SAMPLE_PER_PIXEL = 32
const MAX_PATH_DEPTH = 8

func traceRefraction(scene Scene, ray Ray, info Intersection, depth int) Vector {
    ratioOut := info.Object.Material.Refraction
    ratioIn := 1.0 / ratioOut

    posA := info.Position
    norA := info.Normal

    var ray2 Ray
    ray2.Origin = posA.Add(norA.MulScal(-1. * EPSYLON))
    ray2.Direction = Refract(ray.Direction, info.Normal, ratioIn)
    ray2.InvertCulling = true

    hit2, info2 := Intersect(ray2, info.Object)
    ray2.InvertCulling = false

    if hit2 && info.Object.ID == info2.Object.ID {
        posB := info2.Position
        norB := info2.Normal

        ray2.Origin = posB.Add(norB.MulScal(1. * EPSYLON))
        ray2.Direction = Refract(ray2.Direction, info2.Normal.Neg(), ratioOut)
    }

    refracted, _ := traceRayDepth(scene, ray2, depth - 1)
    return refracted
}

func backtraceLight(scene Scene, info Intersection) (Vector) {
    if !info.Object.Material.Emission.IsZero() {
        return info.Object.Material.Emission
    }

    var ray Ray
    var hit bool
    mask := Vector{ 1, 1, 1 }
    light := Vector{ 0, 0, 0 }

    for i := 0; i < MAX_PATH_DEPTH; i++ {
        ray.Origin = info.Position.Add(info.Normal.MulScal(EPSYLON))
        ray.Direction = RandomHemisphereVector(scene.Random, info.Normal)

        hit, info = IntersectObjects(ray, scene.Objects)
        if !hit {
            break;
        }

        emissive := info.Object.Material.Emission
        if !emissive.IsZero() {
            light = emissive.Scale(mask)
            break
        }

        if info.Object.Material.Opacity >= 1. - EPSYLON {
            mask = mask.Scale(info.Object.Material.Diffuse)
            continue
        }

        ratioOut := info.Object.Material.Refraction
        ratioIn := 1.0 / ratioOut

        posA := info.Position
        norA := info.Normal

        var rayIn Ray
        rayIn.Origin = posA.Add(norA.MulScal(-1. * EPSYLON))
        rayIn.Direction = Refract(ray.Direction, info.Normal, ratioIn)
        rayIn.InvertCulling = true

        hitB, infoB := Intersect(rayIn, info.Object)

        if !hitB || info.Object.ID != infoB.Object.ID {
            ray = rayIn
            continue;
        }

        posB := infoB.Position
        norB := infoB.Normal

        var rayOut Ray
        rayOut.Origin = posB.Add(norB.MulScal(-1. * EPSYLON))
        rayOut.Direction = Refract(rayIn.Direction, infoB.Normal.Neg(), ratioOut)
        rayOut.InvertCulling = false

        ray = rayOut
    }

    return light
}

func traceRayDepth(scene Scene, ray Ray, depth int) (Vector, float64) {
    if depth <= 0 {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    hit, info := IntersectObjects(ray, scene.Objects)
    if !hit {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    diffuse := MtlGetDiffuse(info)
    light := backtraceLight(scene, info)
    diffuse = diffuse.Scale(light)

    specularLevel := info.Object.Material.SpecularLevel

    var reflected Vector
    if specularLevel > 0. {
        var reflRay Ray
        reflRay.Origin = info.Position.Add(info.Normal.MulScal(EPSYLON))
        reflRay.Direction = Reflect(ray.Direction, info.Normal)
        reflected, _ = traceRayDepth(scene, reflRay, depth - 1)
    }
    reflected = reflected.MulScal(specularLevel * 0.001)

    var fresnel float64
    var refracted Vector
    opacity := info.Object.Material.Opacity
    if opacity < 1. {
        refracted = traceRefraction(scene, ray, info, depth)
        fresnel = Fresnel(ray.Direction, info.Normal, 1.0 / 1.5161)
    } else {
        refracted = diffuse
        fresnel = 1.
    }
    refracted = refracted.MulScal(1. - fresnel)

    output := diffuse.MulScal(1. - specularLevel * 0.001)
    output = output.Add(reflected).MulScal(fresnel)
    output = output.Add(refracted)

    return output, info.Distance
}

func TraceRay(scene Scene, ray Ray) (Vector, float64) {

    color := Vector{ 0, 0, 0 }
    distance := 0.

    for i := 0; i < SAMPLE_PER_PIXEL; i++ {
        sample, dist := traceRayDepth(scene, ray, 8)
        if i == 0 {
            distance = dist
        }

        color = color.Add(sample.MulScal(1. / float64(SAMPLE_PER_PIXEL)))
    }

    return Saturate(color), distance
}

