package shader

import (
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func CreateProgramWithGeom(vertShader string, fragShader string, geomShader string) ProgramID {
	vert := CreateShader(vertShader, gl.VERTEX_SHADER)
	frag := CreateShader(fragShader, gl.FRAGMENT_SHADER)
	geom := CreateShader(geomShader, gl.GEOMETRY_SHADER)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.AttachShader(shaderProgram, uint32(geom))
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

type EmbeddedShaderWithGeom struct {
	id         ProgramID
	vertShader string
	fragShader string
	geomShader string
}

func NewEmbeddedShaderWithGeom(vertShader string, fragShader string, geomShader string) *EmbeddedShaderWithGeom {
	id := CreateProgramWithGeom(vertShader, fragShader, geomShader)

	s := EmbeddedShaderWithGeom{
		id:         id,
		vertShader: vertShader,
		fragShader: fragShader,
		geomShader: geomShader,
	}

	return &s
}

func (s *EmbeddedShaderWithGeom) Use() {
	UseProgram(s.id)
}
func (s *EmbeddedShaderWithGeom) CheckShadersForChanges() {}

func (s *EmbeddedShaderWithGeom) SetBool(name string, value bool) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	if value {
		gl.Uniform1i(loc, 1)
	} else {
		gl.Uniform1i(loc, 0)
	}
}
func (s *EmbeddedShaderWithGeom) SetInt(name string, value int32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1i(loc, value)
}
func (s *EmbeddedShaderWithGeom) SetFloat(name string, value float32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1f(loc, value)
}
func (s *EmbeddedShaderWithGeom) SetMatrix4(name string, value mgl32.Mat4) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	m4 := [16]float32(value)
	gl.UniformMatrix4fv(loc, 1, false, &m4[0])
}
func (s *EmbeddedShaderWithGeom) SetVec3(name string, value mgl32.Vec3) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	v3 := [3]float32(value)
	gl.Uniform3fv(loc, 1, &v3[0])
}
