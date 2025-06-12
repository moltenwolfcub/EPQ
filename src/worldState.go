package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/gogl-utils"
)

type WorldState []*WorldObject

type WorldObject struct {
	model    *model.Model
	shader   gogl.Shader
	modelMat mgl32.Mat4

	uniformSetter func(gogl.Shader) gogl.Shader
}

func NewWorldObject(model *model.Model, shader gogl.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		model:    model,
		shader:   shader,
		modelMat: mgl32.Translate3D(pos.Elem()),
	}
}

func (o WorldObject) Draw(proj mgl32.Mat4, view mgl32.Mat4) {
	o.shader.CheckShadersForChanges()
	o.shader.Use()

	o.shader.SetMatrix4("proj", proj)
	o.shader.SetMatrix4("view", view)
	o.shader.SetMatrix4("model", o.modelMat)

	if o.uniformSetter != nil {
		o.shader = o.uniformSetter(o.shader)
	}
	o.model.Draw(o.shader)
}
