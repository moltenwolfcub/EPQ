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
	modelMat  mgl32.Mat4
}

func NewWorldObject(obj gogl.Object, shader gogl.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		renderObj: obj,
		shader:    shader,
		modelMat:  mgl32.Translate3D(pos.Elem()),
	}
}

func (o WorldObject) Draw(proj mgl32.Mat4, view mgl32.Mat4) {
	o.shader.CheckShadersForChanges()
	o.shader.Use()

	o.shader.SetMatrix4("proj", proj)
	o.shader.SetMatrix4("view", view)

	o.renderObj.Draw(o.shader, o.modelMat)
}
