package state

import "github.com/go-gl/mathgl/mgl32"

type Light interface {
	GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3)
}

type PointLight struct {
	Pos mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32
}

func (p PointLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return p.Ambient, p.Diffuse, p.Specular
}

type DirLight struct {
	Dir mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

func (d DirLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return d.Ambient, d.Diffuse, d.Specular
}

type SpotLight struct {
	Pos mgl32.Vec3
	Dir mgl32.Vec3

	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3

	ConstantAttenuation  float32
	LinearAttenuation    float32
	QuadraticAttenuation float32

	Cutoff      float32
	OuterCutoff float32
}

func (s SpotLight) GetLightComponents() (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	return s.Ambient, s.Diffuse, s.Specular
}
