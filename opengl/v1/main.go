package main

import (
	"os"
	"runtime"
	"time"
	"strings"
	"log"
	"io/ioutil"
	"fmt"
	"unsafe"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	_ "golang.org/x/image/webp"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	GameNumer = 16
	GameDenom = 9
	GameTick = time.Millisecond * 16
	GameMinMultiplier = 60
)

func init() {
	runtime.LockOSThread()
}

func main() {
	g := &Game{
		VSync: runtime.GOOS != "linux",
		Height: 600,
	}

	if g.Run() != nil {
		os.Exit(1)
	}
}
type Game struct {
	VSync bool
	Fullscreen bool
	Height int

	x int
	y int
	framerate int64
}

func (g *Game) resized(w *glfw.Window, width, height int) {
	if !g.Fullscreen {
		_, g.Height = w.GetSize()
	}
	var offsetWidth, offsetHeight int
	if width*GameDenom/GameNumer > height {
		offsetHeight = height - (height % (GameDenom * GameNumer))
		offsetWidth = offsetHeight * GameNumer / GameDenom
	} else {
		offsetWidth = width - (width % (GameNumer * GameDenom))
		offsetHeight = offsetWidth * GameDenom / GameNumer
	}
	gl.Viewport(int32(width-offsetWidth)/2, int32(height-offsetHeight)/2, int32(offsetWidth), int32(offsetHeight))
	log.Printf("frame buffer resized to %dx%d; with offsets of: %d, %d\n", offsetWidth, offsetHeight, (width-offsetWidth)/2, (height-offsetHeight)/2)
}

func (g *Game) fullscreen(w *glfw.Window, s bool) {
	if m := w.GetMonitor(); m == nil && s {
		m = glfw.GetPrimaryMonitor()
		g.x, g.y = w.GetPos()
		vm := m.GetVideoMode()
		w.SetMonitor(m, glfw.DontCare, glfw.DontCare, vm.Width, vm.Height, glfw.DontCare)
	} else if m != nil && !s {
		h := g.Height
		w.SetMonitor(nil, g.x, g.y, h*GameNumer/GameDenom, h, glfw.DontCare)
	}
}

func (g *Game) vsync(s bool) {
	if s {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

func (g *Game) keys(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	} else if key == glfw.KeyEscape {
		w.SetShouldClose(true)
	} else if key == glfw.KeyEnter && mods == glfw.ModAlt {
		g.Fullscreen = !g.Fullscreen
		g.fullscreen(w, g.Fullscreen)
	} else if key == glfw.KeyV {
		g.VSync = !g.VSync
		g.vsync(g.VSync)
	} else if key == glfw.KeyD {
		log.Printf("FPS: %d\n", int64(time.Second)/g.framerate)
	}
}

func (g *Game) input(w *glfw.Window, c rune) {}

func (g *Game) click(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}
	posX, posY := w.GetCursorPos()
	if button == glfw.MouseButtonLeft {
		log.Printf("left click received at %.3f, %.3f\n", posX, posY)
	} else if button == glfw.MouseButtonRight {
		log.Printf("right click received at %.3f, %.3f\n", posX, posY)
	}
}

func (g *Game) update() {}

func (g *Game) draw() {}

func (g *Game) createTexture(index uint32, file string) (uint32, error) {
	fin, err := os.Open(file)
	if err != nil {
		log.Printf("failed to load bg: %s\n", err)
		return 0, err
	}
	defer fin.Close()

	img, _, err := image.Decode(fin)
	if err != nil {
		log.Printf("failed to decode bg: %s\n", err)
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	flipped := image.NewRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			flipped.Set(x, img.Bounds().Dy() - y, rgba.At(x, y))
		}
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(index)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(flipped.Pix),
	)

	return texture, nil
}

func (g *Game) compileShader(shaderType uint32, file string) (uint32, error) {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(string(source))
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var ll int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &ll)

		log := strings.Repeat("\x00", int(ll+1))
		gl.GetShaderInfoLog(shader, ll, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %s: %v", source, log)
	}

	return shader, nil

}

func (g *Game) createProgram(vertexShaderFile, fragmentShaderFile string) (uint32, error) {
	vs, err := g.compileShader(gl.VERTEX_SHADER, vertexShaderFile)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vs)

	fs, err := g.compileShader(gl.FRAGMENT_SHADER, fragmentShaderFile)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fs)

	program := gl.CreateProgram()

	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var ll int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &ll)

		log := strings.Repeat("\x00", int(ll+1))
		gl.GetProgramInfoLog(program, ll, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func (g *Game) glDebug(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	if gltype != gl.DEBUG_TYPE_ERROR {
		return
	}

	log.Printf("GL Error: %s\n", message)
}

func (g *Game) Run() error {
	var err error
	if err = glfw.Init(); err != nil {
		log.Printf("failed to initialize glfw: %s\n", err)
		return err
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(g.Height*GameNumer/GameDenom, g.Height, "GOGL", nil, nil)
	if err != nil {
		log.Printf("failed to create window: %s\n", err)
		return err
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	window.SetAspectRatio(GameNumer, GameDenom)
	window.SetSizeLimits(GameMinMultiplier*GameNumer, GameMinMultiplier*GameDenom, glfw.DontCare, glfw.DontCare)
	window.SetFramebufferSizeCallback(g.resized)
	window.SetKeyCallback(g.keys)
	window.SetCharCallback(g.input)
	window.SetMouseButtonCallback(g.click)

	g.vsync(g.VSync)
	g.fullscreen(window, g.Fullscreen)

	if err = gl.Init(); err != nil {
		log.Printf("failed to initialize opengl (%s)...\n", err)
		return err
	}
	log.Printf("OpenGL Version: %s\n", gl.GoStr(gl.GetString(gl.VERSION)))

	gl.DebugMessageCallback(g.glDebug, unsafe.Pointer(nil))
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(0, 0, 0, 1)

	program, err := g.createProgram("assets/shaders/shader.vert", "assets/shaders/shader.frag")
	if err != nil {
		return err
	}
	defer gl.DeleteProgram(program)

	gl.UseProgram(program)

	projection := mgl32.Ortho(0, 1920, 0, 1080, -1, 1)
	uProjection := gl.GetUniformLocation(program, gl.Str("u_projection\x00"))
	gl.UniformMatrix4fv(uProjection, 1, false, &projection[0])

	camera := mgl32.Translate3D(-100, 0, 0)
	uCamera := gl.GetUniformLocation(program, gl.Str("u_camera\x00"))
	gl.UniformMatrix4fv(uCamera, 1, false, &camera[0])

	model := mgl32.Translate3D(200, 200, 0)
	uModel := gl.GetUniformLocation(program, gl.Str("u_model\x00"))
	gl.UniformMatrix4fv(uModel, 1, false, &model[0])

	uColor := gl.GetUniformLocation(program, gl.Str("u_color\x00"))
	gl.Uniform4f(uColor, 0.2, 0.3, 0.8, 1.0)

	uTexture := gl.GetUniformLocation(program, gl.Str("u_texture\x00"))

	positions := []float32{
		100, 100, 0.0, 0.0,
		200, 100, 1.0, 0.0,
		200, 200, 1.0, 1.0,
		100, 200, 0.0, 1.0,

	}
	var stride int32 = 4 * 4

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	defer gl.DeleteVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * 4 * 4, gl.Ptr(positions), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &vbo)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(8))

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &ibo)

	var uColorIncrement float32 = 0.01
	var uColorRed float32 = 0.2

	gl.UseProgram(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	bg, err := g.createTexture(gl.TEXTURE0, "assets/sprites/me.jpg")
	if err != nil {
		return err
	}
	defer gl.DeleteTextures(1, &bg)

	var current time.Time = time.Now()
	var previous time.Time = current
	var elapsed time.Duration
	var accumulator time.Duration
	for !window.ShouldClose() {
		glfw.PollEvents()

		current = time.Now()
		elapsed = current.Sub(previous)
		previous = current
		g.framerate = int64((float64(g.framerate) * 0.9) + (float64(elapsed.Nanoseconds()) * 0.1))

		accumulator += elapsed
		for accumulator >= GameTick {
			g.update()
			accumulator -= GameTick
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		g.draw()

		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		if uColorRed >= 1.0 {
			uColorIncrement = -0.01
		} else if uColorRed <= 0.0 {
			uColorIncrement = 0.01
		}
		uColorRed += uColorIncrement
		gl.Uniform4f(uColor, uColorRed, 0.3, 0.8, 0.2)
		gl.Uniform1i(uTexture, 0)

		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))

		window.SwapBuffers()
	}

	log.Println("Thank you for playing!")
	return nil
}
