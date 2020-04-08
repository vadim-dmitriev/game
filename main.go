package main

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"net"
	"os"
	"sync"
	"time"

	_ "image/png"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	windowWidth  = 500
	windowHeight = 500

	fps   = 60
	speed = 3
)

type bullet struct {
	X, Y  float64
	Angle float64
}

type object struct {
	Hp   int
	X, Y float64
	// Xm, Ym  float64
	Angle   float64
	Bullets []bullet
}

// User .
type User struct {
	uri     string
	mode    string
	iam     object
	clients sync.Map
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./game addr:port")
		os.Exit(1)
	}

	var o = &User{
		uri:     os.Args[1],
		mode:    "client",
		clients: sync.Map{},
		iam: object{
			Hp: 100,
		},
	}

	o.startServer()

	pixelgl.Run(o.run)
}

func (o *User) startServer() {
	conn, err := net.Dial("tcp", o.uri)
	if err != nil {
		fmt.Println("server not fount. I AM A SERVER!")
		o.mode = "server"
		go listenAndServe(o)
		return
	}

	go client(conn, o)
}

func (o *User) run() {
	cfg := pixelgl.WindowConfig{
		Title:  o.mode,
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, 0))

	second := time.Tick(time.Second)
	frames := 0

	imd := imdraw.New(nil)
	pic, err := loadPicture("spaceship.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	imd.Color = pixel.RGB(1, 0, 0)
	fpsTimer := time.Tick(1000 / fps * time.Millisecond)
	for !win.Closed() && !win.JustPressed(pixelgl.KeyEscape) {
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			fmt.Println(len(o.iam.Bullets))
			frames = 0

		case <-fpsTimer:
			win.Clear(colornames.White)
			imd.Clear()
			if o.iam.Hp <= 0 {
				os.Exit(0)
			}

			if win.Pressed(pixelgl.KeyW) {
				o.iam.Y += speed
			}
			if win.Pressed(pixelgl.KeyS) {
				o.iam.Y -= speed
			}
			if win.Pressed(pixelgl.KeyA) {
				o.iam.X -= speed
			}
			if win.Pressed(pixelgl.KeyD) {
				o.iam.X += speed
			}
			pos := win.MousePosition()
			if pos.X > o.iam.X {
				o.iam.Angle = math.Atan((pos.Y - o.iam.Y) / (pos.X - o.iam.X))
			} else {
				o.iam.Angle = math.Atan((pos.Y-o.iam.Y)/(pos.X-o.iam.X)) - math.Pi
			}
			if win.JustPressed(pixelgl.MouseButton1) {
				o.iam.Bullets = append(o.iam.Bullets, bullet{o.iam.X, o.iam.Y, o.iam.Angle})
			}
			sprite.Draw(win, pixel.IM.Rotated(pixel.V(0, 0), o.iam.Angle).Moved(pixel.V(o.iam.X, o.iam.Y)))

			newBullets := o.iam.Bullets[:0]
			for i, b := range o.iam.Bullets {
				if b.X < windowWidth && b.X > 0 && b.Y < windowHeight && b.Y > 0 {
					newBullets = append(newBullets, b)
				}
				o.clients.Range(func(k, v interface{}) bool {
					if isAim(b, v.(object)) {
						fmt.Println("hited")
						// o.iam.Hp--
						newClient := v.(object)
						newClient.Hp--
						o.clients.Store(k, newClient)
					}
					return true // if false, Range stops
				})
				imd.Push(pixel.V(b.X, b.Y))
				imd.Circle(2, 0)

				o.iam.Bullets[i].X += math.Cos(b.Angle)
				o.iam.Bullets[i].Y += math.Sin(b.Angle)
			}
			o.iam.Bullets = newBullets

			o.clients.Range(func(k, v interface{}) bool {
				sprite.Draw(win, pixel.IM.Rotated(pixel.V(0, 0), v.(object).Angle).Moved(pixel.V(v.(object).X, v.(object).Y)))
				for _, b := range v.(object).Bullets {
					imd.Push(pixel.V(b.X, b.Y))
					imd.Circle(2, 0)
				}
				return true // if false, Range stops
			})

			imd.Draw(win)

			win.Update()

			frames++
		}
	}
}

func listenAndServe(o *User) {
	listener, err := net.Listen("tcp", o.uri)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		o.clients.Store(conn.RemoteAddr().String(), object{})

		go handleConn(conn, o)
	}
}

func handleConn(conn net.Conn, o *User) {
	defer conn.Close()

	type msg map[string]object

	var err error
	ticker := time.Tick(1000 / fps * time.Millisecond)

	for {
		select {
		case <-ticker:
			var payload = map[string]object{
				o.uri: o.iam,
			}

			o.clients.Range(func(k, v interface{}) bool {
				if conn.RemoteAddr().String() != k {
					payload[k.(string)] = v.(object)
				}
				return true
			})

			// СЕРВЕР ОТПРАВЛЯЕТ КАЖДОМУ КЛИЕНТУ ДАННЫЕ ОБ !!!ОСТАЛЬНЫХ!!!
			// КЛИЕНТАХ, ВКЛЮЧАЯ СЕБЯ - СЕРВЕР!
			err = json.NewEncoder(conn).Encode(payload)
			if err != nil {
				continue
			}

			// СЕРВЕР, В РАМКАХ ОДНОГО СОЕДИНЕНИЯ, ДЕКОДИРУЕТ НОВЫЕ ДАННЫЕ
			// ОТ КЛИЕНТА И СОХРАНЯЕТ ИХ К СЕБЕ!
			// !!!ДОЛГО ЕБАЛСЯ!!!! С КОНКУРЕНТНОЙ ЗАПИСЬЮ!!!!
			m := msg{}
			err = json.NewDecoder(conn).Decode(&m)
			if err != nil {
				continue
			}

			o.clients.Store(conn.RemoteAddr().String(), m[conn.RemoteAddr().String()])
		}
	}

}

func client(conn net.Conn, o *User) {
	defer conn.Close()

	payload := map[string]object{
		conn.LocalAddr().String(): o.iam,
	}
	err := json.NewEncoder(conn).Encode(payload)
	if err != nil {
		panic(err)
	}
	type msg map[string]object

	ticker := time.Tick(1000 / fps * time.Millisecond)

	for {
		select {
		case <-ticker:
			m := msg{}
			err = json.NewDecoder(conn).Decode(&m)
			if err != nil {
				continue
			}
			for k, v := range m {
				o.clients.Store(k, v)
			}

			payload := map[string]object{
				conn.LocalAddr().String(): o.iam,
			}
			err = json.NewEncoder(conn).Encode(payload)
			if err != nil {
				continue
			}

		}
	}
}

func isAim(b bullet, enemy object) bool {
	if b.X >= (enemy.X-10) && b.X <= (enemy.X+10) {
		if b.X >= (enemy.Y-10) && b.Y <= (enemy.Y+10) {
			return true
		}
	}
	return false
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
