package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type WorldState []*WorldObject

type WorldObject struct {
	model        model.Model
	animator     model.Animator
	hasAnimation bool

	shader   shader.Shader
	modelMat mgl32.Mat4
}

func NewAnimatedWorldObject(m model.Model, anim model.Animator, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		model:        m,
		animator:     anim,
		hasAnimation: true,
		shader:       shader,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}
}

func NewWorldObject(m model.Model, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		model:        m,
		hasAnimation: false,
		shader:       shader,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}
}

func (o *WorldObject) Update(deltaTime float32) {
	o.animator.UpdateAnimation(deltaTime)
}

func (o WorldObject) Draw(proj mgl32.Mat4, view mgl32.Mat4) {
	o.shader.CheckShadersForChanges()
	o.shader.Use()

	o.shader.SetMatrix4("proj", proj)
	o.shader.SetMatrix4("view", view)
	o.shader.SetMatrix4("model", o.modelMat)

	if o.hasAnimation {
		transforms := o.animator.GetFinalBoneMatrices()
		for i, mat := range transforms {
			o.shader.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
		}
	}
	o.model.Draw(o.shader)
}
