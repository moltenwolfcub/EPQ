#version 460 core

//TODO make these more dynamic and passed in from code
const vec3 lightColor = vec3(1, 1, 1);
const float ambientStrength = 0.25;
const vec3 objectColor = vec3(0.1, 0.3, 1);
const vec3 lightPos = vec3(3, 4, -2);

in vec3 normal;
in vec3 fragPos;

// struct Material {
// 	sampler2D texture_diffuse1;
// };
// uniform Material material;

out vec4 FragColor;

void main() {
	vec3 ambient = ambientStrength * lightColor;

	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(lightPos - fragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * lightColor;

	vec3 result = (ambient + diffuse) * objectColor;

	FragColor = vec4(result, 1.0);
}
