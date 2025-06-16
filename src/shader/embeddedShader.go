package shader

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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
