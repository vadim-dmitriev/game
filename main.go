package main

import (
	"fmt"
	"image"
	"math/rand"
	"net"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth  = 1920
	windowHeight = 1080

	fps   = 60
	speed = 3
)

// Spaceship структура игрового объекта корабля
type Spaceship struct {
	X     float64
	Y     float64
	Angle float64
}

// NewSpaceship создает новый объект игрового объекта
func NewSpaceship() *Spaceship {
	rand.Seed(time.Now().UnixNano())

	this := &Spaceship{
		X:     rand.Float64() * windowWidth,
		Y:     rand.Float64() * windowHeight,
		Angle: 0,
	}

	return this
}

// Connection структура, описывающая соединение
type Connection net.Conn

// NewConnection создает новый объект соединения
func NewConnection() Connection {
	this := new(Connection)

	return *this
}

// Network TODO
type Network struct {
}

// NewNetwork TODO
func NewNetwork(uri string) *Network {
	this := &Network{}

	return this
}

// Run TODO
func (n *Network) Run() {
	// ..
}

// User структура описывающая игрока
type User struct {
	Username string
	Connection
	*Spaceship

	// players содержит игровые объекты других игроков
	players []Spaceship
}

// NewUser создает объект нового пользователя
func NewUser(username string) *User {
	this := &User{
		username,
		NewConnection(),
		NewSpaceship(),
		make([]Spaceship, 0),
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
	network := NewNetwork(uri)

	go network.Run()

	startGame(user)
}

func startGame(u *User) {
	pixelgl.Run(
		u.run,
	)
}

func isServerMode(uri string) bool {
	// if isServerMode(uri) {
	// 	// go listenAndServe(user)
	// 	go user.Server()
	// } else {
	// 	go user.Client()
	// }
	return true
}

// Server запускает игру в роли сервера
func (u *User) Server() {
	// listener, err := net.Listen("tcp", u)
	// if err != nil {
	// 	panic(err)
	// }
	// defer listener.Close()

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	u.players.Store(conn.RemoteAddr().String(), object{})

	// 	go handleConn(conn, o)
	// }
}

// Client запускает игру в роли клиента
func (u *User) Client() {

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
	win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, 0))

	second := time.Tick(time.Second)
	frames := 0

	imd := imdraw.New(nil)
	pic, err := loadPicture("assets/spaceship.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	bg, err := loadPicture("assets/space1.png")
	if err != nil {
		panic(err)
	}
	spriteBG := pixel.NewSprite(bg, bg.Bounds())

	imd.Color = pixel.RGB(1, 0, 0)
	fpsTimer := time.Tick(1000 / fps * time.Millisecond)
	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0

		case <-fpsTimer:
			win.Clear(colornames.White)
			// win.MakePicture()
			imd.Clear()
			spriteBG.Draw(win, pixel.IM.Moved(pixel.V(windowWidth/2, windowHeight/2)))
			sprite.Draw(win, pixel.IM.Scaled(pixel.V(0, 0), 5).Moved(pixel.V(windowWidth/2, windowHeight/2)))

			imd.Draw(win)

			win.Update()

			frames++
		}
	}
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
