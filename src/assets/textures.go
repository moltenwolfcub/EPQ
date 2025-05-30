package assets

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
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

var ()
