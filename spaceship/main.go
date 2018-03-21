package main

import (
	"image"
	_ "image/png"
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"encoding/base64"
	"bytes"
	"math"
)

type Ship struct {
	Position [2]float64
	// X, Y representing direction and speed
	Velocity [2]float64
	// Radians representing direction ship is pointing
	Heading float64
	// When we "burn" our engine, by how much can we change our velocity
	Thrust float64
	// How quickly can we change our heading. Radians.
	TurnSpeed float64
	// Tick controls how frequently the ship moves. Confusingly, it
	// doesn't control the turn speed. That is controlled by frame rate.
	Tick time.Duration
	// Set to true when the user has pushed the up
	// arrow, and before we have done the movement.
	Thrusting bool
}

func (s *Ship) TurnLeft () {
	s.Heading += s.TurnSpeed
}

func (s *Ship) TurnRight () {
	s.Heading -= s.TurnSpeed
}

func (s *Ship) Burn () {
	s.Thrusting = true
}

func (s *Ship) Move () {
	go func () {
		for {
			if s.Thrusting {
				s.Velocity[0] += math.Cos(s.Heading)*s.Thrust
				s.Velocity[1] += math.Sin(s.Heading)*s.Thrust
				s.Thrusting = false
			}
			s.Position[0] += s.Velocity[0]
			s.Position[1] += s.Velocity[1]
			time.Sleep(s.Tick)
		}
	}()
}

// Effectively main()
func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Spaceship",
		Bounds: pixel.R(0, 0, 600, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	ship := Ship{
		Thrust: 0.1,
		TurnSpeed: math.Pi/50,
		Tick: time.Millisecond*10,
		Heading: math.Pi/2,
	}
	ship.Move()

	// Load pictures
	shipSprite := spriteFromImage(File_shippng)
	burnSprite := spriteFromImage(File_burnpng)
	leftSprite := spriteFromImage(File_leftpng)
	rightSprite := spriteFromImage(File_rightpng)

	for !win.Closed() {

		if win.Pressed(pixelgl.KeyLeft) {
				ship.TurnLeft()
		}
		if win.Pressed(pixelgl.KeyRight) {
				ship.TurnRight()
		}
		if win.Pressed(pixelgl.KeyUp) {
			ship.Burn()
		}

		win.Clear(colornames.Darkblue)

		mat := pixel.IM
		mat = mat.Rotated(pixel.ZV, ship.Heading-(math.Pi/2))
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.Moved(pixel.V(ship.Position[0], ship.Position[1]))
		shipSprite.Draw(win, mat)
		if win.Pressed(pixelgl.KeyUp) {
				burnSprite.Draw(win, mat)
		}
		if win.Pressed(pixelgl.KeyLeft) {
				leftSprite.Draw(win, mat)
		}
		if win.Pressed(pixelgl.KeyRight) {
				rightSprite.Draw(win, mat)
		}

		win.Update()
		time.Sleep(time.Second/120)
	}
}

func main() {
	pixelgl.Run(run)
}

func spriteFromImage (imgData string) *pixel.Sprite {
	byteData, err := base64.StdEncoding.DecodeString(imgData)
    if err != nil {
        panic(err)
    }
	reader := bytes.NewReader(byteData)
	img, _, err := image.Decode(reader)
    if err != nil {
        panic(err)
    }
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds())
}
