package main

import (
    "math"
);

func getLambertDirectLight(scene Scene, info Intersection) (light Vector) {
    var ray Ray;
    ray.Origin = info.Position.Add(info.Normal.MulScal(BIAS))

    light = info.Object.Material.Emission

    for _, obj := range scene.Objects {
        if obj.ID == info.Object.ID {
            continue
        }

        if obj.Material.Emission.IsZero() {
            continue
        }

        ray.Direction = obj.Center.Sub(ray.Origin).Normalize()

        hit, it := IntersectObjects(ray, scene.Objects)

        if !hit || it.Object.ID != obj.ID {
            continue
        }

        factor := Max(0., info.Normal.Dot(ray.Direction))
        factor *= 1. / (math.Pi * 4. * it.Distance * it.Distance)
        factor *= 20. //FIXME: Power tweak for raytracer
        light = light.Add(obj.Material.Emission.MulScal(factor))

    }

    return
}

func getDefaultLighting(scene Scene, info Intersection) Vector {
    direction := Vector{ 0.1, -0.9, -0.1 }.Normalize()

    light := Vector{ 1., 1., 1.}.MulScal(info.Normal.Dot(direction.Neg()))
    light = MaxVec(Vector{ 0.1, 0.1, 0.1 }, light)

    return light
}

func rtTraceRayDepth(scene Scene, ray Ray, depth int) (Vector, float64) {
    if depth <= 0 {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    hit, info := IntersectObjects(ray, scene.Objects)
    if !hit {
        return Vector{0, 0, 0}, math.Inf(1)
    }

    diffuse := MtlGetDiffuse(info)

    var light Vector
    if scene.HasLight {
        light = getLambertDirectLight(scene, info)
    } else {
        light = getDefaultLighting(scene, info)
    }
    diffuse = diffuse.Scale(light)

    specularLevel := info.Object.Material.SpecularLevel

    var reflected Vector
    if specularLevel > 0. {
        var reflRay Ray
        reflRay.Origin = info.Position.Add(info.Normal.MulScal(EPSYLON))
        reflRay.Direction = Reflect(ray.Direction, info.Normal)

        reflected, _ = rtTraceRayDepth(scene, reflRay, depth - 1)
    }
    reflected = reflected.MulScal(specularLevel * 0.001)

    var fresnel float64
    var refracted Vector
    opacity := info.Object.Material.Opacity
    if opacity < 1. {
        _, rayOut := GetRefractedRay(scene, ray, info)
        refracted, _ = rtTraceRayDepth(scene, rayOut, depth - 1)
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

func RaytracerRender(scene Scene, ray Ray) (Vector, float64) {

    color, distance := rtTraceRayDepth(scene, ray, 8)
    return Saturate(color), distance
}

