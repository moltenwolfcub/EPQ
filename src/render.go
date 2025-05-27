package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/gogl-utils"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	window     *sdl.Window
	destructor func()

	camera *Camera
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.setupWindow()

	r.camera = NewCamera()
	r.camera.Pos = mgl32.Vec3{-10, -10, -10}

	return r
}

func (r *Renderer) setupWindow() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow(
		WINDOW_TITLE,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH, WINDOW_HEIGHT,
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
	r.destructor = func() {
		sdl.Quit()
		window.Destroy()
	}
}

// TODO: needs a lot of changing for when a player is implemented and a worldstate is made
func (r *Renderer) alignCamera(focus mgl32.Vec3) {
	newPos := focus.Add(mgl32.Vec3{-10, -10, -10})
	r.camera.Pos = newPos
}

func (r *Renderer) Draw(playerPos mgl32.Vec3, shader1 gogl.Shader, pent gogl.Object, shader2 gogl.Shader, cube gogl.Object) {
	gl.ClearColor(0.0, 0.2, 0.3, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.alignCamera(playerPos)
	proj, view := r.camera.GetMatricies()

	shader1.Use()

	shader1.SetMatrix4("proj", proj)
	shader1.SetMatrix4("view", view)

	modelMat := mgl32.Translate3D(0, 0, 0)
	pent.Draw(shader1, modelMat)

	shader2.Use()

	shader1.SetMatrix4("proj", proj)
	shader1.SetMatrix4("view", view)

	modelMat = mgl32.Translate3D(5, 0, 0)
	cube.Draw(shader2, modelMat)

	r.window.GLSwap()
	shader1.CheckShadersForChanges()
	shader2.CheckShadersForChanges()
}

func (r *Renderer) Close() {
	r.destructor()
}
