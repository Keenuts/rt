package main

func (tree KDTree) Insert(triangles []Triangle) KDTree {
    if tree.Left != nil || tree.Right != nil {
        panic("Called insert twice on the same node")
    }

    bounds, _ := MeshFindBounds(triangles)
    max := bounds.Max.Sub(bounds.Min)

    tree.BoundingBox = bounds

    if len(triangles) < 16 {
        tree.Triangles = triangles
    } else {
        if max.X >= max.Y && max.X >= max.Z {
            tree = tree.TreeInsertAxis(0, triangles)
        } else if max.Y >= max.X && max.Y >= max.Z {
            tree = tree.TreeInsertAxis(1, triangles)
        } else {
            tree = tree.TreeInsertAxis(2, triangles)
        }
    }

    return tree
}

func (tree KDTree) TreeInsertAxis(axis int, triangles []Triangle) KDTree {
    var left, right []Triangle

    middleV := MeshFindCenter(triangles)
    middle := [3]float32 { middleV.X, middleV.Y, middleV.Z }

    for _, tri := range triangles {
        mtriV := TriangleFindCenter(tri)
        mtri := [3]float32 { mtriV.X, mtriV.Y, mtriV.Z }

        if mtri[axis] < middle[axis] {
            left = append(left, tri)
        } else {
            right = append(right, tri)
        }
    }

    if len(left) > 0 {
        tree.Left = new(KDTree)
        *tree.Left = tree.Left.Insert(left)
    }

    if len(right) > 0 {
        tree.Right = new(KDTree)
        *tree.Right = tree.Right.Insert(right)
    }

    return tree
}
