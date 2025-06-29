package shader

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderWithPaths struct {
	id          ProgramID
	vertPath    string
	vertModTime time.Time
	fragPath    string
	fragModTime time.Time
}

func NewShaderFromFilePaths(vertPath string, fragPath string) *ShaderWithPaths {
	id := CreateProgram(vertPath, fragPath, "")

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
		id := CreateProgram(s.vertPath, s.fragPath, "")

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

func getFileModTime(path string) time.Time {
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}
