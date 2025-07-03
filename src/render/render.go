package render

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/settings"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	window *sdl.Window

	camera *camera
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.setupWindow()

	r.camera = newCamera()
	r.camera.pos = mgl32.Vec3{-10, -10, -10}

	return r
}

func (r *Renderer) setupWindow() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 6)

	window, err := sdl.CreateWindow(
		settings.WINDOW_TITLE,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		settings.WINDOW_WIDTH, settings.WINDOW_HEIGHT,
		sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
	)

	if err != nil {
		panic(err)
	}
	window.GLCreateContext()

	gl.Init()
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	r.window = window
}

func (r *Renderer) Resize(nexX, newY int32) {
	settings.WINDOW_WIDTH = nexX
	settings.WINDOW_HEIGHT = newY

	gl.Viewport(0, 0, settings.WINDOW_WIDTH, settings.WINDOW_HEIGHT)

	r.camera.preCalculateMatricies()
}

func (r *Renderer) Draw(camPos mgl32.Vec3, world []Renderable) {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.camera.pos = camPos
	proj, view := r.camera.getMatricies()

	for _, obj := range world {
		obj.Draw(proj, view, r.camera.pos)
	}

	r.window.GLSwap()
}

func (r *Renderer) Close() {
	sdl.Quit()
	r.window.Destroy()
}

type Renderable interface {
	Draw(proj mgl32.Mat4, view mgl32.Mat4, camPos mgl32.Vec3)
}
