package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

func (r *Renderer) Draw(playerPos mgl32.Vec3, world WorldState) {
	gl.ClearColor(0.0, 0.2, 0.3, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.alignCamera(playerPos)
	proj, view := r.camera.GetMatricies()

	for _, obj := range world {
		obj.Shader.CheckShadersForChanges()
		obj.Shader.Use()

		obj.Shader.SetMatrix4("proj", proj)
		obj.Shader.SetMatrix4("view", view)

		modelMat := mgl32.Translate3D(obj.Pos.Elem())

		obj.RenderObj.Draw(obj.Shader, modelMat)
	}

	r.window.GLSwap()
}

func (r *Renderer) Close() {
	r.destructor()
}
