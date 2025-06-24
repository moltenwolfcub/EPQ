package assets

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func LoadImage(file string) (image.Image, error) {
	embeddedImage, err := textures.ReadFile("textures/" + file)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(bytes.NewReader(embeddedImage))
	if err != nil {
		return nil, err
	}
	return image, nil
}

func MustLoadImage(file string) image.Image {
	img, err := LoadImage(file)
	if err != nil {
		panic("Failed to load PNG: " + err.Error())
	}
	return img
}

func LoadSDLImage(file string) (*sdl.Surface, error) {
	embeddedImage, err := textures.ReadFile("textures/" + file)
	if err != nil {
		return nil, err
	}

	rwops, err := sdl.RWFromMem(embeddedImage)
	if err != nil {
		return nil, err
	}
	defer rwops.Close()

	surface, err := img.LoadRW(rwops, false)
	if err != nil {
		return nil, err
	}
	return surface, nil
}

func MustLoadSDLImage(file string) image.Image {
	img, err := LoadSDLImage(file)
	if err != nil {
		panic("Failed to load PNG: " + err.Error())
	}
	return img
}

var ()
