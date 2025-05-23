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

	keyboardState []uint8

	x, y, z float32

	projMat mgl32.Mat4
	viewMat mgl32.Mat4
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.setupWindow()

	r.x = -10
	r.y = -10
	r.z = -10

	aspect := float32(WINDOW_WIDTH) / float32(WINDOW_HEIGHT)
	r.projMat = mgl32.Ortho(-aspect*ORTHO_SCALE/2, aspect*ORTHO_SCALE/2, -ORTHO_SCALE/2, ORTHO_SCALE/2, 0.1, 100)
	// r.projMat = mgl32.Perspective(mgl32.DegToRad(45), aspectRatio, 0.1, 100) // perspective version

	r.viewMat = mgl32.HomogRotate3DX(mgl32.DegToRad(30)).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-45)))

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
	r.keyboardState = sdl.GetKeyboardState()
}

func (r *Renderer) Draw(shader gogl.Shader, vao gogl.BufferID, pent gogl.Object) {
	gl.ClearColor(0.0, 0.2, 0.3, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shader.Use()

	// VERY HACKY BASIC MOVEMENT FOR TESTING
	// MUST BE REMOVED FOR THE FINAL IMPLEMENTATION
	if r.keyboardState[sdl.SCANCODE_W] != 0 {
		r.z += .1
	}
	if r.keyboardState[sdl.SCANCODE_S] != 0 {
		r.z -= .1
	}
	if r.keyboardState[sdl.SCANCODE_A] != 0 {
		r.x += .1
	}
	if r.keyboardState[sdl.SCANCODE_D] != 0 {
		r.x -= .1
	}
	if r.keyboardState[sdl.SCANCODE_LSHIFT] != 0 {
		r.y += .1
	}
	if r.keyboardState[sdl.SCANCODE_SPACE] != 0 {
		r.y -= .1
	}

	translatedViewMat := r.viewMat.Mul4(mgl32.Translate3D(r.x, r.y, r.z))

	shader.SetMatrix4("proj", r.projMat)
	shader.SetMatrix4("view", translatedViewMat)

	modelMat := mgl32.Translate3D(0, 0, 0)
	pent.Draw(shader, modelMat)

	r.window.GLSwap()
	shader.CheckShadersForChanges()
}

func (r *Renderer) Close() {
	r.destructor()
}
