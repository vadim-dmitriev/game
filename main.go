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
	opensimplex "github.com/ojrac/opensimplex-go"
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

	scale  = float32(.049)
	scale2 = float32(.079)
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

// user структура описывающая игрока
type user struct {
	Username string
	GO       *Spaceship

	// players содержит игровые объекты других игроков
	others []*Spaceship
}

// NewUser создает объект нового пользователя
func newUser(username string) *user {
	this := &user{
		username,
		NewSpaceship(),
		make([]*Spaceship, 0),
	}

	return this
}

func main() {
	if !isEnoughArguments() {
		exitWithFail("Использование: ./spaceshipWars адрес:порт имя_пользователя")
	}

	g, err := newGame()
	if err != nil {
		exitWithFail("Не удалось запустить игру", err)
	}

	g.start()
}

func createWindow() (*pixelgl.Window, error) {
	// TODO: Вынести настройки окна в конфиг-файл.
	windowConfig := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, windowWidth, windowHeight),
		VSync:       true,
		Undecorated: true,
	}
	win, err := pixelgl.NewWindow(windowConfig)

	return win, err
}

func isEnoughArguments() bool {
	return len(os.Args) > 2
}

// exitWithFail записывает сообщение failMessage в стандартный поток ошибок
// и завершает выполнение программы с статусом завершения 1.
// Оператор return не позволяет завершить программу с кодом завершения
// отличным от нуля, поэтому здесь используется os.Exit(code int).
// Tip: exitWithFail завершает программу в тот же момент, игнорируя
// операторы defer.
func exitWithFail(failMessage string, errs ...error) {
	if len(errs) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", failMessage)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", failMessage, errs[0].Error())
	os.Exit(1)
}

func updateBG(p opensimplex.Noise32, pRect pixel.Rect, pixels, pixels2 *[]uint8, time float32) {
	newPixels := make([]uint8, len(*pixels))
	newPixels2 := make([]uint8, len(*pixels))

	pixelsLen := int(4 * windowWidth / 10 * windowHeight / 10)
	for i := 0; i < pixelsLen; i += 4 {
		os := uint8(p.Eval3(float32(int(i/4)/(windowHeight/10))*scale, float32((i/4)%(windowWidth/10))*scale, time*0.0025)*150) + 1

		newPixels[i+0] = 175 / os
		newPixels[i+1] = 85 / os
		newPixels[i+2] = 157 / os
		newPixels[i+3] = os

		os = uint8(p.Eval3(float32(int(i/4)/(windowHeight/10))*scale2, float32((i/4)%(windowWidth/10))*scale2, time*0.003)*100) + 1

		newPixels2[i+0] = 125 / os
		newPixels2[i+1] = 249 / os
		newPixels2[i+2] = 255 / os
		newPixels2[i+3] = os

	}
	copy(*pixels, newPixels)
	copy(*pixels2, newPixels2)

}

// Draw отрисовывает спрайт моего игрока
func (u *user) draw(win *pixelgl.Window) {
	u.GO.sprite.Draw(
		win,
		pixel.IM.Rotated(pixel.V(0, 0), u.GO.Angle).Moved(u.GO.Pos),
	)
}

// DrawOthers отрисовывает спрайты других игроков
func (u *user) DrawOthers(win *pixelgl.Window) {
	for _, other := range u.others {
		other.sprite.Draw(
			win,
			pixel.IM.Moved(other.Pos),
		)
	}

}

func (u *user) handleInput(win *pixelgl.Window) {
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

func (u *user) setSprite() {
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
