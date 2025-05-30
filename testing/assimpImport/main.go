package main

// #cgo pkg-config: assimp zlib
// #include <assimp/cimport.h>
// #include <assimp/scene.h>
// #include <assimp/postprocess.h>
import "C"
import "fmt"

func main() {

	scene := C.aiImportFile(C.CString("multiMeshNO.glb"), C.uint(C.aiProcess_CalcTangentSpace|
		C.aiProcess_Triangulate|
		C.aiProcess_JoinIdenticalVertices|
		C.aiProcess_SortByPType))

	if scene == nil {
		fmt.Println("ERROR::ASSIMP::" + C.GoString(C.aiGetErrorString()))
	}

	// fmt.Println(scene.mNumMeshes)
	// foo(scene.mRootNode)
}

func foo(bar *C.struct_aiNode) {

}
