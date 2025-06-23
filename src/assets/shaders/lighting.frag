#version 460 core

const int POINT = 0;
const int DIRECTION = 1;
const int SPOT = 2;

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
	vec3 pos;
	vec3 direction;

	int lightType;

	vec3 ambient;
	vec3 diffuse;
	vec3 specular;

	float constant;
	float linear;
	float quadratic;

	float cutoff;
	float outerCutoff;
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

	vec3 reflectDir = reflect(-lightDir, normal);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), 256);

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

vec3 CalcDirLight(Light light, vec3 normal, vec3 viewDir, vec3 diffuseColor, vec3 specularColor, float shininess) {
	vec3 lightDir = normalize(-light.direction);

	float diff = max(dot(normal, lightDir), 0.0);

	vec3 reflectDir = reflect(-lightDir, normal);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);

	vec3 ambient = light.ambient * diffuseColor;
	vec3 diffuse = light.diffuse * diffuseColor * diff;
	vec3 specular = light.specular * specularColor * spec;

	return ambient + diffuse + specular;
}

vec3 CalcSpotLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir, vec3 diffuseColor, vec3 specularColor, float shininess) {
	vec3 lightDir = normalize(light.pos - fragPos);

	float theta = dot(lightDir, normalize(-light.direction));
	float epsilon = light.cutoff - light.outerCutoff;
	float intensity = smoothstep(0.0, 1.0, (theta - light.outerCutoff) / epsilon);

	float diff = max(dot(normal, lightDir), 0.0);

	vec3 reflectDir = reflect(-lightDir, normal);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);

	float dist = length(light.pos - fragPos);
	float attenuation = 1/(light.constant + light.linear*dist + light.quadratic*dist*dist);

	vec3 ambient = light.ambient * diffuseColor;
	vec3 diffuse = light.diffuse * diffuseColor * diff;
	vec3 specular = light.specular * specularColor * spec;
	ambient *= attenuation;
	diffuse *= attenuation;
	specular *= attenuation;

	diffuse *= intensity;
	specular *= intensity;

	return ambient + diffuse + specular;
}

void main() {
	float gamma = 2.2;

	vec3 diffuseColor = material.diffuse;
	if(material.hasTexDiffuse) {
		diffuseColor = texture(material.texture_diffuse1, texCoord).rgb;
		diffuseColor = pow(diffuseColor, vec3(gamma));
	}

	vec3 specularColor = material.specular;
	if(material.hasTexSpecular) {
		specularColor = texture(material.texture_specular1, texCoord).rgb;
		specularColor = pow(specularColor, vec3(gamma));
	}

	float shine = material.shininess;
	if(material.hasTexRoughness) {
		float roughness = texture(material.texture_roughness1, texCoord).r;
		roughness = clamp(roughness, 0.0, 0.99);
		shine = 1-roughness;
	}
	float shininessValue = mix(2, 512, shine);

	vec3 norm = normalize(normal);
	vec3 viewDir = normalize(camera - fragPos);

	vec3 result = vec3(0);
	for (int i = 0; i < lights.length(); i++) {
		Light l = lights[i];
		if (l.lightType == POINT) {
			result += CalcPointLight(l, norm, fragPos, viewDir, diffuseColor, specularColor, shininessValue);
		} else if (l.lightType == DIRECTION) {
			result += CalcDirLight(l, norm, viewDir, diffuseColor, specularColor, shininessValue);
		} else if (l.lightType == SPOT) {
			result += CalcSpotLight(l, norm, fragPos, viewDir, diffuseColor, specularColor, shininessValue);
		} else {
			continue;
		}
	}

	FragColor = vec4(result, 1.0);

	FragColor.rgb = pow(FragColor.rgb, vec3(1.0/gamma));
}
