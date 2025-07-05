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

func TestBCMeshVertexNormals(t *testing.T) {
	m := model.NewModel("testdata/basicCube.obj", false)
	mesh := m.Meshes[0]

	normals := make([]mgl32.Vec3, 0, len(mesh.Vertices))
	for _, vertex := range mesh.Vertices {
		normals = append(normals, vertex.Normal)
	}

	want := []mgl32.Vec3{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},

		{0, 0, 1},
		{0, 0, 1},
		{0, 0, 1},
		{0, 0, 1},

		{-1, 0, 0},
		{-1, 0, 0},
		{-1, 0, 0},
		{-1, 0, 0},

		{0, -1, 0},
		{0, -1, 0},
		{0, -1, 0},
		{0, -1, 0},

		{1, 0, 0},
		{1, 0, 0},
		{1, 0, 0},
		{1, 0, 0},

		{0, 0, -1},
		{0, 0, -1},
		{0, 0, -1},
		{0, 0, -1},
	}

	if !cmp.Equal(normals, want) {
		t.Errorf("NewModel(basicCube): mesh.vertices.normals. Got %v. Want %v.", normals, want)
	}
}

func TestBCMeshVertexTexCoords(t *testing.T) {
	m := model.NewModel("testdata/basicCube.obj", false)
	mesh := m.Meshes[0]

	texCoords := make([]mgl32.Vec2, 0, len(mesh.Vertices))
	for _, vertex := range mesh.Vertices {
		texCoords = append(texCoords, vertex.TexCoords)
	}

	// 1-Ycoord because of the flipUVs flag
	want := []mgl32.Vec2{
		{0.625, 1 - 0.50},
		{0.875, 1 - 0.50},
		{0.875, 1 - 0.75},
		{0.625, 1 - 0.75},

		{0.375, 1 - 0.75},
		{0.625, 1 - 0.75},
		{0.625, 1 - 1.00},
		{0.375, 1 - 1.00},

		{0.375, 1 - 0.00},
		{0.625, 1 - 0.00},
		{0.625, 1 - 0.25},
		{0.375, 1 - 0.25},

		{0.125, 1 - 0.50},
		{0.375, 1 - 0.50},
		{0.375, 1 - 0.75},
		{0.125, 1 - 0.75},

		{0.375, 1 - 0.50},
		{0.625, 1 - 0.50},
		{0.625, 1 - 0.75},
		{0.375, 1 - 0.75},

		{0.375, 1 - 0.25},
		{0.625, 1 - 0.25},
		{0.625, 1 - 0.50},
		{0.375, 1 - 0.50},
	}

	if !cmp.Equal(texCoords, want) {
		t.Errorf("NewModel(basicCube): mesh.vertices.normals. Got %v. Want %v.", texCoords, want)
	}
}
