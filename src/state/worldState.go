package state

import (
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/render"
)

type WorldState struct {
	Objects []*WorldObject
	Lights  []Light

	Player *Player

	lightingSSBO uint32
}

func NewWorldState() *WorldState {
	state := &WorldState{
		Objects: make([]*WorldObject, 0),
	}

	gl.GenBuffers(1, &state.lightingSSBO)

	return state
}

func (s *WorldState) FinaliseLoad() {
	s.Player.finaliseLoad()
	for _, o := range s.Objects {
		o.finaliseLoad()
	}
}

func (s *WorldState) BindLights() {
	//padding because vec3s need to be aligned to 16 bytes in SSBOS
	type internalLight struct {
		pos         [3]float32
		_pad1       float32
		dir         [3]float32
		lightType   int32
		ambient     [3]float32
		_pad2       float32
		diffuse     [3]float32
		_pad3       float32
		specular    [3]float32
		constant    float32
		linear      float32
		quadratic   float32
		cutoff      float32
		outerCutoff float32
	}

	internalLights := []internalLight{}
	for _, light := range s.Lights {
		var il internalLight

		switch l := light.(type) {
		case PointLight:
			il = internalLight{
				lightType: 0,
				pos:       l.Pos,
				ambient:   l.Ambient,
				diffuse:   l.Diffuse,
				specular:  l.Specular,
				constant:  l.ConstantAttenuation,
				linear:    l.LinearAttenuation,
				quadratic: l.QuadraticAttenuation,
			}
		case DirLight:
			il = internalLight{
				lightType: 1,
				dir:       l.Dir,
				ambient:   l.Ambient,
				diffuse:   l.Diffuse,
				specular:  l.Specular,
			}
		case SpotLight:
			il = internalLight{
				lightType:   2,
				pos:         l.Pos,
				dir:         l.Dir,
				ambient:     l.Ambient,
				diffuse:     l.Diffuse,
				specular:    l.Specular,
				constant:    l.ConstantAttenuation,
				linear:      l.LinearAttenuation,
				quadratic:   l.QuadraticAttenuation,
				cutoff:      float32(math.Cos(float64(mgl32.DegToRad(l.Cutoff)))),
				outerCutoff: float32(math.Cos(float64(mgl32.DegToRad(l.OuterCutoff)))),
			}
		}

		internalLights = append(internalLights, il)
	}

	lightsPtr := gl.Ptr(new(int)) //dummy pointer if length is 0 because C needs an address even if it's unused
	if len(internalLights) != 0 {
		lightsPtr = gl.Ptr(internalLights)
	}

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, s.lightingSSBO)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(internalLights)*int(unsafe.Sizeof(internalLight{})), lightsPtr, gl.STATIC_DRAW)
}

func (w WorldState) ToRender() []render.Renderable {
	r := make([]render.Renderable, 0, len(w.Objects)+1)

	for _, obj := range w.Objects {
		r = append(r, obj)
	}

	r = append(r, w.Player)

	return r
}
