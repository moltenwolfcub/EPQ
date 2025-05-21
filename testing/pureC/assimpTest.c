#include <assimp/cimport.h>
#include <assimp/scene.h>
#include <assimp/postprocess.h>
#include <stdio.h>

int main()
{
	const struct aiScene *scene =
		aiImportFile("multiMesh.glb",
					 aiProcess_CalcTangentSpace |
						 aiProcess_Triangulate |
						 aiProcess_JoinIdenticalVertices |
						 aiProcess_SortByPType);

	printf("%d", scene->mNumMeshes);

	aiReleaseImport(scene);
}
