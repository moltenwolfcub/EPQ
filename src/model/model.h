#pragma once

#include <assimp/cfileio.h>

C_STRUCT aiFileIO *CreateMemoryFileIO();

extern char *getRawModel(char *path, int *size);
