package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"github.com/vadim-dmitriev/game/effect"
	"golang.org/x/image/colornames"
)

type effecter interface {
	Draw(*pixelgl.Window)
	CreateCanvas()
}

type world struct {
	backgroundSprite *pixel.Sprite
	backgroundMatrix pixel.Matrix

	effects []effecter

	perlinCanvas *pixelgl.Canvas
}

func newWorld() (*world, error) {
	bgPic, err := loadPicture("assets/space1.png")
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать фоновое изображение")
	}
	bgSprite := pixel.NewSprite(bgPic, bgPic.Bounds())

	this := &world{
		backgroundSprite: bgSprite,
		backgroundMatrix: pixel.IM.ScaledXY(pixel.V(0, 0), pixel.V(windowWidth/bgSprite.Frame().W(), windowHeight/bgSprite.Frame().H())).Moved(pixel.V(windowWidth/2, windowHeight/2)),
		effects:          make([]effecter, 0),
	}

	this.effects = append(this.effects,
		effect.NewOpenSimplexNoise(pixel.R(0, 0, windowWidth, windowHeight), 0.07, color.RGBA{155, 151, 208, 255}),
		effect.NewOpenSimplexNoise(pixel.R(0, 0, windowWidth, windowHeight), 0.01, color.RGBA{27, 26, 38, 255}),
	)

	return this, nil

}

func (w *world) update(win *pixelgl.Window) {
	win.Clear(colornames.White)
	w.backgroundSprite.Draw(win, w.backgroundMatrix)

	for _, e := range w.effects {
		e.Draw(win)
	}

}

func (w *world) InitCanvases() {
	for _, effect := range w.effects {
		effect.CreateCanvas()
	}
}
