#version 460 core
out vec4 FragColor;

in float gradient;

void main() {
	FragColor = vec4(0.28, 0.64, 0.93, 1.0) - vec4(gradient);
}
