package main

import (
    "math"
)

const EPSYLON = 1e-5
const BIAS = 1e-5
const DEG2RAD = ((math.Pi * 2.) / 360.)
const KDTREE_DEPTH_HINT = 18
const KDTREE_BUCKET_SIZE_HINT = 16

const MAX_PATH_DEPTH = 8
const PT_SAMPLE_PER_PIXEL = 256
