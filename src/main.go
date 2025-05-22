package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	r := NewRenderer()
	defer r.Close()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		r.Draw()
	}
}
