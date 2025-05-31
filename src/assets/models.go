package assets

func LoadModel(file string) ([]byte, error) {

	embeddedShader, err := models.ReadFile("models/" + file)
	if err != nil {
		return nil, err
	}
	return embeddedShader, nil
}

func MustLoadModel(file string) []byte {
	model, err := LoadModel(file)
	if err != nil {
		panic("Failed to load model: " + err.Error())
	}
	return model
}
