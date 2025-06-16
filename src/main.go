package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/gogl-utils"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	renderer *Renderer

	keyboardState []uint8

	state     WorldState
	playerPos mgl32.Vec3

	vampireAnimator model.Animator
}

func NewGame() *Game {
	g := Game{}

	g.renderer = NewRenderer()
	g.keyboardState = sdl.GetKeyboardState()

	// orangeShader := gogl.Shader(gogl.NewEmbeddedShader(assets.OrangeVert, assets.OrangeFrag))
	// blueShader := gogl.Shader(gogl.NewEmbeddedShader(assets.BlueVert, assets.BlueFrag))

	// terrain := model.NewModel("terrain.obj")
	// assimpModelShader := gogl.Shader(gogl.NewEmbeddedShader(assets.AssimpModelVert, assets.AssimpModelFrag))

	// simpleAnim := model.NewModel("simpleAnimatedModel.gltf")

	// simpleShader := gogl.Shader(gogl.NewEmbeddedShader(assets.SimpleVert, assets.SimpleFrag))
	animatedShader := gogl.Shader(gogl.NewEmbeddedShader(assets.AnimatedModelVert, assets.AssimpModelFrag))

	vampire := model.NewModel("dancing_vampire.dae")
	vampireAnimation := model.NewAnimation("dancing_vampire.dae", &vampire)
	g.vampireAnimator = model.NewAnimator(&vampireAnimation)

	vampireObject := NewWorldObject(&vampire, animatedShader, mgl32.Vec3{0, 0, 0})
	// vampireObject.modelMat = vampireObject.modelMat.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05))
	vampireObject.uniformSetter = func(s gogl.Shader) gogl.Shader {
		transforms := g.vampireAnimator.GetFinalBoneMatrices()
		for i, mat := range transforms {
			s.SetMatrix4(fmt.Sprintf("finalBonesMatrices[%d]", i), mat)
		}

		return s
	}
	// cube := model.NewCubeModel(1)
	// bigCuge := model.NewCubeModel(2)

	g.state = WorldState{
		// NewWorldObject(&terrain, assimpModelShader, mgl32.Vec3{0, 0, 0}),
		// NewWorldObject(&cube, orangeShader, mgl32.Vec3{0, 0, 0}),
		// NewWorldObject(&bigCuge, blueShader, mgl32.Vec3{5, 0, 0}),
		// NewWorldObject(&cube, blueShader, mgl32.Vec3{0, 3, 0}),
		// NewWorldObject(&bigCuge, orangeShader, mgl32.Vec3{0, 0, -6}),
		// NewWorldObject(simpleAnim, simpleShader, mgl32.Vec3{0, 10, 0}),
		vampireObject,
	}

	g.playerPos = mgl32.Vec3{}

	return &g
}

func (g *Game) close() {
	g.renderer.Close()
}

func (g *Game) runGame() {
	var lastFrame uint64
	for {
		now := sdl.GetTicks64()
		deltaTime := float32(now-lastFrame) / 1000
		lastFrame = now

		if g.handleEvents() != 0 {
			return
		}

		g.vampireAnimator.UpdateAnimation(deltaTime)

		translationVec := mgl32.Vec3{
			float32(g.keyboardState[sdl.SCANCODE_A]) - float32(g.keyboardState[sdl.SCANCODE_D]),
			float32(g.keyboardState[sdl.SCANCODE_LSHIFT]) - float32(g.keyboardState[sdl.SCANCODE_SPACE]),
			float32(g.keyboardState[sdl.SCANCODE_W]) - float32(g.keyboardState[sdl.SCANCODE_S]),
		}
		g.playerPos = g.playerPos.Add(translationVec.Mul(MOVEMENT_SPEED))

		g.renderer.Draw(g.playerPos, g.state)
	}
}

func (g *Game) handleEvents() int {
	for rawEvent := sdl.PollEvent(); rawEvent != nil; rawEvent = sdl.PollEvent() {
		switch event := rawEvent.(type) {
		case *sdl.QuitEvent:
			return 1
		case *sdl.WindowEvent:
			if event.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
				g.renderer.Resize(event.Data1, event.Data2)
			}
		}
	}
	return 0
}

func main() {
	g := NewGame()
	defer g.close()
	g.runGame()
}
