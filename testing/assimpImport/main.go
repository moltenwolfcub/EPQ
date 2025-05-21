package main

// #cgo pkg-config: assimp zlib
// #include <assimp/cimport.h>
// #include <assimp/scene.h>
// #include <assimp/postprocess.h>
import "C"
import "fmt"

func main() {

	scene := C.aiImportFile(C.CString("multiMesh.glb"), C.uint(C.aiProcess_CalcTangentSpace|
		C.aiProcess_Triangulate|
		C.aiProcess_JoinIdenticalVertices|
		C.aiProcess_SortByPType))

	fmt.Println(scene.mNumMeshes)
}
