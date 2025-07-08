package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type Player struct {
	state *WorldState

	pos mgl32.Vec3

	model      *model.Model
	animations map[string]*model.Animation
	animator   *model.Animator

	shader shader.Shader
}

func NewPlayer(state *WorldState, generalShader shader.Shader) *Player {
	p := Player{
		state:  state,
		shader: generalShader,
	}

	p.model = model.NewModel("player.glb", true)

	p.animations = model.LoadAllAnimations(p.model)
	p.animator = model.NewAnimator(p.animations["idle"])

	return &p
}

func (p *Player) finaliseLoad() {
	for _, m := range p.model.Meshes {
		m.SetupMesh()
		m.BindTextures()
	}
}

func (p *Player) Update(deltaTime float32) {
	if p.animator == nil {
		panic("Player animator wasn't set. Nil Pointer")
	}
	p.animator.UpdateAnimation(deltaTime)
}

func (p Player) Draw(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3) {
	p.shader.CheckShadersForChanges()
	p.shader.Use()

	p.shader.SetMatrix4("proj", proj)
	p.shader.SetMatrix4("view", view)
	p.shader.SetMatrix4("model", mgl32.Translate3D(p.pos.Elem()))
	p.shader.SetVec3("camera", mgl32.Vec3{}.Sub(camPos))

	if p.state.lightingSSBO != 0 {
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, p.state.lightingSSBO)
	} else {
		panic("Lighting was never bound on state")
	}

	p.shader.SetBool("hasAnimation", true)
	transforms := p.animator.GetFinalBoneMatrices()
	for i, mat := range transforms {
		p.shader.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
	}
	p.model.Draw(p.shader)
}
