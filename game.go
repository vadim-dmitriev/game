package main

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
)

type game struct {
	user    *user
	world   *world
	network *Network
}

func newGame() (*game, error) {
	w, err := newWorld()
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать мир")
	}
	this := &game{
		user:    newUser("aa"),
		world:   w,
		network: nil,
	}

	return this, nil
}

func (g *game) start() {

	g.user.setSprite()

	pixelgl.Run(
		g.run,
	)
}

func (g *game) run() {
	win, err := createWindow()
	if err != nil {
		exitWithFail("Не удалось запустить игру", err)
	}
	defer win.Destroy()

	fpsTimer := time.Tick(1000 / fps * time.Millisecond)

	world := g.world
	user := g.user

	world.InitCanvases()

	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		<-fpsTimer
		world.update(win)

		user.handleInput(win)
		user.draw(win)

		win.Update()
	}

}
