package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jtbonhomme/golife/pkg/game"
)

const (
	ScreenWidth  int = 1280
	ScreenHeight int = 720
)

func main() {
	emptyImage := ebiten.NewImage(3, 3)
	emptyImage.Fill(color.White)

	g := game.New(ScreenWidth, ScreenHeight)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
