package model

// #cgo pkg-config: assimp zlib
// #include <assimp/cimport.h>
// #include <assimp/scene.h>
// #include <assimp/postprocess.h>
// #include <stdlib.h>
// #include <model.h>
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

//export getRawModel
func getRawModel(path *C.char, size *C.int) *C.char {
	goPath := C.GoString(path)
	data := assets.MustLoadModel(goPath)

	*size = C.int(len(data))
	return (*C.char)(C.CBytes(data))
}

func Mat4assimp2mgl(mat C.struct_aiMatrix4x4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(mat.a1), float32(mat.b1), float32(mat.c1), float32(mat.d1),
		float32(mat.a2), float32(mat.b2), float32(mat.c2), float32(mat.d2),
		float32(mat.a3), float32(mat.b3), float32(mat.c3), float32(mat.d3),
		float32(mat.a4), float32(mat.b4), float32(mat.c4), float32(mat.d4),
	}
}

func Lerp(x mgl32.Vec3, y mgl32.Vec3, a float32) mgl32.Vec3 {
	return x.Mul(1 - a).Add(y.Mul(a))
}

type Model struct {
	Meshes           []Mesh
	textureDirectory string
	texturesLoaded   []Texture

	boneInfoMap map[string]BoneInfo
	boneCounter int
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

// very anti-go style so prolly gonna fix later.
// following a c++ tutorial again
func (m Model) getBoneInfoMap() map[string]BoneInfo {
	return m.boneInfoMap
}
func (m Model) getBoneCount() int {
	return m.boneCounter
}

func (m *Model) loadModel(path string) {
	fileIO := C.CreateMemoryFileIO()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	scene := C.aiImportFileEx(
		cpath,
		C.uint(C.aiProcess_Triangulate|C.aiProcess_FlipUVs),
		fileIO,
	)
	defer C.aiReleaseImport(scene)

	if scene == nil || (scene.mFlags&C.AI_SCENE_FLAGS_INCOMPLETE) != 0 || scene.mRootNode == nil {
		fmt.Println("ERROR::ASSIMP::" + C.GoString(C.aiGetErrorString()))
		return
	}

	m.textureDirectory = strings.Split(path, ".")[0]

	m.boneInfoMap = make(map[string]BoneInfo)

	m.processNode(scene.mRootNode, scene)
}

func (m *Model) processNode(node *C.struct_aiNode, scene *C.struct_aiScene) {
	nodeMeshes := unsafe.Slice(node.mMeshes, node.mNumMeshes)
	sceneMeshes := unsafe.Slice(scene.mMeshes, scene.mNumMeshes)

	for _, nodeMesh := range nodeMeshes {
		mesh := sceneMeshes[nodeMesh]
		processedMesh := m.processMesh(mesh, scene)

		m.Meshes = append(m.Meshes, processedMesh)
	}

	nodeChildren := unsafe.Slice(node.mChildren, node.mNumChildren)

	for _, child := range nodeChildren {
		m.processNode(child, scene)
	}
}

func (m *Model) processMesh(mesh *C.struct_aiMesh, scene *C.struct_aiScene) Mesh {
	var vertices []Vertex
	var indices []uint32
	var textures []Texture

	// vertices
	meshVertices := unsafe.Slice(mesh.mVertices, mesh.mNumVertices)
	meshNormals := unsafe.Slice(mesh.mNormals, mesh.mNumVertices)
	for i := range int(mesh.mNumVertices) {
		var vertex Vertex

		for i := range max_bone_influence {
			vertex.BoneIDs[i] = -1
			vertex.Weights[i] = 0
		}

		vertex.Position = mgl32.Vec3{
			float32(meshVertices[i].x),
			float32(meshVertices[i].y),
			float32(meshVertices[i].z),
		}

		vertex.Normal = mgl32.Vec3{
			float32(meshNormals[i].x),
			float32(meshNormals[i].y),
			float32(meshNormals[i].z),
		}

		if mesh.mTextureCoords[0] != nil {
			meshTextureCoords := unsafe.Slice(mesh.mTextureCoords[0], mesh.mNumVertices)

			vertex.TexCoords = mgl32.Vec2{
				float32(meshTextureCoords[i].x),
				float32(meshTextureCoords[i].y),
			}

		} else {
			vertex.TexCoords = mgl32.Vec2{0, 0}
		}

		vertices = append(vertices, vertex)
	}
	//indices
	meshFaces := unsafe.Slice(mesh.mFaces, mesh.mNumFaces)
	for _, face := range meshFaces {
		faceIndices := unsafe.Slice(face.mIndices, face.mNumIndices)
		for _, faceIndex := range faceIndices {
			indices = append(indices, uint32(faceIndex))
		}
	}

	//material
	sceneMaterials := unsafe.Slice(scene.mMaterials, scene.mNumMaterials)
	material := sceneMaterials[mesh.mMaterialIndex]

	diffuseMaps := m.loadMaterialTextures(material, C.aiTextureType_DIFFUSE, "texture_diffuse")
	textures = append(textures, diffuseMaps...)

	specularMaps := m.loadMaterialTextures(material, C.aiTextureType_SPECULAR, "texture_specular")
	textures = append(textures, specularMaps...)

	// bone data
	meshBones := unsafe.Slice(mesh.mBones, mesh.mNumBones)

	for _, bone := range meshBones {
		boneId := -1
		boneName := C.GoString(&bone.mName.data[0])

		info, ok := m.boneInfoMap[boneName]
		if !ok {
			newInfo := BoneInfo{
				Id:     m.boneCounter,
				Offset: Mat4assimp2mgl(bone.mOffsetMatrix),
			}
			m.boneInfoMap[boneName] = newInfo
			boneId = m.boneCounter
			m.boneCounter++

		} else {
			boneId = info.Id
		}

		meshWeights := unsafe.Slice(bone.mWeights, bone.mNumWeights)

		for _, weight := range meshWeights {
			vertex := vertices[weight.mVertexId]
			for i := range max_bone_influence {
				if vertex.BoneIDs[i] < 0 {
					vertex.Weights[i] = float32(weight.mVertexId)
					vertex.BoneIDs[i] = boneId
					break
				}
			}
		}
	}

	return NewMesh(vertices, indices, textures)
}

func (m *Model) loadMaterialTextures(mat *C.struct_aiMaterial, texture_type C.enum_aiTextureType, typeName string) []Texture {
	var textures []Texture
	for i := range C.aiGetMaterialTextureCount(mat, texture_type) {
		var cstr C.struct_aiString
		C.aiGetMaterialTexture(mat, texture_type, i, &cstr, nil, nil, nil, nil, nil, nil)
		path := C.GoString(&cstr.data[0])

		skip := false
		for _, tex := range m.texturesLoaded {
			if tex.Path == path {
				textures = append(textures, tex)
				skip = true
				break
			}
		}
		if !skip {
			var texture Texture
			texture.Id = TextureFromFile(path, m.textureDirectory)
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

const max_bone_influence int = 4

type BoneInfo struct {
	Id     int        // index in finalBoneMatricies
	Offset mgl32.Mat4 // transforms vertex from model to bone space
}

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2

	BoneIDs [max_bone_influence]int
	Weights [max_bone_influence]float32
}

type Texture struct {
	Id          uint32
	TextureType string
	Path        string
}

type Mesh struct {
	Vertices      []Vertex
	Indices       []uint32
	Textures      []Texture
	vao, vbo, ebo uint32
}

func NewMesh(verts []Vertex, indices []uint32, textures []Texture) Mesh {
	m := Mesh{
		Vertices: verts,
		Indices:  indices,
		Textures: textures,
	}
	m.setupMesh()

	return m
}

func (m Mesh) Draw(shader gogl.Shader) {
	diffuseNr := 1
	specularNr := 1

	for i, tex := range m.Textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		var number string
		name := tex.TextureType
		if name == "texture_diffuse" {
			number = fmt.Sprintf("%d", diffuseNr)
			diffuseNr++
		} else if name == "texture_specular" {
			number = fmt.Sprintf("%d", specularNr)
			diffuseNr++
		}

		shader.SetInt("material."+name+number, int32(i))
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

	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*int(unsafe.Sizeof(Vertex{})), gl.Ptr(m.Vertices), gl.STATIC_DRAW)

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
	// bone IDs
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribIPointer(3, 4, gl.INT, int32(unsafe.Sizeof(Vertex{})), gl.Ptr(unsafe.Offsetof(Vertex{}.BoneIDs)))
	//TODO: change boneIDs to int[4] rather than an ivec4 because it is slightly missleading
	// weights
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 4, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.Ptr(unsafe.Offsetof(Vertex{}.Weights)))

	gl.BindVertexArray(0)
}

type KeyPosition struct {
	pos       mgl32.Vec3
	timeStamp float32
}
type KeyRotation struct {
	rot       mgl32.Quat
	timeStamp float32
}
type KeyScale struct {
	scale     mgl32.Vec3
	timeStamp float32
}

type Bone struct {
	positions      []KeyPosition
	rotations      []KeyRotation
	scales         []KeyScale
	numPositions   int
	numRotations   int
	numScales      int
	localTransform mgl32.Mat4
	name           string
	id             int
}

func NewBone(name string, id int, channel *C.struct_aiNodeAnim) Bone {
	b := Bone{
		name:           name,
		id:             id,
		localTransform: mgl32.Ident4(),
	}

	b.numPositions = int(channel.mNumPositionKeys)
	positionArray := unsafe.Slice(channel.mPositionKeys, channel.mNumPositionKeys)
	for _, pos := range positionArray {
		data := KeyPosition{
			pos: mgl32.Vec3{
				float32(pos.mValue.x),
				float32(pos.mValue.y),
				float32(pos.mValue.z),
			},
			timeStamp: float32(pos.mTime),
		}
		b.positions = append(b.positions, data)
	}

	b.numRotations = int(channel.mNumRotationKeys)
	rotationArray := unsafe.Slice(channel.mRotationKeys, channel.mNumRotationKeys)
	for _, rot := range rotationArray {
		data := KeyRotation{
			rot: mgl32.Quat{
				W: float32(rot.mValue.w),
				V: mgl32.Vec3{
					float32(rot.mValue.x),
					float32(rot.mValue.y),
					float32(rot.mValue.z),
				},
			},
			timeStamp: float32(rot.mTime),
		}
		b.rotations = append(b.rotations, data)
	}

	b.numScales = int(channel.mNumScalingKeys)
	scaleArray := unsafe.Slice(channel.mScalingKeys, channel.mNumScalingKeys)
	for _, scales := range scaleArray {
		data := KeyScale{
			scale: mgl32.Vec3{
				float32(scales.mValue.x),
				float32(scales.mValue.y),
				float32(scales.mValue.z),
			},
			timeStamp: float32(scales.mTime),
		}
		b.scales = append(b.scales, data)
	}

	return b
}

func (b *Bone) Update(animationTime float32) {
	translation := b.interpolatePosition(animationTime)
	rotation := b.interpolateRotation(animationTime)
	scale := b.interpolateScaling(animationTime)
	b.localTransform = translation.Mul4(rotation.Mul4(scale))
}

func (b Bone) GetPositionIndex(animationTime float32) int {
	for index := range b.numPositions - 1 {
		if animationTime < b.positions[index+1].timeStamp {
			return index
		}
	}
	panic(fmt.Errorf("no position index found for animationTime %f", animationTime))
}
func (b Bone) GetRotationIndex(animationTime float32) int {
	for index := range b.numRotations - 1 {
		if animationTime < b.rotations[index+1].timeStamp {
			return index
		}
	}
	panic(fmt.Errorf("no rotation index found for animationTime %f", animationTime))
}
func (b Bone) GetScaleIndex(animationTime float32) int {
	for index := range b.numScales - 1 {
		if animationTime < b.scales[index+1].timeStamp {
			return index
		}
	}
	panic(fmt.Errorf("no scale index found for animationTime %f", animationTime))
}

func (b Bone) getScaleFactor(lastTimeStamp float32, nextTimeStamp float32, animationTime float32) float32 {
	midwayLength := animationTime - lastTimeStamp
	framesDiff := nextTimeStamp - lastTimeStamp
	return midwayLength / framesDiff
}

func (b *Bone) interpolatePosition(animationTime float32) mgl32.Mat4 {
	if b.numPositions == 1 {
		return mgl32.Translate3D(b.positions[0].pos.Elem())
	}
	p0Index := b.GetPositionIndex(animationTime)
	p1Index := p0Index + 1
	scaleFactor := b.getScaleFactor(b.positions[p0Index].timeStamp, b.positions[p1Index].timeStamp, animationTime)

	finalPos := Lerp(b.positions[p0Index].pos, b.positions[p1Index].pos, scaleFactor)
	return mgl32.Translate3D(finalPos.Elem())
}
func (b *Bone) interpolateRotation(animationTime float32) mgl32.Mat4 {
	if b.numRotations == 1 {
		return b.rotations[0].rot.Normalize().Mat4()
	}
	p0Index := b.GetRotationIndex(animationTime)
	p1Index := p0Index + 1
	scaleFactor := b.getScaleFactor(b.rotations[p0Index].timeStamp, b.rotations[p1Index].timeStamp, animationTime)

	finalRotation := mgl32.QuatSlerp(b.rotations[p0Index].rot, b.rotations[p1Index].rot, scaleFactor)
	return finalRotation.Normalize().Mat4()
}
func (b *Bone) interpolateScaling(animationTime float32) mgl32.Mat4 {
	if b.numScales == 1 {
		return mgl32.Scale3D(b.scales[0].scale.Elem())
	}
	p0Index := b.GetScaleIndex(animationTime)
	p1Index := p0Index + 1
	scaleFactor := b.getScaleFactor(b.scales[p0Index].timeStamp, b.scales[p1Index].timeStamp, animationTime)

	finalScale := Lerp(b.scales[p0Index].scale, b.scales[p1Index].scale, scaleFactor)
	return mgl32.Scale3D(finalScale.Elem())
}
