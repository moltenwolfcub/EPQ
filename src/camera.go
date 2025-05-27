package main

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	Pos mgl32.Vec3

	projMat mgl32.Mat4
	viewMat mgl32.Mat4
}

func NewCamera() *Camera {
	c := &Camera{}
	c.preCalculateMatricies()

	return c
}

func (c *Camera) preCalculateMatricies() {
	aspectRatio := float32(WINDOW_WIDTH) / float32(WINDOW_HEIGHT)

	c.projMat = mgl32.Ortho(-aspectRatio*ORTHO_SCALE/2, aspectRatio*ORTHO_SCALE/2, -ORTHO_SCALE/2, ORTHO_SCALE/2, 0.1, 100)
	c.viewMat = mgl32.HomogRotate3DX(mgl32.DegToRad(30)).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-45)))
}

// gets the projection and view matricies
func (c *Camera) GetMatricies() (proj mgl32.Mat4, view mgl32.Mat4) {
	translatedViewMat := c.viewMat.Mul4(mgl32.Translate3D(c.Pos.Elem()))

	return c.projMat, translatedViewMat
}
