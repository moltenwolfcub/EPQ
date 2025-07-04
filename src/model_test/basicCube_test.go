package model_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/go-cmp/cmp"
	"github.com/moltenwolfcub/EPQ/src/model"
)

func TestBCMeshCount(t *testing.T) {
	m := model.NewModel("testdata/basicCube.obj", false)

	if len(m.Meshes) != 1 {
		t.Errorf("NewModel(basicCube): len(meshes). Got %d meshes. Want 1.", len(m.Meshes))
	}

	t.Errorf("FORCED BROKEN FOR TESTING")
}

func TestBCMeshVertexCount(t *testing.T) {
	m := model.NewModel("testdata/basicCube.obj", false)
	mesh := m.Meshes[0]

	if len(mesh.Vertices) != 6*4 {
		t.Errorf("NewModel(basicCube): len(mesh.vertices). Got %d vertices. Want %d.", len(mesh.Vertices), 6*4)
	}
}

func TestBCMeshVertexPositions(t *testing.T) {
	m := model.NewModel("testdata/basicCube.obj", false)
	mesh := m.Meshes[0]

	positions := make([]mgl32.Vec3, 0, len(mesh.Vertices))
	for _, vertex := range mesh.Vertices {
		positions = append(positions, vertex.Position)
	}

	want := []mgl32.Vec3{
		{1, 1, -1},
		{-1, 1, -1},
		{-1, 1, 1},
		{1, 1, 1},

		{1, -1, 1},
		{1, 1, 1},
		{-1, 1, 1},
		{-1, -1, 1},

		{-1, -1, 1},
		{-1, 1, 1},
		{-1, 1, -1},
		{-1, -1, -1},

		{-1, -1, -1},
		{1, -1, -1},
		{1, -1, 1},
		{-1, -1, 1},

		{1, -1, -1},
		{1, 1, -1},
		{1, 1, 1},
		{1, -1, 1},

		{-1, -1, -1},
		{-1, 1, -1},
		{1, 1, -1},
		{1, -1, -1},
	}

	if !cmp.Equal(positions, want) {
		t.Errorf("NewModel(basicCube): mesh.vertices.positions. Got %v. Want %v.", positions, want)
	}
}
