#version 460 core
out vec4 FragColor;

in float height;

void main() {
	FragColor = vec4(0.91, 0.35, 0.03, 1.0) - vec4(height);
}
