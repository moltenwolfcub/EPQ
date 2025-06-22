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

func NewWorldObjectFromModel(m model.Model, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		model:        m,
		hasAnimation: false,
		shader:       shader,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}
}

func NewWorldObject(modelFile string, hasAnimation bool, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	o := WorldObject{
		model:        model.NewModel(modelFile),
		hasAnimation: hasAnimation,
		shader:       shader,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}

	if o.hasAnimation {
		animation := model.NewAnimation(modelFile, &o.model)
		o.animator = model.NewAnimator(&animation)
	}

	return &o
}

func (o *WorldObject) Update(deltaTime float32) {
	o.animator.UpdateAnimation(deltaTime)
}

func (o WorldObject) Draw(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3) {
	o.shader.CheckShadersForChanges()
	o.shader.Use()

	o.shader.SetMatrix4("proj", proj)
	o.shader.SetMatrix4("view", view)
	o.shader.SetMatrix4("model", o.modelMat)
	o.shader.SetVec3("camera", mgl32.Vec3{}.Sub(camPos))

	//TMP
	// o.shader.SetVec3("material.ambient", mgl32.Vec3{0.25, 0.25, 0.25})
	// o.shader.SetVec3("material.diffuse", mgl32.Vec3{0.4, 0.4, 0.4})
	// o.shader.SetVec3("material.specular", mgl32.Vec3{0.774597, 0.774597, 0.774597})
	// o.shader.SetFloat("material.shininess", 0.6*128)

	o.shader.SetVec3("light.pos", mgl32.Vec3{-2, 5, -2})
	o.shader.SetVec3("light.ambient", mgl32.Vec3{0.2, 0.2, 0.2})
	o.shader.SetVec3("light.diffuse", mgl32.Vec3{1, 1, 1})
	o.shader.SetVec3("light.specular", mgl32.Vec3{1, 1, 1})

	if o.hasAnimation {
		transforms := o.animator.GetFinalBoneMatrices()
		for i, mat := range transforms {
			o.shader.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
		}
	}
	o.model.Draw(o.shader)
}
