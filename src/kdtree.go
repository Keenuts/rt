package main

import "math"

func treeCreateBuckets(axis int, triangles []Triangle) (left, right []Triangle) {
    middleV := MeshFindCenter(triangles)
    middle := [3]float64 { middleV.X, middleV.Y, middleV.Z }

    for _, tri := range triangles {
        mtriV := TriangleFindCenter(tri)
        mtri := [3]float64 { mtriV.X, mtriV.Y, mtriV.Z }

        if mtri[axis] < middle[axis] {
            left = append(left, tri)
        } else {
            right = append(right, tri)
        }
    }

    return
}

func TreeCreate(triangles []Triangle) (root KDTree) {

    root.Value = triangles
    queue := []*KDTree { &root }
    maxBuilt := uint64(math.Pow(2, KDTREE_DEPTH_HINT)) - 1
    built := uint64(0);

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]

        tris := node.Value.([]Triangle)
        bounds, _ := MeshFindBounds(tris)
        max := bounds.Max.Sub(bounds.Min)
        node.BoundingBox = bounds
        built += 1

        if built > maxBuilt || len(tris) <= KDTREE_BUCKET_SIZE_HINT {
            node.Left = nil
            node.Right = nil
            continue
        }

        var lbucket, rbucket []Triangle
        if max.X >= max.Y && max.X >= max.Z {
            lbucket, rbucket = treeCreateBuckets(0, tris)
        } else if max.Y >= max.X && max.Y >= max.Z {
            lbucket, rbucket = treeCreateBuckets(1, tris)
        } else {
            lbucket, rbucket = treeCreateBuckets(2, tris)
        }

        var lchild, rchild KDTree

        if len(lbucket) > 0 {
            lchild.Value = lbucket
            node.Left = &lchild
            queue = append(queue, &lchild)
        }

        if len(rbucket) > 0 {
            rchild.Value = rbucket
            node.Right = &rchild
            queue = append(queue, &rchild)
        }

        node.Value = make([]Triangle, 0)
    }

    return
}
