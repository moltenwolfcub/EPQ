package main

// #cgo pkg-config: assimp
// #cgo LDFLAGS: -static -lstdc++ -lz -lm -ldl
// #include <assimp/cimport.h>
// #include <assimp/scene.h>
// #include <assimp/postprocess.h>
import "C"
import "fmt"

func main() {

	// scene := C.aiImportFile(C.CString("test.obj"), C.uint(C.aiProcess_Triangulate|C.aiProcess_FlipUVs))

	// fmt.Println(C.aiImportFile(nil, C.uint(0)))
	fmt.Println(uint(C.aiProcess_Triangulate))
}
