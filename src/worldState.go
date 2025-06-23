package main

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type WorldState struct {
	Objects []*WorldObject
	Lights  []Light

	lightingSSBO uint32
}

func NewWorldState() *WorldState {
	state := &WorldState{}

	gl.GenBuffers(1, &state.lightingSSBO)

	return state
}

func (s *WorldState) BindLights() {
	//padding because vec3s need to be aligned to 16 bytes in SSBOS
	type internalLight struct {
		pos       [3]float32
		_pad1     float32
		ambient   [3]float32
		_pad2     float32
		diffuse   [3]float32
		_pad3     float32
		specular  [3]float32
		constant  float32
		linear    float32
		quadratic float32
		_pad4     float32
		_pad5     float32
	}

	internalLights := []internalLight{}
	for _, light := range s.Lights {
		internalLights = append(internalLights, internalLight{
			pos: light.Pos,

			ambient:  light.Ambient,
			diffuse:  light.Diffuse,
			specular: light.Specular,

			constant:  light.ConstantAttenuation,
			linear:    light.LinearAttenuation,
			quadratic: light.QuadraticAttenuation,
		})
	}

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, s.lightingSSBO)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(internalLights)*int(unsafe.Sizeof(internalLight{})), gl.Ptr(internalLights), gl.STATIC_DRAW)
}

type WorldObject struct {
	state *WorldState

	model        model.Model
	animator     model.Animator
	hasAnimation bool

	shader   shader.Shader
	modelMat mgl32.Mat4
}

func NewWorldObjectFromModel(state *WorldState, m model.Model, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		state:        state,
		model:        m,
		hasAnimation: false,
		shader:       shader,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}
}

func NewWorldObject(state *WorldState, modelFile string, hasAnimation bool, shader shader.Shader, pos mgl32.Vec3) *WorldObject {
	o := WorldObject{
		state:        state,
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

	if o.state.lightingSSBO != 0 {
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, o.state.lightingSSBO)
	} else {
		panic("Lighting was never bound on state")
	}

	if o.hasAnimation {
		transforms := o.animator.GetFinalBoneMatrices()
		for i, mat := range transforms {
			o.shader.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
		}
	}
	o.model.Draw(o.shader)
}

type Light struct {
	Pos mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32
}
