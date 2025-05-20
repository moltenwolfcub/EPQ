#include <assimp/cimport.h> // Plain-C interface

int main()
{
	const struct aiScene *scene = aiImportFile("", 0);

	aiReleaseImport(scene);
	return 0;
}
