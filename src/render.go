package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	window     *sdl.Window
	destructor func()
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.setupWindow()

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

	r.window = window
	r.destructor = func() {
		sdl.Quit()
		window.Destroy()
	}
}

func (r *Renderer) Draw() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r.window.GLSwap()
}

func (r *Renderer) Close() {
	r.destructor()
}
