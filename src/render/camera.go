package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/settings"
)

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
	aspectRatio := float32(settings.WINDOW_WIDTH) / float32(settings.WINDOW_HEIGHT)

	if settings.DEBUG_PERSPECTIVE {
		c.projMat = mgl32.Perspective(45, aspectRatio, 0.1, 100)
	} else {
		c.projMat = mgl32.Ortho(-aspectRatio*settings.ORTHO_SCALE/2, aspectRatio*settings.ORTHO_SCALE/2, -settings.ORTHO_SCALE/2, settings.ORTHO_SCALE/2, 0.1, 100)
	}
	c.viewMat = mgl32.HomogRotate3DX(mgl32.DegToRad(30)).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-45)))
}

// gets the projection and view matricies
func (c *Camera) GetMatricies() (proj mgl32.Mat4, view mgl32.Mat4) {
	translatedViewMat := c.viewMat.Mul4(mgl32.Translate3D(c.Pos.Mul(-1).Elem()))

	return c.projMat, translatedViewMat
}
