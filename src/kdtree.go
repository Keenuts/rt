package main

import "math"

func TreeCreate(triangles []Triangle) (root KDTree) {

    root.Triangles = triangles
    queue := []*KDTree { &root }
    maxBuilt := uint64(math.Pow(2, KDTREE_DEPTH_HINT)) - 1
    built := uint64(0);

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]

        bounds, _ := MeshFindBounds(node.Triangles)
        max := bounds.Max.Sub(bounds.Min)
        node.BoundingBox = bounds
        built += 1

        if built > maxBuilt || len(node.Triangles) <= KDTREE_BUCKET_SIZE_HINT {
            node.Left = nil
            node.Right = nil
            continue
        }

        var lbucket, rbucket []Triangle
        if max.X >= max.Y && max.X >= max.Z {
            lbucket, rbucket = TreeCreateBuckets(0, node.Triangles)
        } else if max.Y >= max.X && max.Y >= max.Z {
            lbucket, rbucket = TreeCreateBuckets(1, node.Triangles)
        } else {
            lbucket, rbucket = TreeCreateBuckets(2, node.Triangles)
        }

        var lchild, rchild KDTree

        if len(lbucket) > 0 {
            lchild.Triangles = lbucket
            node.Left = &lchild
            queue = append(queue, &lchild)
        }

        if len(rbucket) > 0 {
            rchild.Triangles = rbucket
            node.Right = &rchild
            queue = append(queue, &rchild)
        }

        node.Triangles = make([]Triangle, 0)
    }

    return
}

func TreeCreateBuckets(axis int, triangles []Triangle) (left, right []Triangle) {
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

    return
}
