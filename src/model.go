package main

// #cgo pkg-config: assimp zlib
// #include <assimp/cimport.h>
// #include <assimp/scene.h>
// #include <assimp/postprocess.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"image"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/EPQ/src/assets"
	"github.com/moltenwolfcub/gogl-utils"
)

type Model struct {
	Meshes         []Mesh
	Directory      string
	texturesLoaded []Texture
}

func NewModel(path string) Model {
	m := Model{}

	m.loadModel(path)

	return m
}

func (m Model) Draw(shader gogl.Shader) {
	for _, mesh := range m.Meshes {
		mesh.Draw(shader)
	}
}

func (m *Model) loadModel(path string) {
	cpath := C.CString("assets/models/" + path)
	defer C.free(unsafe.Pointer(cpath))

	scene := C.aiImportFile(
		cpath,
		C.uint(C.aiProcess_Triangulate|C.aiProcess_FlipUVs),
	)
	defer C.aiReleaseImport(scene)

	if scene == nil || (scene.mFlags&C.AI_SCENE_FLAGS_INCOMPLETE) != 0 || scene.mRootNode == nil {
		fmt.Println("ERROR::ASSIMP::" + C.GoString(C.aiGetErrorString()))
		return
	}

	dirIndex := strings.LastIndex(path, "/")
	if dirIndex != -1 {
		m.Directory = path[:dirIndex]
	} else {
		m.Directory = path
	}

	m.processNode(scene.mRootNode, scene)
}

func (m *Model) processNode(node *C.struct_aiNode, scene *C.struct_aiScene) {
	nodeMeshes := unsafe.Slice((*C.uint)(unsafe.Pointer(node.mMeshes)), node.mNumMeshes)
	sceneMeshes := unsafe.Slice((**C.struct_aiMesh)(unsafe.Pointer(scene.mMeshes)), scene.mNumMeshes)
	// meshIndex := *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(node.mMeshes)) + uintptr(i)*unsafe.Sizeof(*node.mMeshes)))
	// mesh := (*C.struct_aiMesh)(unsafe.Pointer(uintptr(unsafe.Pointer(scene.mMeshes)) + uintptr(meshIndex)*unsafe.Sizeof(*scene.mMeshes)))

	for i := range int(node.mNumMeshes) {
		mesh := sceneMeshes[nodeMeshes[i]]
		processedMesh := m.processMesh(mesh, scene)

		m.Meshes = append(m.Meshes, processedMesh)
	}

	nodeChildren := unsafe.Slice((**C.struct_aiNode)(unsafe.Pointer(node.mChildren)), node.mNumChildren)

	for i := range int(node.mNumChildren) {
		m.processNode(nodeChildren[i], scene)
	}
}

func (m *Model) processMesh(mesh *C.struct_aiMesh, scene *C.struct_aiScene) Mesh {
	var verticies []Vertex
	var indices []uint32
	var textures []Texture

	// verticies
	meshVerticies := unsafe.Slice((*C.struct_aiVector3D)(unsafe.Pointer(mesh.mVertices)), mesh.mNumVertices)
	meshNormals := unsafe.Slice((*C.struct_aiVector3D)(unsafe.Pointer(mesh.mNormals)), mesh.mNumVertices)
	for i := range int(mesh.mNumVertices) {
		var vertex Vertex

		vertex.Position = mgl32.Vec3{
			float32(meshVerticies[i].x),
			float32(meshVerticies[i].y),
			float32(meshVerticies[i].z),
		}

		vertex.Normal = mgl32.Vec3{
			float32(meshNormals[i].x),
			float32(meshNormals[i].y),
			float32(meshNormals[i].z),
		}

		if mesh.mTextureCoords[0] != nil {
			meshTextureCoords := unsafe.Slice((*C.struct_aiVector3D)(unsafe.Pointer(mesh.mTextureCoords[0])), mesh.mNumVertices)

			vertex.TexCoords = mgl32.Vec2{
				float32(meshTextureCoords[i].x),
				float32(meshTextureCoords[i].y),
			}

		} else {
			vertex.TexCoords = mgl32.Vec2{0, 0}
		}

		verticies = append(verticies, vertex)
	}
	//indicies
	meshFaces := unsafe.Slice((*C.struct_aiFace)(unsafe.Pointer(mesh.mFaces)), mesh.mNumFaces)
	for i := range int(mesh.mNumFaces) {
		face := meshFaces[i]

		faceIndicies := unsafe.Slice((*C.uint)(unsafe.Pointer(face.mIndices)), face.mNumIndices)
		for j := range int(face.mNumIndices) {
			indices = append(indices, uint32(faceIndicies[j]))
		}
	}

	//material
	sceneMaterials := unsafe.Slice((**C.struct_aiMaterial)(unsafe.Pointer(scene.mMaterials)), scene.mNumMaterials)
	material := sceneMaterials[mesh.mMaterialIndex]

	diffuseMaps := m.loadMaterialTextures(material, C.aiTextureType_DIFFUSE, "texture_diffuse")
	textures = append(textures, diffuseMaps...)

	specularMaps := m.loadMaterialTextures(material, C.aiTextureType_SPECULAR, "texture_specular")
	textures = append(textures, specularMaps...)

	return NewMesh(verticies, indices, textures)
}

func (m *Model) loadMaterialTextures(mat *C.struct_aiMaterial, texture_type C.enum_aiTextureType, typeName string) []Texture {
	var textures []Texture
	for i := range C.aiGetMaterialTextureCount(mat, texture_type) {
		var cstr C.struct_aiString
		C.aiGetMaterialTexture(mat, texture_type, i, &cstr, nil, nil, nil, nil, nil, nil)
		path := C.GoString(&cstr.data[0])

		skip := false
		for j := range len(m.texturesLoaded) {
			if m.texturesLoaded[j].Path == path {
				textures = append(textures, m.texturesLoaded[j])
				skip = true
				break
			}
		}
		if !skip {
			var texture Texture
			texture.Id = TextureFromFile(path, m.Directory)
			texture.TextureType = typeName
			texture.Path = path

			textures = append(textures, texture)
			m.texturesLoaded = append(m.texturesLoaded, texture)
		}
	}
	return textures
}

func TextureFromFile(path, directory string) uint32 {
	loc := fmt.Sprintf("%s/%s", directory, path)
	img := assets.MustLoadImage(loc)

	convertedImg := flipVertical(img)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(convertedImg.Bounds().Dx()),
		int32(convertedImg.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(convertedImg.Pix),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return textureID
}

func flipVertical(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	flipped := image.NewNRGBA(bounds)
	width := bounds.Dx()
	height := bounds.Dy()

	for y := range height {
		for x := range width {
			flipped.Set(x, height-1-y, img.At(x, y))
		}
	}
	return flipped
}

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2
}

type Texture struct {
	Id          uint32
	TextureType string
	Path        string
}

type Mesh struct {
	Verticies     []Vertex
	Indices       []uint32
	Textures      []Texture
	vao, vbo, ebo uint32
}

func NewMesh(verts []Vertex, indices []uint32, textures []Texture) Mesh {
	m := Mesh{
		Verticies: verts,
		Indices:   indices,
		Textures:  textures,
	}
	m.setupMesh()

	return m
}

func (m Mesh) Draw(shader gogl.Shader) {
	diffuseNr := 1
	specularNr := 1

	for i := uint32(0); i < uint32(len(m.Textures)); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + i)
		var number string
		name := m.Textures[i].TextureType
		if name == "texture_diffuse" {
			number = fmt.Sprintf("%d", diffuseNr)
			diffuseNr++
		} else if name == "texture_specular" {
			number = fmt.Sprintf("%d", specularNr)
			diffuseNr++
		}

		shader.SetInt(("material." + name + number), int32(i))
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[i].Id)
	}
	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
	gl.BindVertexArray(0)
}

func (m *Mesh) setupMesh() {
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(m.Verticies)*int(unsafe.Sizeof(Vertex{})), gl.Ptr(m.Verticies), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(m.Indices), gl.STATIC_DRAW)

	// pos
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.Ptr(uintptr(0)))
	// normal
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.Ptr(unsafe.Offsetof(Vertex{}.Normal)))
	// texture Coords
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.Ptr(unsafe.Offsetof(Vertex{}.TexCoords)))

	gl.BindVertexArray(0)
}
