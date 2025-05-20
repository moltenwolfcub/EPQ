#include <assimp/cimport.h>		// Plain-C interface
#include <assimp/scene.h>		// Output data structure
#include <assimp/postprocess.h> // Post processing flags
#include <iostream>
using namespace std;

int main()
{
	const struct aiScene *scene =
		aiImportFile("multiMesh.glb",
					 aiProcess_CalcTangentSpace |
						 aiProcess_Triangulate |
						 aiProcess_JoinIdenticalVertices |
						 aiProcess_SortByPType);

	cout << (scene->mNumMeshes);

	aiReleaseImport(scene);
	return 0;
}
