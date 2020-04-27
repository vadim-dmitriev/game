package effect

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/pixelgl"
	opensimplex "github.com/ojrac/opensimplex-go"
)

// OpenSimplexNoise TODO
type OpenSimplexNoise struct {
	generator    opensimplex.Noise32
	area         pixel.Rect
	noiseColor   color.RGBA
	canvas       *pixelgl.Canvas
	canvasPixels []uint8
	scale        float32
	currentTime  float32
}

// NewOpenSimplexNoise TODO
func NewOpenSimplexNoise(area pixel.Rect, scale float32, color color.RGBA) *OpenSimplexNoise {
	seed := time.Now().Unix()

	this := &OpenSimplexNoise{
		generator:    opensimplex.NewNormalized32(seed),
		area:         area,
		noiseColor:   color,
		canvas:       nil,
		canvasPixels: make([]uint8, int(4*area.W()/10*area.H()/10)),
		scale:        scale,
	}

	return this
}

// Draw имплементирует интерфейс Effecter
func (n *OpenSimplexNoise) Draw(win *pixelgl.Window) {
	n.update()
	n.canvas.DrawColorMask(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(n.area.Center(), 10), n.noiseColor)

	n.currentTime++
}

func (n *OpenSimplexNoise) update() {
	var openSimplexValue uint8
	var x, y int

	for i := 0; i < len(n.canvasPixels); i += 4 {
		x = (i / 4) / int(n.area.H()/10.)
		y = (i / 4) % int(n.area.W()/10.)

		openSimplexValue = uint8(n.generator.Eval3(float32(x)*n.scale, float32(y)*n.scale, n.currentTime*0.0065)*94) + 1

		n.canvasPixels[i] = openSimplexValue
		n.canvasPixels[i+1] = openSimplexValue
		n.canvasPixels[i+2] = openSimplexValue
		n.canvasPixels[i+3] = openSimplexValue
	}

	n.canvas.SetPixels(n.canvasPixels)
}

// CreateCanvas TODO
func (n *OpenSimplexNoise) CreateCanvas() {
	area := pixel.R(0, 0, n.area.W()/10, n.area.H()/10)
	n.canvas = pixelgl.NewCanvas(area)
}
