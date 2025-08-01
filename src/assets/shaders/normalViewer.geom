#version 330 core

const float MAGNITUDE = 0.4;

layout (triangles) in;
layout (line_strip, max_vertices = 6) out;

in VS_OUT {
	vec3 normal;
} gs_in[];

uniform mat4 proj;

void GenerateLine(int index) {
	gl_Position = proj * gl_in[index].gl_Position;
	EmitVertex();
	gl_Position = proj * (gl_in[index].gl_Position + 
		vec4(gs_in[index].normal, 0.0) * MAGNITUDE);
	EmitVertex();
	EndPrimitive();
}

void main() {
	GenerateLine(0); // first vertex normal
	GenerateLine(1); // second vertex normal
	GenerateLine(2); // third vertex normal
} 
