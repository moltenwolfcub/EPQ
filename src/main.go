package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/moltenwolfcub/gogl-utils"
)

func main() {
	r := NewRenderer()
	defer r.Close()

	shaderProgram := gogl.Shader(gogl.NewEmbeddedShader(assets.TriangleVert, assets.TriangleFrag))

	pent := gogl.Pentahedron(1)

	gl.BindVertexArray(0)

	for {
		if handleEvents() != 0 {
			return
		}

		r.Draw(shaderProgram, 0, pent)
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
