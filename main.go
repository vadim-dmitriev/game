package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/vadim-dmitriev/game/network"
)

const (
	windowWidth  = 1920
	windowHeight = 1080

	// windowWidth  = 500
	// windowHeight = 250

	fps   = 600
	speed = 3

	speedLimit = 5
	friction   = 0.96
)

var (
	// OX единичный вектор оси X
	OX = pixel.Vec{
		X: 1,
		Y: 0,
	}

	// OY единичный вектор оси Y
	OY = pixel.Vec{
		X: 0,
		Y: 1,
	}
)

// Spaceship структура игрового объекта корабля
type Spaceship struct {
	Pos   pixel.Vec
	Angle float64

	velocityVec pixel.Vec
	sprite      *pixel.Sprite
}

// NewSpaceship создает новый объект игрового объекта
func NewSpaceship() *Spaceship {
	rand.Seed(time.Now().UnixNano())

	this := &Spaceship{
		Pos: pixel.Vec{
			X: rand.Float64() * windowWidth,
			Y: rand.Float64() * windowHeight,
		},
		Angle: 0,

		velocityVec: pixel.Vec{},
		sprite:      nil,
	}

	return this
}

// User структура описывающая игрока
type User struct {
	Username string
	GO       *Spaceship

	// players содержит игровые объекты других игроков
	others []*Spaceship
}

// NewUser создает объект нового пользователя
func NewUser(username string) *User {
	this := &User{
		username,
		NewSpaceship(),
		make([]*Spaceship, 0),
	}

	return this
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./game addr:port username")
		return
	}

	uri, username := os.Args[1], os.Args[2]

	user := NewUser(username)
	network := network.New(uri)

	network.Run()

	startGame(user)
}

func startGame(u *User) {
	u.setSprite()

	pixelgl.Run(
		u.run,
	)
}

func (u *User) run() {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, 0))

	second := time.Tick(time.Second)
	frames := 0

	fpsTimer := time.Tick(1000 / fps * time.Millisecond)
	bg, err := loadPicture("assets/space2.png")
	if err != nil {
		panic(err)
	}
	spriteBG := pixel.NewSprite(bg, bg.Bounds())
	bgIM := pixel.IM.ScaledXY(pixel.V(0, 0), pixel.V(windowWidth/spriteBG.Frame().W(), windowHeight/spriteBG.Frame().H()))

	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0

		case <-fpsTimer:
			win.Clear(colornames.White)

			u.handleInput(win)

			spriteBG.Draw(win, bgIM.Moved(pixel.V(windowWidth/2, windowHeight/2)))
			u.Draw(win)
			u.DrawOthers(win)

			win.Update()

			frames++
		}
	}
}

// Draw отрисовывает спрайт моего игрока
func (u *User) Draw(win *pixelgl.Window) {
	u.GO.sprite.Draw(
		win,
		pixel.IM.Rotated(pixel.V(0, 0), u.GO.Angle).Moved(u.GO.Pos),
	)
}

// DrawOthers отрисовывает спрайты других игроков
func (u *User) DrawOthers(win *pixelgl.Window) {
	for _, other := range u.others {
		other.sprite.Draw(
			win,
			pixel.IM.Moved(other.Pos),
		)
	}

}

func (u *User) handleInput(win *pixelgl.Window) {
	// вычисляем и устанавливаем направление взора
	mousePosition := win.MousePosition()
	u.GO.Angle = calcDirectionAngle(u.GO.Pos, mousePosition)

	var isBtnPressed bool

	if win.Pressed(pixelgl.KeyW) && u.GO.velocityVec.Project(OY).Y < speedLimit {
		u.GO.velocityVec.Y++
		isBtnPressed = true

	} else if win.Pressed(pixelgl.KeyS) && u.GO.velocityVec.Project(OY).Y > -speedLimit {
		u.GO.velocityVec.Y--
		isBtnPressed = true

	}

	if win.Pressed(pixelgl.KeyA) && u.GO.velocityVec.Project(OX).X > -speedLimit {
		u.GO.velocityVec.X--
		isBtnPressed = true

	} else if win.Pressed(pixelgl.KeyD) && u.GO.velocityVec.Project(OX).X < speedLimit {
		u.GO.velocityVec.X++
		isBtnPressed = true
	}

	if !isBtnPressed {
		u.GO.velocityVec = u.GO.velocityVec.Scaled(friction)
	}

	u.GO.Pos.X += u.GO.velocityVec.
		Project(pixel.V(1, 0)).
		Scaled(speed).
		X

	u.GO.Pos.Y += u.GO.velocityVec.
		Project(pixel.V(0, 1)).
		Scaled(speed).
		Y

}

func (u *User) setSprite() {
	pic, err := loadPicture("assets/spaceship.png")
	if err != nil {
		panic(err)
	}
	u.GO.sprite = pixel.NewSprite(pic, pic.Bounds())
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
