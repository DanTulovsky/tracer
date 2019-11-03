package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	points = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}
	colors = []float32{
		1, 0, 0, // red
		0, 1, 0, // green
		0, 0, 1, // blue
	}
)

const (
	width  = 500
	height = 500

	// executed for each vertex
	vertexShaderSource = `
#version 410
// “layout(location = 0)” refers to the buffer we 
// use to feed the vertexPosition_modelspace attribute
layout(location = 0) in vec3 vp;
layout(location = 1) in vec3 vertexColor;

out vec3 fragmentColor;

// called for each vertex
void main() {
	// set the vertex position to whatever was in the buffer
	gl_Position = vec4(vp, 1.0);
	gl_PointSize = 5.0;

	fragmentColor = vertexColor;
}
` + "\x00"

	// executed for each sample
	// since we use 4x antialising, we have 4 samples in each pixel
	fragmentShaderSource = `
#version 410
// passed in from the vertexShader above
in vec3 fragmentColor;

out vec4 color;

void main() {
	color = vec4(fragmentColor, 1.0);
}
` + "\x00"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Tracer", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Allows drawing points
	gl.Enable(gl.PROGRAM_POINT_SIZE)
	// Depth test
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.DrawArrays(gl.POINTS, 0, int32(len(points)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVbo gives our data to OpenGL
func makeVbo(points []float32) {
	// points buffer
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// the 0 here is the buffer in the VertexAttribPointer
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		0,        // must match the layout in the shader
		3,        // size
		gl.FLOAT, // type
		false,    // normalized?
		0,        // stride
		nil)      // array buffer offset

	// colors buffer
	var colorbo uint32
	gl.GenBuffers(1, &colorbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)

	// the 1 here is the buffer in the VertexAttribPointer
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(
		1,        // must match the layout in the shader
		3,        // size
		gl.FLOAT, // type
		false,    // normalized?
		0,        // stride
		nil)      // array buffer offset
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	// What does this do?
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// make buffers from out data and tell OpenGL about them
	makeVbo(points)

	for !window.ShouldClose() {
		draw(window, program)

	}
}
