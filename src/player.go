package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
)

type Player struct {
	state *WorldState

	pos      mgl32.Vec3
	rotation mgl32.Quat

	model            *model.Model
	animations       map[string]*model.Animation
	animator         *model.Animator
	currentAnimation string

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

func (p *Player) Move(velocity mgl32.Vec3) {
	p.pos = p.pos.Add(velocity)

	if velocity.Len() == 0 {
		if p.currentAnimation != "idle" {
			p.animator.PlayAnimation(p.animations["idle"])
			p.currentAnimation = "idle"
		}
	} else {

		forwards := mgl32.Vec2{0, 1}
		lookingDir := mgl32.Vec2{-velocity.X(), velocity.Z()}
		if lookingDir.Len() != 0 { //to catch only vertical movement
			theta := angleBetween(forwards, lookingDir)

			p.rotation = mgl32.QuatRotate(theta, mgl32.Vec3{0, 1, 0})
		}

		if p.currentAnimation != "run" {
			p.animator.PlayAnimation(p.animations["run"])
			p.currentAnimation = "run"
		}
	}
}

func angleBetween(a, b mgl32.Vec2) float32 {
	normal := mgl32.Vec3{0, 0, 1}

	a3 := a.Vec3(0)
	b3 := b.Vec3(0)

	cross := a3.Cross(b3)
	dot := a3.Dot(b3)

	normalCross := cross.Dot(normal)

	theta := float32(math.Atan2(float64(normalCross), float64(dot)))

	return theta
}

func (p Player) Draw(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3) {
	p.shader.CheckShadersForChanges()
	p.shader.Use()

	p.shader.SetMatrix4("proj", proj)
	p.shader.SetMatrix4("view", view)

	model := mgl32.Translate3D(p.pos.Elem()).Mul4(p.rotation.Normalize().Mat4())
	p.shader.SetMatrix4("model", model)

	p.shader.SetVec3("camera", camPos)

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
