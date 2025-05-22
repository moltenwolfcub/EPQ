package main

import (
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	r := NewRenderer()
	defer r.Close()

	vertShader := gl.CreateShader(gl.VERTEX_SHADER)
	csource, free := gl.Strs(assets.TriangleVert)
	gl.ShaderSource(vertShader, 1, csource, nil)
	free()
	gl.CompileShader(vertShader)
	var status int32
	gl.GetShaderiv(vertShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertShader, logLength, nil, gl.Str(log))
		panic("Failed to compile shader:\n" + log)
	}

	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource, free = gl.Strs(assets.TriangleFrag)
	gl.ShaderSource(fragShader, 1, csource, nil)
	free()
	gl.CompileShader(fragShader)
	gl.GetShaderiv(fragShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragShader, logLength, nil, gl.Str(log))
		panic("Failed to compile shader:\n" + log)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic("Failed to link program:\n" + log)
	}
	gl.DeleteShader(uint32(vertShader))
	gl.DeleteShader(uint32(fragShader))

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

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
