#version 460 core

const int POINT = 0;
const int DIRECTION = 1;

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
struct Light {
	// int lightType;

	vec3 pos;
	// vec3 direction;

	vec3 ambient;
	vec3 diffuse;
	vec3 specular;

	float constant;
	float linear;
	float quadratic;
};

in vec3 normal;
in vec3 fragPos;
in vec2 texCoord;

uniform vec3 camera;
uniform Material material;

layout(std430, binding = 2) buffer LightBuffer {
	Light lights[];
};

out vec4 FragColor;

vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir, vec3 diffuseColor, vec3 specularColor, float shininess) {
	vec3 lightDir = normalize(light.pos - fragPos);

	float diff = max(dot(normal, lightDir), 0.0);

	float spec = 0;
	if (shininess != 0) {
		vec3 reflectDir = reflect(-lightDir, normal);
		spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);
	}

	float dist = length(light.pos - fragPos);
	float attenuation = 1/(light.constant + light.linear*dist + light.quadratic*dist*dist);

	vec3 ambient = light.ambient * diffuseColor;
	vec3 diffuse = light.diffuse * diffuseColor * diff;
	vec3 specular = light.specular * specularColor * spec;
	ambient *= attenuation;
	diffuse *= attenuation;
	specular *= attenuation;

	return ambient + diffuse + specular;
}

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
	vec3 viewDir = normalize(camera - fragPos);

	vec3 result = vec3(0);
	for (int i = 0; i < lights.length(); i++) {
		result += CalcPointLight(lights[i], norm, fragPos, viewDir, diffuseColor, specularColor, shininessValue);
	}

	FragColor = vec4(result, 1.0);

	float gamma = 2.2;
	FragColor.rgb = pow(FragColor.rgb, vec3(1.0/gamma));
}
