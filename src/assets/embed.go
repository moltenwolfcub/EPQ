package assets

import "embed"

var (
	//go:embed shaders
	shaders embed.FS

	//go:embed textures
	textures embed.FS
)
