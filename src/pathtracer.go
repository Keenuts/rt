package main

import (
    "math"
);

func ptBacktraceLight(scene Scene, info Intersection) (Vector) {
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

func ptTraceRayDepth(scene Scene, ray Ray, depth int) (Vector, float64) {
    if depth <= 0 {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    hit, info := IntersectObjects(ray, scene.Objects)
    if !hit {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    diffuse := MtlGetDiffuse(info)
    light := ptBacktraceLight(scene, info)
    diffuse = diffuse.Scale(light)

    specularLevel := info.Object.Material.SpecularLevel

    var reflected Vector
    if specularLevel > 0. {
        var reflRay Ray
        reflRay.Origin = info.Position.Add(info.Normal.MulScal(EPSYLON))
        reflRay.Direction = Reflect(ray.Direction, info.Normal)
        reflected, _ = ptTraceRayDepth(scene, reflRay, depth - 1)
    }
    reflected = reflected.MulScal(specularLevel * 0.001)

    var fresnel float64
    var refracted Vector
    opacity := info.Object.Material.Opacity
    if opacity < 1. {
        _, rayOut := GetRefractedRay(scene, ray, info)
        refracted, _ = ptTraceRayDepth(scene, rayOut, depth - 1)
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

func PathtracerRender(scene Scene, ray Ray) (Vector, float64) {

    color := Vector{ 0, 0, 0 }
    distance := 0.

    for i := 0; i < PT_SAMPLE_PER_PIXEL; i++ {
        sample, dist := ptTraceRayDepth(scene, ray, 8)
        if i == 0 {
            distance = dist
        }

        color = color.Add(sample.MulScal(1. / float64(PT_SAMPLE_PER_PIXEL)))
    }

    return Saturate(color), distance
}

