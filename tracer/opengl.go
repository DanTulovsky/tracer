package tracer

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
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
	gl_PointSize = 2.0;

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

func draw(window *glfw.Window, program uint32, canvas *Canvas, colorbo uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// gl.BindBuffer(gl.ARRAY_BUFFER, colorbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(canvas.Colors()), gl.Ptr(canvas.Colors()), gl.STATIC_DRAW)

	gl.DrawArrays(gl.POINTS, 0, int32(len(canvas.Points())/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

// RenderLive renders using OpenGL to the screen
func RenderLive(w *World) {
	// setup OpenGL

	camera := w.Camera()
	width, height := int(camera.Hsize), int(camera.Vsize)
	canvas := NewCanvas(width, height)

	window := initGlfw(width, height)
	defer glfw.Terminate()

	program := initOpenGL()

	// What does this do?
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// make buffers from out data and tell OpenGL about them
	colorbo := makeVbo(canvas)

	go w.Render(camera, canvas)

	for !window.ShouldClose() {
		draw(window, program, canvas, colorbo)

	}
}

// Render renders to a file
func Render(w *World, output string) {
	camera := w.Camera()
	width, height := int(camera.Hsize), int(camera.Vsize)
	canvas := NewCanvas(width, height)

	w.Render(camera, canvas)

	f, err := os.Create(output)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

// makeVbo gives our data to OpenGL
// return the color buffer to update it later
func makeVbo(canvas *Canvas) uint32 {
	// points buffer
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(canvas.Points()), gl.Ptr(canvas.Points()), gl.STATIC_DRAW)

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
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(canvas.Colors()), gl.Ptr(canvas.Colors()), gl.STATIC_DRAW)

	// the 1 here is the buffer in the VertexAttribPointer
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(
		1,        // must match the layout in the shader
		3,        // size
		gl.FLOAT, // type
		false,    // normalized?
		0,        // stride
		nil)      // array buffer offset

	return colorbo
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

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(width, height int) *glfw.Window {
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
