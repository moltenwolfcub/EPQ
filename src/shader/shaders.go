package shader

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ProgramID uint32
type ShaderID uint32

func CreateProgram(vertPath string, fragPath string) ProgramID {
	vert := LoadShader(vertPath, gl.VERTEX_SHADER)
	frag := LoadShader(fragPath, gl.FRAGMENT_SHADER)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
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

func CreateProgramFromShaders(vertShader string, fragShader string) ProgramID {
	vert := CreateShader(vertShader, gl.VERTEX_SHADER)
	frag := CreateShader(fragShader, gl.FRAGMENT_SHADER)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
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

type ShaderWithPaths struct {
	id          ProgramID
	vertPath    string
	vertModTime time.Time
	fragPath    string
	fragModTime time.Time
}

func NewShaderFromFilePaths(vertPath string, fragPath string) *ShaderWithPaths {
	id := CreateProgram(vertPath, fragPath)

	s := ShaderWithPaths{
		id:       id,
		vertPath: vertPath,
		fragPath: fragPath,

		vertModTime: getFileModTime(vertPath),
		fragModTime: getFileModTime(fragPath),
	}

	return &s
}

func (s *ShaderWithPaths) Use() {
	UseProgram(s.id)
}

func (s *ShaderWithPaths) CheckShadersForChanges() {
	vertModTime := getFileModTime(s.vertPath)
	fragModTime := getFileModTime(s.fragPath)
	if v, f := !vertModTime.Equal(s.vertModTime), !fragModTime.Equal(s.fragModTime); v || f {
		if v {
			fmt.Printf("A vertex shader file has been modified: %s\n", s.vertPath)
			s.vertModTime = vertModTime
		}
		if f {
			fmt.Printf("A fragment shader file has been modified: %s\n", s.fragPath)
			s.fragModTime = fragModTime
		}
		id := CreateProgram(s.vertPath, s.fragPath)

		gl.DeleteProgram(uint32(s.id))
		s.id = id
	}
}

func (s *ShaderWithPaths) SetBool(name string, value bool) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	if value {
		gl.Uniform1i(loc, 1)
	} else {
		gl.Uniform1i(loc, 0)
	}
}
func (s *ShaderWithPaths) SetInt(name string, value int32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1i(loc, value)
}
func (s *ShaderWithPaths) SetFloat(name string, value float32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1f(loc, value)
}
func (s *ShaderWithPaths) SetMatrix4(name string, value mgl32.Mat4) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	m4 := [16]float32(value)
	gl.UniformMatrix4fv(loc, 1, false, &m4[0])
}
func (s *ShaderWithPaths) SetVec3(name string, value mgl32.Vec3) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	v3 := [3]float32(value)
	gl.Uniform3fv(loc, 1, &v3[0])
}

type EmbeddedShader struct {
	id         ProgramID
	vertShader string
	fragShader string
}

func NewEmbeddedShader(vertShader string, fragShader string) *EmbeddedShader {
	id := CreateProgramFromShaders(vertShader, fragShader)

	s := EmbeddedShader{
		id:         id,
		vertShader: vertShader,
		fragShader: fragShader,
	}

	return &s
}

func (s *EmbeddedShader) Use() {
	UseProgram(s.id)
}
func (s *EmbeddedShader) CheckShadersForChanges() {}

func (s *EmbeddedShader) SetBool(name string, value bool) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	if value {
		gl.Uniform1i(loc, 1)
	} else {
		gl.Uniform1i(loc, 0)
	}
}
func (s *EmbeddedShader) SetInt(name string, value int32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1i(loc, value)
}
func (s *EmbeddedShader) SetFloat(name string, value float32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1f(loc, value)
}
func (s *EmbeddedShader) SetMatrix4(name string, value mgl32.Mat4) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	m4 := [16]float32(value)
	gl.UniformMatrix4fv(loc, 1, false, &m4[0])
}
func (s *EmbeddedShader) SetVec3(name string, value mgl32.Vec3) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	v3 := [3]float32(value)
	gl.Uniform3fv(loc, 1, &v3[0])
}

func getFileModTime(path string) time.Time {
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}

func UseProgram(id ProgramID) {
	gl.UseProgram(uint32(id))
}
