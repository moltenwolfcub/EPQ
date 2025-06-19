#version 460 core

layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aNormal;
layout(location = 2) in vec2 aTexCoord;
layout(location = 3) in int aBoneOffset;
layout(location = 4) in int aBoneCount;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;

const int MAX_BONES = 100;
uniform mat4 finalBonesMatrices[MAX_BONES];

layout(std430, binding = 0) buffer BoneIDBuffer {
int boneIDs[];
};
layout(std430, binding = 1) buffer BoneWeightBuffer {
float boneWeights[];
};

out vec2 TexCoord;

void main() {
	vec4 riggedPos = vec4(0);

	for(int i = 0; i < aBoneCount; i++) {
		int boneIndex = boneIDs[aBoneOffset + i];
		float weight = boneWeights[aBoneOffset + i];
		if(boneIndex < 0)
			continue;
		if(boneIndex >= MAX_BONES) {
			riggedPos = vec4(aPos, 1);
			break;
		}

		vec4 localPos = finalBonesMatrices[boneIndex] * vec4(aPos, 1);
		riggedPos += localPos * weight;
	}

	gl_Position = proj * view * model * riggedPos;
	TexCoord = aTexCoord;
}
