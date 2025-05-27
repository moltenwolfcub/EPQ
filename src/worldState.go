package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/gogl-utils"
)

type WorldState []*WorldObject

// for now just a wrapper for a gogl-util object but eventually writing my own to function with assimp
type WorldObject struct {
	RenderObj gogl.Object
	Shader    gogl.Shader
	Pos       mgl32.Vec3
}

func NewWorldObject(obj gogl.Object, shader gogl.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		RenderObj: obj,
		Shader:    shader,
		Pos:       pos,
	}
}
