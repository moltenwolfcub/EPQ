#version 460 core

in vec3 normal;
in vec3 fragPos;

uniform vec3 camera;

struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float shininess;
};
uniform Material material;

struct Light {
	vec3 pos;

	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
};
uniform Light light;

out vec4 FragColor;

void main() {
	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(light.pos - fragPos);

	vec3 ambient = light.ambient * material.ambient;

	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = light.diffuse * diff * material.diffuse;

	vec3 viewDir = normalize(camera - fragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
	vec3 specular = light.specular * spec * material.specular;

	vec3 result = ambient + diffuse + specular;

	FragColor = vec4(result, 1.0);
}
