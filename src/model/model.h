#pragma once

#include <assimp/cfileio.h>

C_STRUCT aiFileIO *CreateMemoryFileIO();

extern char *GetRawModel(char *path, int *size);
