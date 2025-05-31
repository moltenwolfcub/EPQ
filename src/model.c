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
	printf("READ\n");
	printf("\treading: size=%d, count=%d\n", size, count);

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
static size_t MyWriteProc(C_STRUCT aiFile *f, const char *, size_t, size_t)
{
	printf("WRITE\n");
	return 0;
}
static size_t MyTellProc(C_STRUCT aiFile *f)
{
	printf("TELL\n");
	return 0;
}
static size_t MyFileSizeProc(C_STRUCT aiFile *f)
{
	printf("FILE_SIZE\n");

	EmbeddedData *embed = (EmbeddedData *)f->UserData;
	if (!embed)
		return 0;
	printf("\tfile-size: %d\n", embed->size);
	return embed->size;
}
static C_ENUM aiReturn MySeekProc(C_STRUCT aiFile *f, size_t offset, C_ENUM aiOrigin origin)
{
	printf("SEEK\n");
	printf("\tsought: offset=%zu, origin=%d\n", offset, origin);

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
static void MyFlushProc(C_STRUCT aiFile *f)
{
	printf("FLUSH\n");
}

static C_STRUCT aiFile *MyOpenProc(C_STRUCT aiFileIO *io, const char *filename, const char *mode)
{
	printf("OPEN\n");
	printf("\topened: %s in %s\n", filename, mode);

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
	printf("CLOSE\n");

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
	printf("CREATE\n");

	C_STRUCT aiFileIO *io = (C_STRUCT aiFileIO *)malloc(sizeof(C_STRUCT aiFileIO));
	io->OpenProc = MyOpenProc;
	io->CloseProc = MyCloseProc;
	io->UserData = NULL;
	return io;
}

// aiFile callbacks
// typedef size_t          (*aiFileWriteProc) (C_STRUCT aiFile*,   const char*, size_t, size_t);
// typedef size_t          (*aiFileReadProc)  (C_STRUCT aiFile*,   char*, size_t,size_t);
// typedef size_t          (*aiFileTellProc)  (C_STRUCT aiFile*);
// typedef void            (*aiFileFlushProc) (C_STRUCT aiFile*);
// typedef C_ENUM aiReturn (*aiFileSeek)      (C_STRUCT aiFile*, size_t, C_ENUM aiOrigin);

// aiFileIO callbacks
// typedef C_STRUCT aiFile* (*aiFileOpenProc)  (C_STRUCT aiFileIO*, const char*, const char*);
// typedef void             (*aiFileCloseProc) (C_STRUCT aiFileIO*, C_STRUCT aiFile*);

// Represents user-defined data
// typedef char* aiUserData;

// struct aiFileIO
// {
//     /** Function used to open a new file
//      */
//     aiFileOpenProc OpenProc;

//     /** Function used to close an existing file
//      */
//     aiFileCloseProc CloseProc;

//     /** User-defined, opaque data */
//     aiUserData UserData;
// };

// struct aiFile {
//     /** Callback to read from a file */
//     aiFileReadProc ReadProc;

//     /** Callback to write to a file */
//     aiFileWriteProc WriteProc;

//     /** Callback to retrieve the current position of
//      *  the file cursor (ftell())
//      */
//     aiFileTellProc TellProc;

//     /** Callback to retrieve the size of the file,
//      *  in bytes
//      */
//     aiFileTellProc FileSizeProc;

//     /** Callback to set the current position
//      * of the file cursor (fseek())
//      */
//     aiFileSeek SeekProc;

//     /** Callback to flush the file contents
//      */
//     aiFileFlushProc FlushProc;

//     /** User-defined, opaque data
//      */
//     aiUserData UserData;
// };
