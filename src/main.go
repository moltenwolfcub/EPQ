package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/moltenwolfcub/EPQ/src/model"
	"github.com/moltenwolfcub/EPQ/src/shader"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	renderer *Renderer

	keyboardState []uint8

	state     *WorldState
	playerPos mgl32.Vec3
}

func NewGame() *Game {
	g := Game{}

	g.renderer = NewRenderer()
	g.keyboardState = sdl.GetKeyboardState()

	orangeShader := shader.Shader(shader.NewEmbeddedShader(assets.OrangeVert, assets.OrangeFrag))
	blueShader := shader.Shader(shader.NewEmbeddedShader(assets.BlueVert, assets.BlueFrag))
	// assimpModelShader := shader.Shader(shader.NewEmbeddedShader(assets.AssimpModelVert, assets.AssimpModelFrag))
	// animatedShader := shader.Shader(shader.NewEmbeddedShader(assets.AnimatedModelVert, assets.AssimpModelFrag))
	lightingShader := shader.Shader(shader.NewEmbeddedShader(assets.LightingVert, assets.LightingFrag))

	cube := model.NewCubeModel(0.5)
	// bigCuge := model.NewCubeModel(2)

	g.state = NewWorldState()

	g.state.Lights = []Light{
		{
			Pos:                  mgl32.Vec3{-2, 5, -2},
			Ambient:              mgl32.Vec3{0.1, 0.1, 0.1},
			Diffuse:              mgl32.Vec3{0.3, 0.5, 1},
			Specular:             mgl32.Vec3{0.5, 0.7, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.09,
			QuadraticAttenuation: 0.032,
		},
		{
			Pos:                  mgl32.Vec3{0, 0, 3},
			Ambient:              mgl32.Vec3{0.1, 0.05, 0},
			Diffuse:              mgl32.Vec3{1, 0.5, 0.2},
			Specular:             mgl32.Vec3{1, 0.5, 0.2},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.045,
			QuadraticAttenuation: 0.0075,
		},
		{
			Pos:                  mgl32.Vec3{4, 1, 0},
			Ambient:              mgl32.Vec3{0.1, 0.1, 0.1},
			Diffuse:              mgl32.Vec3{1, 1, 1},
			Specular:             mgl32.Vec3{1, 1, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.027,
			QuadraticAttenuation: 0.0028,
		},
	}
	g.state.BindLights()

	g.state.Objects = []*WorldObject{
		// NewWorldObject(g.state, "terrain.obj", false, assimpModelShader, mgl32.Vec3{0, 0, 0}),
		NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{-2, 5, -2}),
		NewWorldObjectFromModel(g.state, cube, orangeShader, mgl32.Vec3{0, 0, 3}),
		NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{4, 1, 0}),
		// NewWorldObjectFromModel(g.state, bigCuge, blueShader, mgl32.Vec3{5, 0, 0}),
		// NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{0, 5, 0}),
		// NewWorldObjectFromModel(g.state, bigCuge, orangeShader, mgl32.Vec3{0, 0, -6}),
		// NewWorldObject(g.state, simpleAnim, simpleShader, mgl32.Vec3{0, 10, 0}),
		// NewWorldObject(g.state, "dancing_vampire.dae", true, animatedShader, mgl32.Vec3{0, 1, 0}),
		NewWorldObject(g.state, "firePit.obj", false, lightingShader, mgl32.Vec3{0, 0, 0}),
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

		for _, object := range g.state.Objects {
			object.Update(deltaTime)
		}

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
