package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/gogl-utils"
)

type WorldState []*WorldObject

// for now just a wrapper for a gogl-util object but eventually writing my own to function with assimp
type WorldObject struct {
	renderObj gogl.Object
	shader    gogl.Shader
	pos       mgl32.Vec3
}

func NewWorldObject(obj gogl.Object, shader gogl.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		renderObj: obj,
		shader:    shader,
		pos:       pos,
	}
}

func (o WorldObject) Draw(proj mgl32.Mat4, view mgl32.Mat4) {
	o.shader.CheckShadersForChanges()
	o.shader.Use()

	o.shader.SetMatrix4("proj", proj)
	o.shader.SetMatrix4("view", view)

	modelMat := mgl32.Translate3D(o.pos.Elem())

	o.renderObj.Draw(o.shader, modelMat)
}
