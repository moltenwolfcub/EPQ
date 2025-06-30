package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

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

	normalShader = shader.NewEmbeddedShaderVFG(assets.NormViewVert, assets.NormViewFrag, assets.NormViewGeom)
	orangeShader := shader.NewEmbeddedShaderVF(assets.OrangeVert, assets.OrangeFrag)
	blueShader := shader.NewEmbeddedShaderVF(assets.BlueVert, assets.BlueFrag)
	generalShader := shader.NewEmbeddedShaderVF(assets.GeneralVert, assets.GeneralFrag)

	cube := model.NewCubeModel(0.5)

	g.state = NewWorldState()

	g.state.Lights = []Light{
		// PointLight{
		// 	Pos:                  mgl32.Vec3{0, 0, 0},
		// 	Ambient:              mgl32.Vec3{1, 1, 1},
		// 	Diffuse:              mgl32.Vec3{0, 0, 0},
		// 	Specular:             mgl32.Vec3{0, 0, 0},
		// 	ConstantAttenuation:  1.0,
		// 	LinearAttenuation:    0,
		// 	QuadraticAttenuation: 0,
		// },
		PointLight{
			Pos:                  mgl32.Vec3{-2, 5, -2},
			Ambient:              mgl32.Vec3{0, 0, 0},
			Diffuse:              mgl32.Vec3{0.3, 0.5, 1},
			Specular:             mgl32.Vec3{0.6, 0.8, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.0075,
			QuadraticAttenuation: 0.07,
		},
		PointLight{
			Pos:                  mgl32.Vec3{0, 1, 3},
			Ambient:              mgl32.Vec3{0.1, 0.05, 0},
			Diffuse:              mgl32.Vec3{0.5, 0.25, 0.1},
			Specular:             mgl32.Vec3{1, 0.5, 0.2},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.14,
			QuadraticAttenuation: 0.07,
		},
		SpotLight{
			Pos:                  mgl32.Vec3{6, 6, 0},
			Dir:                  mgl32.Vec3{-1, -1, 0},
			Ambient:              mgl32.Vec3{0, 0, 0},
			Diffuse:              mgl32.Vec3{1, 1, 1},
			Specular:             mgl32.Vec3{1, 1, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.09,
			QuadraticAttenuation: 0.032,
			cutoff:               12.5,
			outerCutoff:          15,
		},
	}
	g.state.BindLights()

	g.state.Objects = append(g.state.Objects,
		NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{-2, 5, -2}),
		NewWorldObjectFromModel(g.state, cube, orangeShader, mgl32.Vec3{0, 1, 3}),
		NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{6, 6, 0}),
		NewWorldObject(g.state, "firePit.obj", false, generalShader, mgl32.Vec3{0, 0, 0}),
		NewWorldObject(g.state, "terrain.obj", false, generalShader, mgl32.Vec3{0, -1, 0}),
		NewWorldObject(g.state, "multiAnimation.glb", true, generalShader, mgl32.Vec3{0, 1, 0}),
	)
	g.state.Player = NewPlayer(g.state, generalShader)
	g.state.FinaliseLoad()

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

		g.state.Player.Update(deltaTime)
		for _, object := range g.state.Objects {
			object.Update(deltaTime)
		}

		translationVec := mgl32.Vec3{
			float32(g.keyboardState[sdl.SCANCODE_LEFT]) - float32(g.keyboardState[sdl.SCANCODE_RIGHT]),
			float32(g.keyboardState[sdl.SCANCODE_KP_0]) - float32(g.keyboardState[sdl.SCANCODE_RCTRL]),
			float32(g.keyboardState[sdl.SCANCODE_UP]) - float32(g.keyboardState[sdl.SCANCODE_DOWN]),
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	g := NewGame()
	defer g.close()
	g.runGame()
}
