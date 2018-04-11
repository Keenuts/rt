package main

// Maths

type Vector struct {
    X, Y, Z float32
}

type Mat3 [3][3]float32

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
    Position Vector
    Rotation Vector
    Scale Vector
}

type Camera struct {
    Position, Forward, Up Vector
    Fov float32
    ZNear, ZFar float32
}

type Scene struct {
    Name string
    OutputSize [2]int

    Camera Camera
    Objects []Object
}

type Config struct {
    MaxThreads int
    BlockSize int
    OutputDir string
    SavePicture bool
    SaveReport bool
    ForceOutputName bool
    PictureName string
    ReportName string
    SceneName string
}

// Objects related

type Box struct {
    Min, Max Vector
    Volume float32
}

type Sphere struct {
    Center Vector
    Radius, Volume float32
}

type Triangle struct {
    Vertex [3]Vector
    Normals [3]Vector
}

type Model struct {
    Name string
    Triangles []Triangle
    Vertex []Vector
}

type Object struct {
    Name string

    Center Vector
    BoundingSphere Sphere
    BoundingBox Box

    Triangles []Triangle
    Tree KDTree
}


// Tracing

type Intersection struct {
    Position, Normal Vector
}

type Ray struct {
    Origin, Direction Vector
}

type KDTree struct {
    Left, Right *KDTree

    BoundingBox Box
    Triangles []Triangle
}
