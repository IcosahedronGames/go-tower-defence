package main

import (
	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	// Update the logical state
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render the screen
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Return the game logical screen size.
	// The screen is automatically scaled.
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 400)
	ebiten.SetWindowTitle("title")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
