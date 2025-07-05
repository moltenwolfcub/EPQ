package model_test

import (
	"testing"

	"github.com/moltenwolfcub/EPQ/src/model"
)

func TestTextures(t *testing.T) {
	m := model.NewModel("testdata/textureTest.obj", false)
	mesh := m.Meshes[0]

	if len(mesh.Textures) != 3 {
		t.Errorf("len(mesh.Textures), got: %v, want: %v", len(mesh.Textures), 3)
	}

	for _, texture := range mesh.Textures {
		switch texture.TextureType {
		case "texture_diffuse":
			got := texture.Path
			want := "Diffuse.png"
			if got != want {
				t.Errorf("mesh.Texture.path, got: %v, want: %v", got, want)
			}
		case "texture_specular":
			got := texture.Path
			want := "Specular.png"
			if got != want {
				t.Errorf("mesh.Texture.path, got: %v, want: %v", got, want)
			}
		case "texture_roughness":
			got := texture.Path
			want := "Roughness.png"
			if got != want {
				t.Errorf("mesh.Texture.path, got: %v, want: %v", got, want)
			}
		default:
			t.Errorf("mesh.Texture.TextureType, got: %v, want: texture_diffuse, texture_specular or texture_roughness", texture.TextureType)
		}
	}
}
