#!/bin/python3

import json
import os
import math

DEG2RAD = ((math.pi * 2.) / 360.)
EPSYLON = 1e-5

class Vector:
    def __init__(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

    def Add(self, vec):
        return Vector(
            self.x + vec.x,
            self.y + vec.y,
            self.z + vec.z
        )

    def Scale(self, vec):
        return Vector(
            self.x * vec.x,
            self.y * vec.y,
            self.z * vec.z
        )

    def Sub(self, vec):
        return Vector(
            self.x - vec.x,
            self.y - vec.y,
            self.z - vec.z
        )

    def MulScal(self, s):
        return Vector(
            self.x * s,
            self.y * s,
            self.z * s,
        )

    def Magnitude(self):
        return math.sqrt(self.x * self.x + self.y * self.y + self.z * self.z)

    def Normalized(self):
        norm = self.Magnitude()
        if norm < EPSYLON and norm > -EPSYLON:
            return Vector(0, 0, 0)
    
        return Vector( self.x / norm, self.y / norm, self.z / norm )

BACKUP_CONFIG_FILE = ".config.json"
BACKUP_SCENE_FILE = ".scene.json"
CONFIG_FILE = "config.json"
SCENE_FILE = "scenes/teapot-simplified.json"
TMP_SCENE_FILE = ".tmp.json"

DISTANCE = 10
ORIGIN = Vector(0, 0, 0)
STEPS = 32


def CopyFile(src, dst):
    srcfile = open(src)
    dstfile = open(dst, "w")

    dstfile.write(srcfile.read())

def SaveFiles():
    CopyFile(SCENE_FILE, BACKUP_SCENE_FILE)
    CopyFile(CONFIG_FILE, BACKUP_CONFIG_FILE)

def RestoreFiles():
    CopyFile(BACKUP_SCENE_FILE, SCENE_FILE)
    CopyFile(BACKUP_CONFIG_FILE, CONFIG_FILE)

    os.remove(BACKUP_SCENE_FILE)
    os.remove(BACKUP_CONFIG_FILE)
    os.remove(TMP_SCENE_FILE)

def SetFilename(config, filename):
    config['pictureName'] = filename
    return config

def WriteScene(scene):
    with open(TMP_SCENE_FILE, "w") as f:
        f.write(json.dumps(scene))

def WriteConfig(config):
    with open(CONFIG_FILE, "w") as f:
        f.write(json.dumps(config))

def SetVector(v):
    return { "x": v.x, "y": v.y, "z": v.z}

def GetVectors(origin, distance, angle):
    angle = angle * DEG2RAD

    s = Vector(0, 0, -1)
    v = Vector(0, 0, 0)
    v.x = math.cos(angle) * s.x - math.sin(angle) * s.z
    v.z = math.sin(angle) * s.x + math.cos(angle) * s.z

    v = v.Normalized().MulScal(distance)

    direction = origin.Sub(v)

    direction = direction.Normalized()
    position = origin.Sub(direction.MulScal(distance))

    return SetVector(position), SetVector(direction)

def SetCamera(scene, origin, distance, angle):
    position, forward = GetVectors(ORIGIN, DISTANCE, angle)

    print(position, forward)
    scene['Camera']["position"] = position
    scene['Camera']["forward"] = forward

    return scene

def Render():
    try:
        os.system("./rt")
    except:
        return False
    return True

def main():
    SaveFiles()
    config = None
    scene = None
    render_id = 0

    with open(CONFIG_FILE) as f:
        config = json.loads(f.read())
    with open(SCENE_FILE) as f:
        scene = json.loads(f.read())


    config['forceOutputName'] = True
    config['outputDir'] = "turnover"
    config['savePicture'] = True
    config['saveReport'] = False
    config['sceneName'] = TMP_SCENE_FILE


    angle = 0
    for i in range(STEPS):
        scene = SetCamera(scene, ORIGIN, DISTANCE, angle)
        config = SetFilename(config, format("turnover-%03d.png" % i))

        WriteConfig(config)
        WriteScene(scene)

        angle += 360 / STEPS
        Render()


    RestoreFiles()

main()
