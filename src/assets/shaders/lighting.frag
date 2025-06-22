#version 460 core

in vec3 normal;
in vec3 fragPos;
in vec2 texCoord;

uniform vec3 camera;

struct Material {
	vec3 diffuse;
	vec3 specular;
	float shininess;

	sampler2D texture_diffuse1;
	bool hasTexDiffuse;
	sampler2D texture_specular1;
	bool hasTexSpecular;
	sampler2D texture_roughness1;
	bool hasTexRoughness;
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
	vec3 diffuseColor = material.diffuse;
	if(material.hasTexDiffuse) {
		diffuseColor = texture(material.texture_diffuse1, texCoord).rgb;
	}

	vec3 specularColor = material.specular;
	if(material.hasTexSpecular) {
		specularColor = texture(material.texture_specular1, texCoord).rgb;
	}

	float shininessValue = material.shininess;
	if(material.hasTexRoughness) {
		float roughness = texture(material.texture_roughness1, texCoord).r;
		shininessValue = mix(0, 1000, 1-roughness);
	}

	vec3 norm = normalize(normal);
	vec3 lightDir = normalize(light.pos - fragPos);

	vec3 ambient = light.ambient * diffuseColor;

	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = light.diffuse * diff * diffuseColor;

	vec3 specular = vec3(0);
	if (shininessValue != 0) {
		vec3 viewDir = normalize(camera - fragPos);
		vec3 reflectDir = reflect(-lightDir, norm);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininessValue);
		specular = light.specular * spec * specularColor;
	}

	vec3 result = ambient + diffuse + specular;

	FragColor = vec4(result, 1.0);

	float gamma = 2.2;
	FragColor.rgb = pow(FragColor.rgb, vec3(1.0/gamma));
}
