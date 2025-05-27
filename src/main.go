package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/moltenwolfcub/gogl-utils"
)

func main() {
	r := NewRenderer()
	defer r.Close()

	keyboardState := sdl.GetKeyboardState()

	orangeShader := gogl.Shader(gogl.NewEmbeddedShader(assets.TriangleVert, assets.TriangleFrag))
	blueShader := gogl.Shader(gogl.NewEmbeddedShader(assets.Shader2Vert, assets.Shader2Frag))

	state := WorldState{
		NewWorldObject(gogl.Pentahedron(1), orangeShader, mgl32.Vec3{0, 0, 0}),
		NewWorldObject(gogl.Cube(2), blueShader, mgl32.Vec3{5, 0, 0}),
	}

	gl.BindVertexArray(0)

	playerPos := mgl32.Vec3{}

	for {
		if handleEvents() != 0 {
			return
		}

		translationVec := mgl32.Vec3{
			float32(keyboardState[sdl.SCANCODE_A]) - float32(keyboardState[sdl.SCANCODE_D]),
			float32(keyboardState[sdl.SCANCODE_LSHIFT]) - float32(keyboardState[sdl.SCANCODE_SPACE]),
			float32(keyboardState[sdl.SCANCODE_W]) - float32(keyboardState[sdl.SCANCODE_S]),
		}
		playerPos = playerPos.Add(translationVec.Mul(MOVEMENT_SPEED))

		r.Draw(playerPos, state)
	}
}

func handleEvents() int {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return 1
		}
	}
	return 0
}
