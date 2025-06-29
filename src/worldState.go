package main

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type WorldState struct {
	Objects []*WorldObject
	Lights  []Light

	lightingSSBO uint32
}

func NewWorldState() *WorldState {
	state := &WorldState{
		Objects: make([]*WorldObject, 0),
	}

	gl.GenBuffers(1, &state.lightingSSBO)

	return state
}

func (s *WorldState) FinaliseLoad() {
	for _, o := range s.Objects {
		o.finaliseLoad()
	}
}

func (s *WorldState) BindLights() {
	//padding because vec3s need to be aligned to 16 bytes in SSBOS
	type internalLight struct {
		pos         [3]float32
		_pad1       float32
		dir         [3]float32
		lightType   int32
		ambient     [3]float32
		_pad2       float32
		diffuse     [3]float32
		_pad3       float32
		specular    [3]float32
		constant    float32
		linear      float32
		quadratic   float32
		cutoff      float32
		outerCutoff float32
	}

	internalLights := []internalLight{}
	for _, light := range s.Lights {
		var il internalLight

		switch l := light.(type) {
		case PointLight:
			il = internalLight{
				lightType: 0,
				pos:       l.Pos,
				ambient:   l.Ambient,
				diffuse:   l.Diffuse,
				specular:  l.Specular,
				constant:  l.ConstantAttenuation,
				linear:    l.LinearAttenuation,
				quadratic: l.QuadraticAttenuation,
			}
		case DirLight:
			il = internalLight{
				lightType: 1,
				dir:       l.Dir,
				ambient:   l.Ambient,
				diffuse:   l.Diffuse,
				specular:  l.Specular,
			}
		case SpotLight:
			il = internalLight{
				lightType:   2,
				pos:         l.Pos,
				dir:         l.Dir,
				ambient:     l.Ambient,
				diffuse:     l.Diffuse,
				specular:    l.Specular,
				constant:    l.ConstantAttenuation,
				linear:      l.LinearAttenuation,
				quadratic:   l.QuadraticAttenuation,
				cutoff:      float32(math.Cos(float64(mgl32.DegToRad(l.cutoff)))),
				outerCutoff: float32(math.Cos(float64(mgl32.DegToRad(l.outerCutoff)))),
			}
		}

		internalLights = append(internalLights, il)
	}

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, s.lightingSSBO)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(internalLights)*int(unsafe.Sizeof(internalLight{})), gl.Ptr(internalLights), gl.STATIC_DRAW)
}

type WorldObject struct {
	state *WorldState

	model        *model.Model
	animator     *model.Animator
	hasAnimation bool

	shader   shader.Shader
	modelMat mgl32.Mat4

	normalShader shader.Shader
}

func NewWorldObjectFromModel(state *WorldState, m *model.Model, s shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		state:        state,
		model:        m,
		hasAnimation: false,
		shader:       s,
		modelMat:     mgl32.Translate3D(pos.Elem()),
		normalShader: shader.NewEmbeddedShaderVFG(assets.NormViewVert, assets.NormViewFrag, assets.NormViewGeom),
	}
}

func NewWorldObject(state *WorldState, modelFile string, hasAnimation bool, s shader.Shader, pos mgl32.Vec3) *WorldObject {
	o := WorldObject{
		state:        state,
		model:        model.NewModel(modelFile, hasAnimation),
		hasAnimation: hasAnimation,
		shader:       s,
		modelMat:     mgl32.Translate3D(pos.Elem()),
		normalShader: shader.NewEmbeddedShaderVFG(assets.NormViewVert, assets.NormViewFrag, assets.NormViewGeom),
	}

	if o.hasAnimation {
		animation := model.NewAnimation(o.model)
		o.animator = model.NewAnimator(animation)
	}

	return &o
}

func (o *WorldObject) finaliseLoad() {
	for _, m := range o.model.Meshes {
		m.SetupMesh()
		m.BindTextures()
	}
}

func (o *WorldObject) Update(deltaTime float32) {
	if o.hasAnimation {
		if o.animator == nil {
			panic("Animator wasn't set but object has animation. Nil Pointer")
		}
		o.animator.UpdateAnimation(deltaTime)
	}
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

	o.shader.SetBool("hasAnimation", o.hasAnimation)
	if o.hasAnimation {
		transforms := o.animator.GetFinalBoneMatrices()
		for i, mat := range transforms {
			o.shader.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
		}
	}
	o.model.Draw(o.shader)
}

func (o WorldObject) DrawWithNormals(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3) {
	o.Draw(proj, view, camPos)

	o.normalShader.Use()
	o.normalShader.SetMatrix4("proj", proj)
	o.normalShader.SetMatrix4("view", view)
	o.normalShader.SetMatrix4("model", o.modelMat)
	o.model.Draw(o.normalShader)
}

type Light interface {
	GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3)
}

type PointLight struct {
	Pos mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32
}

func (p PointLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return p.Ambient, p.Diffuse, p.Specular
}

type DirLight struct {
	Dir mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

func (d DirLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return d.Ambient, d.Diffuse, d.Specular
}

type SpotLight struct {
	Pos mgl32.Vec3
	Dir mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32

	cutoff      float32
	outerCutoff float32
}

func (s SpotLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return s.Ambient, s.Diffuse, s.Specular
}
