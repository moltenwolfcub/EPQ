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
	"github.com/moltenwolfcub/EPQ/src/render"
	"github.com/moltenwolfcub/EPQ/src/settings"
	"github.com/moltenwolfcub/EPQ/src/shader"
	"github.com/moltenwolfcub/EPQ/src/state"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	renderer *render.Renderer

	keyboardState []uint8

	state *state.WorldState

	camPos         mgl32.Vec3
	detachedCamera bool
}

func NewGame() *Game {
	g := Game{}

	g.renderer = render.NewRenderer()
	g.keyboardState = sdl.GetKeyboardState()

	orangeShader := shader.NewEmbeddedShaderVF(assets.OrangeVert, assets.OrangeFrag)
	blueShader := shader.NewEmbeddedShaderVF(assets.BlueVert, assets.BlueFrag)
	generalShader := shader.NewEmbeddedShaderVF(assets.GeneralVert, assets.GeneralFrag)

	cube := model.NewCubeModel(0.5)

	g.state = state.NewWorldState()

	g.state.Lights = []state.Light{
		// PointLight{
		// 	Pos:                  mgl32.Vec3{0, 0, 0},
		// 	Ambient:              mgl32.Vec3{1, 1, 1},
		// 	Diffuse:              mgl32.Vec3{0, 0, 0},
		// 	Specular:             mgl32.Vec3{0, 0, 0},
		// 	ConstantAttenuation:  1.0,
		// 	LinearAttenuation:    0,
		// 	QuadraticAttenuation: 0,
		// },
		state.PointLight{
			Pos:                  mgl32.Vec3{-2, 5, -2},
			Ambient:              mgl32.Vec3{0, 0, 0},
			Diffuse:              mgl32.Vec3{0.3, 0.5, 1},
			Specular:             mgl32.Vec3{0.6, 0.8, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.0075,
			QuadraticAttenuation: 0.07,
		},
		state.PointLight{
			Pos:                  mgl32.Vec3{0, 1, 3},
			Ambient:              mgl32.Vec3{0.1, 0.05, 0},
			Diffuse:              mgl32.Vec3{0.5, 0.25, 0.1},
			Specular:             mgl32.Vec3{1, 0.5, 0.2},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.14,
			QuadraticAttenuation: 0.07,
		},
		state.SpotLight{
			Pos:                  mgl32.Vec3{6, 6, 0},
			Dir:                  mgl32.Vec3{-1, -1, 0},
			Ambient:              mgl32.Vec3{0, 0, 0},
			Diffuse:              mgl32.Vec3{1, 1, 1},
			Specular:             mgl32.Vec3{1, 1, 1},
			ConstantAttenuation:  1.0,
			LinearAttenuation:    0.09,
			QuadraticAttenuation: 0.032,
			Cutoff:               12.5,
			OuterCutoff:          15,
		},
	}
	g.state.BindLights()

	g.state.Objects = append(g.state.Objects,
		state.NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{-2, 5, -2}),
		state.NewWorldObjectFromModel(g.state, cube, orangeShader, mgl32.Vec3{0, 1, 3}),
		state.NewWorldObjectFromModel(g.state, cube, blueShader, mgl32.Vec3{6, 6, 0}),
		state.NewWorldObject(g.state, "firePit.obj", false, generalShader, mgl32.Vec3{0, 0, 0}),
		state.NewWorldObject(g.state, "terrain.obj", false, generalShader, mgl32.Vec3{0, -1, 0}),
		state.NewWorldObject(g.state, "multiAnimation.glb", true, generalShader, mgl32.Vec3{0, 1, 0}),
		state.NewWorldObject(g.state, "axis.obj", false, generalShader, mgl32.Vec3{0, -5, 0}),
	)
	g.state.Player = state.NewPlayer(g.state, generalShader)
	g.state.FinaliseLoad()

	g.detachedCamera = false
	g.alignCamera()

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

		playerAcceleration := []mgl32.Vec3{}
		translationVec := mgl32.Vec3{
			float32(g.keyboardState[sdl.SCANCODE_A]) - float32(g.keyboardState[sdl.SCANCODE_D]),
			float32(g.keyboardState[sdl.SCANCODE_SPACE]) - float32(g.keyboardState[sdl.SCANCODE_LSHIFT]),
			float32(g.keyboardState[sdl.SCANCODE_W]) - float32(g.keyboardState[sdl.SCANCODE_S]),
		}
		if g.detachedCamera {
			deltaPos := translationVec.Mul(settings.PLAYER_ACCELLERATION * deltaTime)
			deltaPos = mgl32.Vec3{
				-deltaPos.X(),
				deltaPos.Y(),
				-deltaPos.Z(),
			}
			g.camPos = g.camPos.Add(deltaPos)
		} else {
			g.alignCamera()

			playerAcceleration = append(playerAcceleration, translationVec.Mul(settings.PLAYER_ACCELLERATION))
		}

		if !g.state.Player.Flying {
			gravity := mgl32.Vec3{0, -settings.GRAVITY, 0}
			playerAcceleration = append(playerAcceleration, gravity)
		}

		g.state.Player.Update(deltaTime, playerAcceleration)
		for _, object := range g.state.Objects {
			object.Update(deltaTime)
		}

		g.renderer.Draw(g.camPos, g.state.ToRender())
	}
}

func (g *Game) alignCamera() {
	g.camPos = g.state.Player.GetPosition().Add(mgl32.Vec3{10, 10, 10})
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
		case *sdl.KeyboardEvent:
			if event.Type == sdl.KEYDOWN {
				switch event.Keysym.Sym {
				case sdl.K_TAB:
					g.detachedCamera = !g.detachedCamera
				case sdl.K_f:
					g.state.Player.Flying = !g.state.Player.Flying
				}
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
