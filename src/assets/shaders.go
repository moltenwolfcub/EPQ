package assets

func LoadShader(file string) (string, error) {
	embeddedShader, err := shaders.ReadFile("shaders/" + file)
	if err != nil {
		return "", err
	}
	return string(embeddedShader), nil
}

func MustLoadShader(file string) string {
	shader, err := LoadShader(file)
	if err != nil {
		panic("Failed to load shader: " + err.Error())
	}
	return shader
}

var (
	OrangeVert string = MustLoadShader("orange.vert")
	OrangeFrag string = MustLoadShader("orange.frag")

	BlueVert string = MustLoadShader("blue.vert")
	BlueFrag string = MustLoadShader("blue.frag")

	AssimpModelVert string = MustLoadShader("assimpModel.vert")
	AssimpModelFrag string = MustLoadShader("assimpModel.frag")

	AnimatedModelVert string = MustLoadShader("animatedModel.vert")

	SimpleVert string = MustLoadShader("simple.vert")
	SimpleFrag string = MustLoadShader("simple.frag")

	LightingVert string = MustLoadShader("lighting.vert")
	LightingFrag string = MustLoadShader("lighting.frag")
)
