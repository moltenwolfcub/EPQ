package model

import (
	"github.com/go-gl/mathgl/mgl32"
)

func NewCubeModel(size float32) Model {
	m := Model{}

	pos := size / 2
	verts := []Vertex{
		{Position: mgl32.Vec3{-pos, -pos, -pos}, Normal: mgl32.Vec3{0, 0, -1}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{+pos, -pos, -pos}, Normal: mgl32.Vec3{0, 0, -1}, TexCoords: mgl32.Vec2{1, 0}},
		{Position: mgl32.Vec3{-pos, +pos, -pos}, Normal: mgl32.Vec3{0, 0, -1}, TexCoords: mgl32.Vec2{0, 1}},
		{Position: mgl32.Vec3{+pos, +pos, -pos}, Normal: mgl32.Vec3{0, 0, -1}, TexCoords: mgl32.Vec2{1, 1}},

		{Position: mgl32.Vec3{-pos, -pos, +pos}, Normal: mgl32.Vec3{0, 0, +1}, TexCoords: mgl32.Vec2{1, 0}},
		{Position: mgl32.Vec3{+pos, -pos, +pos}, Normal: mgl32.Vec3{0, 0, +1}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{-pos, +pos, +pos}, Normal: mgl32.Vec3{0, 0, +1}, TexCoords: mgl32.Vec2{1, 1}},
		{Position: mgl32.Vec3{+pos, +pos, +pos}, Normal: mgl32.Vec3{0, 0, +1}, TexCoords: mgl32.Vec2{0, 1}},

		{Position: mgl32.Vec3{-pos, -pos, -pos}, Normal: mgl32.Vec3{0, -1, 0}, TexCoords: mgl32.Vec2{0, 1}},
		{Position: mgl32.Vec3{+pos, -pos, -pos}, Normal: mgl32.Vec3{0, -1, 0}, TexCoords: mgl32.Vec2{1, 1}},
		{Position: mgl32.Vec3{-pos, -pos, +pos}, Normal: mgl32.Vec3{0, -1, 0}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{+pos, -pos, +pos}, Normal: mgl32.Vec3{0, -1, 0}, TexCoords: mgl32.Vec2{1, 0}},

		{Position: mgl32.Vec3{-pos, +pos, -pos}, Normal: mgl32.Vec3{0, +1, 0}, TexCoords: mgl32.Vec2{0, 1}},
		{Position: mgl32.Vec3{+pos, +pos, -pos}, Normal: mgl32.Vec3{0, +1, 0}, TexCoords: mgl32.Vec2{1, 1}},
		{Position: mgl32.Vec3{-pos, +pos, +pos}, Normal: mgl32.Vec3{0, +1, 0}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{+pos, +pos, +pos}, Normal: mgl32.Vec3{0, +1, 0}, TexCoords: mgl32.Vec2{1, 0}},

		{Position: mgl32.Vec3{-pos, -pos, -pos}, Normal: mgl32.Vec3{-1, 0, 0}, TexCoords: mgl32.Vec2{0, 1}},
		{Position: mgl32.Vec3{-pos, +pos, -pos}, Normal: mgl32.Vec3{-1, 0, 0}, TexCoords: mgl32.Vec2{1, 1}},
		{Position: mgl32.Vec3{-pos, -pos, +pos}, Normal: mgl32.Vec3{-1, 0, 0}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{-pos, +pos, +pos}, Normal: mgl32.Vec3{-1, 0, 0}, TexCoords: mgl32.Vec2{1, 0}},

		{Position: mgl32.Vec3{+pos, -pos, -pos}, Normal: mgl32.Vec3{+1, 0, 0}, TexCoords: mgl32.Vec2{0, 1}},
		{Position: mgl32.Vec3{+pos, +pos, -pos}, Normal: mgl32.Vec3{+1, 0, 0}, TexCoords: mgl32.Vec2{1, 1}},
		{Position: mgl32.Vec3{+pos, -pos, +pos}, Normal: mgl32.Vec3{+1, 0, 0}, TexCoords: mgl32.Vec2{0, 0}},
		{Position: mgl32.Vec3{+pos, +pos, +pos}, Normal: mgl32.Vec3{+1, 0, 0}, TexCoords: mgl32.Vec2{1, 0}},
	}
	indices := []uint32{
		2, 3, 0,
		3, 1, 0,

		4, 7, 6,
		4, 5, 7,

		8, 9, 10,
		9, 11, 10,

		14, 13, 12,
		14, 15, 13,

		18, 17, 16,
		18, 19, 17,

		20, 21, 22,
		21, 23, 22,
	}
	textures := []Texture{}
	material := Material{}

	m.Meshes = []Mesh{NewMesh(verts, indices, textures, material, []int32{-3}, []float32{0})}
	return m
}
