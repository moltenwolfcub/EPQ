package assets

func LoadShader(file string) (string, error) {

	embeddedShader, err := shaders.ReadFile("shaders/" + file)
	if err != nil {
		return "", err
	}
	return string(embeddedShader), nil
}

func MustLoadShader(file string) string {
	img, err := LoadShader(file)
	if err != nil {
		panic("Failed to load shader: " + err.Error())
	}
	return img
}

var (
	TriangleVert string = MustLoadShader("triangle.vert")
	TriangleFrag string = MustLoadShader("triangle.frag")

	Shader2Vert string = MustLoadShader("shader2.vert")
	Shader2Frag string = MustLoadShader("shader2.frag")
)
