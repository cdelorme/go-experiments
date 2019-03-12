package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, _ := glfw.CreateWindow(windowWidth, windowHeight, "glfontExample", nil, nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	// @note: load font (fontfile, font scale, window width, window height
	// vera bitstream license is 100% legal to redistribute as a component of
	// other software, making it not only a pretty font face but a great asset!
	font, err := glfont.LoadFont("VeraMono.ttf", 52, windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}
	font.SetColor(1.0, 1.0, 1.0, 1.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !window.ShouldClose() {
		glfw.PollEvents()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		font.Printf(16, 32, 0.5, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
		window.SwapBuffers()
	}
}
