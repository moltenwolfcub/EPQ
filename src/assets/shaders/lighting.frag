#version 460 core

//TODO make these more dynamic and passed in from code
const vec3 ambientLightColor = vec3(1,1,1);
const float ambientStrength = 0.4;
const vec3 objectColor = vec3(0.1,0.3,1);

in vec3 normal;

// struct Material {
// 	sampler2D texture_diffuse1;
// };
// uniform Material material;

out vec4 FragColor;

void main() {
	vec3 ambient = ambientStrength * ambientLightColor;

	vec3 result = ambient * objectColor;

	FragColor = vec4(result, 1.0);
}
