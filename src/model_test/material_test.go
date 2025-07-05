package model_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/model"
)

func TestDiffuse(t *testing.T) {
	var tests = []struct {
		model string

		want mgl32.Vec3
	}{
		{
			"white.obj",
			mgl32.Vec3{1, 1, 1},
		},
		{
			"blue.obj",
			mgl32.Vec3{0, 0, 1},
		},
		{
			"pink.obj",
			mgl32.Vec3{0.75, 0.1, 1},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.model, func(t *testing.T) {
			m := model.NewModel("testdata/"+testCase.model, false)
			mesh := m.Meshes[0]

			got := mesh.Material.Diffuse

			if got != testCase.want {
				t.Errorf("Material.Diffuse, got: %v, expected: %v", got, testCase.want)
			}
		})
	}
}

func TestSpecular(t *testing.T) {
	var tests = []struct {
		model string

		want mgl32.Vec3
	}{
		{
			"fullWhiteSpec.obj",
			mgl32.Vec3{1, 1, 1},
		},
		{
			"lowGreenSpec.obj",
			mgl32.Vec3{0.2, 0.2, 0.2}, //blender doesn't actually export specular tint
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.model, func(t *testing.T) {
			m := model.NewModel("testdata/"+testCase.model, false)
			mesh := m.Meshes[0]

			got := mesh.Material.Specular

			if got != testCase.want {
				t.Errorf("Material.Specular, got: %v, expected: %v", got, testCase.want)
			}
		})
	}
}

func TestShininess(t *testing.T) {
	var tests = []struct {
		model string

		want float32
	}{
		{
			"fullShine.obj",
			1000,
		},
		{
			"fullRough.obj",
			0,
		},
		{
			"partialRough.obj",
			62.5,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.model, func(t *testing.T) {
			m := model.NewModel("testdata/"+testCase.model, false)
			mesh := m.Meshes[0]

			got := mesh.Material.Shininess

			if got != testCase.want {
				t.Errorf("Material.Shininess, got: %v, expected: %v", got, testCase.want)
			}
		})
	}
}
