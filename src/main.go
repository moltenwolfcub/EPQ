package main

import (
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
}

func NewGame() *Game {
	g := Game{}

	g.renderer = NewRenderer()
	g.keyboardState = sdl.GetKeyboardState()

	orangeShader := gogl.Shader(gogl.NewEmbeddedShader(assets.OrangeVert, assets.OrangeFrag))
	blueShader := gogl.Shader(gogl.NewEmbeddedShader(assets.BlueVert, assets.BlueFrag))

	terrain := model.NewModel("terrain.obj")
	assimpModelShader := gogl.Shader(gogl.NewEmbeddedShader(assets.AssimpModelVert, assets.AssimpModelFrag))

	g.state = WorldState{
		NewWorldObject(terrain, assimpModelShader, mgl32.Vec3{0, 0, 0}),
		NewWorldObject(model.NewCubeModel(1), orangeShader, mgl32.Vec3{0, 0, 0}),
		NewWorldObject(model.NewCubeModel(2), blueShader, mgl32.Vec3{5, 0, 0}),
		NewWorldObject(model.NewCubeModel(1), blueShader, mgl32.Vec3{0, 3, 0}),
		NewWorldObject(model.NewCubeModel(2), orangeShader, mgl32.Vec3{0, 0, -6}),
	}

	g.playerPos = mgl32.Vec3{}

	return &g
}

func (g *Game) close() {
	g.renderer.Close()
}

func (g *Game) runGame() {
	for {
		if g.handleEvents() != 0 {
			return
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
