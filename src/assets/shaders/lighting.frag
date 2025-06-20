#version 460 core

//TODO make these more dynamic and passed in from code
const vec3 lightColor = vec3(1, 1, 1);
const float ambientStrength = 0.25;
const float specularStrength = 0.75;
const vec3 objectColor = vec3(0.1, 0.3, 0.8);
const vec3 lightPos = vec3(-2, 5, -2);
const int shininess = 32;

in vec3 normal;
in vec3 fragPos;

uniform vec3 camera;

// struct Material {
// 	sampler2D texture_diffuse1;
// };
// uniform Material material;

out vec4 FragColor;

void main() {
	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(lightPos - fragPos);

	vec3 ambient = ambientStrength * lightColor;

	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * lightColor;

	vec3 viewDir = normalize(camera - fragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
	vec3 specular = specularStrength * spec * lightColor;

	vec3 result = (ambient + diffuse + specular) * objectColor;

	FragColor = vec4(result, 1.0);
}
