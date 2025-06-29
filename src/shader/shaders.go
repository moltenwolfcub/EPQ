package shader

import (
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ProgramID uint32
type ShaderID uint32

func CreateProgram(vertPath string, fragPath string, geomPath string) ProgramID {
	var vert, frag, geom ShaderID
	if vertPath != "" {
		vert = LoadShader(vertPath, gl.VERTEX_SHADER)
	}
	if fragPath != "" {
		frag = LoadShader(fragPath, gl.FRAGMENT_SHADER)
	}
	if geomPath != "" {
		geom = LoadShader(geomPath, gl.GEOMETRY_SHADER)
	}

	shaderProgram := gl.CreateProgram()

	if vertPath != "" {
		gl.AttachShader(shaderProgram, uint32(vert))
	}
	if fragPath != "" {
		gl.AttachShader(shaderProgram, uint32(frag))
	}
	if geomPath != "" {
		gl.AttachShader(shaderProgram, uint32(geom))
	}

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
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram)
}

func LoadShader(path string, shaderType uint32) ShaderID {
	shaderFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	shaderId := CreateShader(string(shaderFile), shaderType)
	return shaderId
}

func CreateProgramFromShaders(vertShader string, fragShader string, geomShader string) ProgramID {
	var vert, frag, geom ShaderID
	if vertShader != "" {
		vert = CreateShader(vertShader, gl.VERTEX_SHADER)
	}
	if fragShader != "" {
		frag = CreateShader(fragShader, gl.FRAGMENT_SHADER)
	}
	if geomShader != "" {
		geom = CreateShader(geomShader, gl.GEOMETRY_SHADER)
	}

	shaderProgram := gl.CreateProgram()

	if vertShader != "" {
		gl.AttachShader(shaderProgram, uint32(vert))
	}
	if fragShader != "" {
		gl.AttachShader(shaderProgram, uint32(frag))
	}
	if geomShader != "" {
		gl.AttachShader(shaderProgram, uint32(geom))
	}

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
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram)
}

func CreateShader(shaderSource string, shaderType uint32) ShaderID {
	shaderId := gl.CreateShader(shaderType)
	shaderSource += "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		panic("Failed to compile shader:\n" + log)
	}
	return ShaderID(shaderId)
}

type Shader interface {
	CheckShadersForChanges()
	Use()

	SetBool(name string, value bool)
	SetInt(name string, value int32)
	SetFloat(name string, value float32)
	SetVec3(name string, value mgl32.Vec3)
	SetMatrix4(name string, value mgl32.Mat4)
}

func UseProgram(id ProgramID) {
	gl.UseProgram(uint32(id))
}
