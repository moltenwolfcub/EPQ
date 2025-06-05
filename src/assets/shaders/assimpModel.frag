#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

struct Material {
	sampler2D texture_diffuse1;
};
uniform Material material;

void main() {
	vec4 mappedColor = texture2D(material.texture_diffuse1, TexCoord);
	if(mappedColor.a < 0.1) {
		discard;
	}
	FragColor = mappedColor;
}
