#version 460 core

const int MAX_BONES = 100;

layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aNormal;
layout(location = 2) in vec2 aTexCoord;
layout(location = 3) in int aBoneOffset;
layout(location = 4) in int aBoneCount;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;
uniform mat4 finalBonesMatrices[MAX_BONES];

layout(std430, binding = 0) buffer BoneIDBuffer {
	int boneIDs[];
};
layout(std430, binding = 1) buffer BoneWeightBuffer {
	float boneWeights[];
};

out vec3 normal;
out vec3 fragPos;

void main() {
	gl_Position = proj * view * model * vec4(aPos, 1.0);
	normal = mat3(transpose(inverse(model))) * aNormal;
	fragPos = vec3(model*vec4(aPos,1.0));
}
