package main

import (
    "image"
    "image/color"
    "math/rand"
);

// Maths

type Vector struct {
    X, Y, Z float64
}

type Mat3 [3][3]float64

// Scene related

type RenderInfo struct {
    SceneName string
    Date string
    OutputSize [2]int
    Duration string
    Threads int
    Config Config
};

type SceneObject struct {
    ObjectID int
    Position, Rotation, Scale Vector
    MaterialLibID int
    MaterialName string
    DebugName string
}

type Camera struct {
    Position, Forward, Up Vector
    Fov float64
    ZNear, ZFar float64
}

type Scene struct {
    Name string
    OutputSize [2]int

    Camera Camera
    Objects []Object
    Random *rand.Rand
}

type Config struct {
    BlockSize int
    ForceOutputName bool
    MaxThreads int
    OutputDir string
    PictureName string
    ReportName string
    SavePicture bool
    SaveReport bool
    SceneName string
    ShowDebug bool
}

// Objects related

type Box struct {
    Min, Max Vector
    Volume float64
}

type Sphere struct {
    Center Vector
    Radius, Volume float64
}

type Triangle struct {
    Vertex, Normals, UV [3]Vector
}

type Model struct {
    Name string
    Triangles []Triangle
    Vertex []Vector
}

type Object struct {
    ID int
    Name string

    Center Vector
    BoundingSphere Sphere
    BoundingBox Box

    Triangles []Triangle
    Tree KDTree
    Material Material
}

// Material related

type MaterialLib map[string]Material

type Texture struct {
    Pixels *image.RGBA
    Width, Height int
}

type Material struct {
    Diffuse, Specular, Emission Vector
    Opacity float64
    Refraction float64
    SpecularLevel float64
    DiffuseTex *Texture
}

// Tracing & Rasterizing

type Intersection struct {
    Position, Normal, UV Vector
    Distance float64
    Object Object
}

type Ray struct {
    Origin, Direction Vector
    InvertCulling bool
}

type KDTree struct {
    Left, Right *KDTree

    BoundingBox Box
    Triangles []Triangle
}

type Frame struct {
    Pixels *image.RGBA
    Depth [][]float64
    Width, Height int
}

type Point struct {
    Position Vector
    Color color.RGBA
}

type Line struct {
    A, B Vector
    Color color.RGBA
}
