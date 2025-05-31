#include <stdio.h>
#include <stdlib.h>
#include <assimp/cfileio.h>
#include <model.h>

typedef struct
{
	char *data;
	int size;
	int offset;
} EmbeddedData;

static size_t MyReadProc(C_STRUCT aiFile *f, char *buffer, size_t size, size_t count)
{
	EmbeddedData *embed = (EmbeddedData *)f->UserData;
	if (!embed || !buffer || size == 0 || count == 0)
		return 0;

	size_t bytes_left = embed->size - embed->offset;
	size_t bytes_requested = size * count;
	size_t bytes_to_read = bytes_requested < bytes_left ? bytes_requested : bytes_left;

	memcpy(buffer, embed->data + embed->offset, bytes_to_read);
	embed->offset += bytes_to_read;

	return bytes_to_read / size;
}
static size_t MyFileSizeProc(C_STRUCT aiFile *f)
{
	EmbeddedData *embed = (EmbeddedData *)f->UserData;
	if (!embed)
		return 0;
	return embed->size;
}
static C_ENUM aiReturn MySeekProc(C_STRUCT aiFile *f, size_t offset, C_ENUM aiOrigin origin)
{
	EmbeddedData *embed = (EmbeddedData *)f->UserData;
	if (!embed)
		return aiReturn_FAILURE;

	int new_offset = 0;
	switch (origin)
	{
	case aiOrigin_SET:
		new_offset = offset;
		break;
	case aiOrigin_CUR:
		new_offset = embed->offset + offset;
		break;
	case aiOrigin_END:
		new_offset = embed->size + offset;
		break;
	default:
		return aiReturn_FAILURE;
	}

	if (new_offset < 0 || new_offset > embed->size)
	{
		return aiReturn_FAILURE;
	}

	embed->offset = new_offset;
	return aiReturn_FAILURE;
}
static void MyFlushProc(C_STRUCT aiFile *) { printf("FLUSH UNIMPLEMENTED\n"); }
static size_t MyWriteProc(C_STRUCT aiFile *, const char *, size_t, size_t) { return printf("WRITE UNIMPLEMENTED\n") * 0; }
static size_t MyTellProc(C_STRUCT aiFile *) { return printf("TELL UNIMPLEMENTED\n") * 0; }

static C_STRUCT aiFile *MyOpenProc(C_STRUCT aiFileIO *io, const char *filename, const char *mode)
{
	int size = 0;
	char *data = GetRawModel((char *)filename, &size);

	EmbeddedData *embed = (EmbeddedData *)malloc(sizeof(EmbeddedData));
	embed->data = data;
	embed->size = size;
	embed->offset = 0;

	C_STRUCT aiFile *file = (C_STRUCT aiFile *)malloc(sizeof(C_STRUCT aiFile));
	file->ReadProc = MyReadProc;
	file->WriteProc = MyWriteProc;
	file->SeekProc = MySeekProc;
	file->TellProc = MyTellProc;
	file->FileSizeProc = MyFileSizeProc;
	file->FlushProc = MyFlushProc;
	file->UserData = embed;
	return file;
}
static void MyCloseProc(C_STRUCT aiFileIO *io, C_STRUCT aiFile *file)
{
	EmbeddedData *embed = (EmbeddedData *)file->UserData;
	if (embed)
	{
		free(embed->data);
		free(embed);
	}

	free(file);
}

C_STRUCT aiFileIO *CreateMemoryFileIO()
{
	C_STRUCT aiFileIO *io = (C_STRUCT aiFileIO *)malloc(sizeof(C_STRUCT aiFileIO));
	io->OpenProc = MyOpenProc;
	io->CloseProc = MyCloseProc;
	io->UserData = NULL;
	return io;
}
