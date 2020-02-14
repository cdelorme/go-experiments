#version 410 core

layout(location = 0) in vec4 position;
layout(location = 1) in vec2 texCoord;

out vec2 v_texCoord;

uniform mat4 u_projection;
uniform mat4 u_camera;
uniform mat4 u_model;

void main() {
	v_texCoord = texCoord;
	gl_Position = u_projection * u_camera * u_model * position;
}
