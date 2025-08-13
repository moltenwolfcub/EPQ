package state

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type WorldObject struct {
	state *WorldState

	model        *model.Model
	animator     *model.Animator
	hasAnimation bool

	shader   shader.Shader
	modelMat mgl32.Mat4
}

func NewWorldObjectFromModel(state *WorldState, m *model.Model, s shader.Shader, pos mgl32.Vec3) *WorldObject {
	return &WorldObject{
		state:        state,
		model:        m,
		hasAnimation: false,
		shader:       s,
		modelMat:     mgl32.Translate3D(pos.Elem()),
	}
}

func NewWorldObject(state *WorldState, modelFile string, hasAnimation bool, s shader.Shader, pos mgl32.Vec3) *WorldObject {
	o := WorldObject{
		state:        state,
		model:        model.NewModel(modelFile, hasAnimation),
		hasAnimation: hasAnimation,
		shader:       s,
		modelMat:     mgl32.Translate3D(pos.Elem()),
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
	if o.state.Objects[0] == o {
		o.modelMat = mgl32.Translate3D(o.state.Player.position.Elem()) //TODO: remove player position debugging
	}
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
	o.shader.SetVec3("camera", camPos)

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

var normalShader shader.Shader

func (o WorldObject) DrawWithNormals(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3) {
	o.Draw(proj, view, camPos)

	if normalShader == nil {
		normalShader = shader.NewEmbeddedShaderVFG(assets.NormViewVert, assets.NormViewFrag, assets.NormViewGeom)
	}

	normalShader.Use()
	normalShader.SetMatrix4("proj", proj)
	normalShader.SetMatrix4("view", view)
	normalShader.SetMatrix4("model", o.modelMat)
	o.model.Draw(normalShader)
}
