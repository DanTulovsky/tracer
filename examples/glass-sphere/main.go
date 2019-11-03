package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"

	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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

var (
	output = flag.String("output", "", "name of the output file, if empty, renders to screen")
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func env() *tracer.World {
	// width, height := 150.0, 100.0
	width, height := 400.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(1, 4, -1), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 1.7, -4.7)
	to := tracer.NewPoint(0, -1, 10)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	return w
}

func floor() *tracer.Plane {
	p := tracer.NewPlane()
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Red), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func ceiling() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().Translate(0, 5, 0))
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func backWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, 10))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Teal), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func frontWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, -5))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Purple), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func rightWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Peachpuff), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func leftWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Lightgreen), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func sphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Scale(.75, .75, .75).Translate(0, 1.75, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1.0
	s.Material().Transparency = 1.0
	s.Material().ShadowCaster = false
	s.Material().RefractiveIndex = 1.573

	return s
}

func pedestal() *tracer.Cube {
	s := tracer.NewUnitCube()
	s.SetTransform(tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = tracer.ColorName(colornames.Gold)

	return s
}

func cone() *tracer.Cone {
	s := tracer.NewClosedCone(-2, 0)
	s.SetTransform(tracer.IdentityMatrix().Translate(0, 2, 0))
	sp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Green), tracer.ColorName(colornames.Violet))
	s.Material().SetPattern(sp)
	return s
}

func background() *tracer.Group {
	g := tracer.NewGroup()
	g.AddMember(cone())
	g.SetTransform(tracer.IdentityMatrix().Translate(0, 1, 6))
	return g
}

func group(s ...tracer.Shaper) *tracer.Group {
	g := tracer.NewGroup()

	for _, s := range s {
		g.AddMember(s)
	}

	return g
}

func scene() {
	w := env()

	w.AddObject(backWall())
	w.AddObject(frontWall())
	w.AddObject(rightWall())
	w.AddObject(leftWall())
	w.AddObject(ceiling())
	w.AddObject(floor())

	w.AddObject(group(sphere(), pedestal()))
	w.AddObject(background())

	if *output != "" {
		render(w, *output)
	} else {
		renderLive(w)
	}
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

// makeVbo gives our data to OpenGL
// return the color buffer to update it later
func makeVbo(canvas *tracer.Canvas) uint32 {
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

func renderLive(w *tracer.World) {
	// setup OpenGL

	camera := w.Camera()
	width, height := int(camera.Hsize), int(camera.Vsize)
	canvas := tracer.NewCanvas(width, height)

	window := initGlfw(width, height)
	defer glfw.Terminate()

	program := initOpenGL()

	// What does this do?
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// make buffers from out data and tell OpenGL about them
	colorbo := makeVbo(canvas)

	go w.RenderLive(camera, canvas)

	for !window.ShouldClose() {
		draw(window, program, canvas, colorbo)

	}

}

func draw(window *glfw.Window, program uint32, canvas *tracer.Canvas, colorbo uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// gl.BindBuffer(gl.ARRAY_BUFFER, colorbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(canvas.Colors()), gl.Ptr(canvas.Colors()), gl.STATIC_DRAW)

	gl.DrawArrays(gl.POINTS, 0, int32(len(canvas.Points())/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

func render(w *tracer.World, output string) {
	canvas := w.Render()

	f, err := os.Create(output)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()
}
