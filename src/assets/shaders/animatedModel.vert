#version 460 core

layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aNormal;
layout(location = 2) in vec2 aTexCoord;
layout(location = 3) in ivec4 aBoneIds;
layout(location = 4) in vec4 aweights;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;

const int MAX_BONES = 100;
const int MAX_BONES_INFLUENCE = 4;
uniform mat4 finalBonesMatrices[MAX_BONES];

out vec2 TexCoord;

void main() {
	vec4 riggedPos = vec4(0);
	// vec4 riggedPos = vec4(aPos, 1);

	for(int i = 0; i < MAX_BONES_INFLUENCE; i++) {
		if(aBoneIds[i] == -1) {
			continue;
		}
		if(aBoneIds[i] >= MAX_BONES) {
			riggedPos = vec4(aPos, 1);
			break;
		}

		vec4 localPos = finalBonesMatrices[aBoneIds[i]] * vec4(aPos, 1);
		riggedPos += localPos * aweights[i];
		// vec3 localNormal = mat3(finalBonesMatrices[aBoneIds[i]]) * aNormal;
	}

	gl_Position = proj * view * model * riggedPos;
	TexCoord = aTexCoord;
}
