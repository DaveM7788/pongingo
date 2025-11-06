package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"image/color"

	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 1280
	screenHeight = 960
	ballSpeed    = 7.0
	paddleSpeed  = 9.0
)

type Object struct {
	X, Y, W, H float32
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt float32
	dydt float32
}

type Game struct {
	paddle    Paddle
	ball      Ball
	score     int
	highScore int
}

func (g *Game) Update() error {
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithPaddle()
	g.CollideWithWall()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen,
		float32(g.paddle.X), float32(g.paddle.Y),
		float32(g.paddle.W), float32(g.paddle.H),
		color.White, false)

	vector.FillRect(screen,
		float32(g.ball.X), float32(g.ball.Y),
		float32(g.ball.W), float32(g.ball.H),
		color.White, false)

	scoreStr := "Score: " + fmt.Sprint(g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 20, color.White)

	highScoreStr := "High Score: " + fmt.Sprint(g.highScore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 40, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= paddleSpeed
	}
}

func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0
	g.score = 0
}

func (g *Game) CollideWithWall() {
	if g.ball.X >= screenWidth {
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	} else if g.ball.Y < 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.Y >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X >= g.paddle.X && !(g.ball.X >= g.paddle.X+g.paddle.W) {
		if g.ball.Y >= g.paddle.Y && g.ball.Y <= g.paddle.Y+g.paddle.H {
			g.ball.dxdt = -g.ball.dxdt
			g.score++
			if g.score > g.highScore {
				g.highScore = g.score
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	paddle := Paddle{
		Object: Object{
			X: 1250,
			Y: 200,
			W: 15,
			H: 100,
		},
	}

	ball := Ball{
		Object: Object{
			X: 0,
			Y: 0,
			W: 15,
			H: 15,
		},
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}

	g := &Game{
		paddle: paddle,
		ball:   ball,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
