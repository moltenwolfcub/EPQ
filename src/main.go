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

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	gogl.GenBindBuffer(gl.ARRAY_BUFFER) // VBO
	VAO := gogl.GenBindVertexArray()

	gogl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	for {
		if handleEvents() != 0 {
			return
		}

		r.Draw(shaderProgram, VAO)
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
