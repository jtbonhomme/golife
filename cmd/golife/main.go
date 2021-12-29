package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"

	"github.com/jtbonhomme/golife/internal/version"
	"github.com/jtbonhomme/golife/pkg/game"
)

const (
	ScreenWidth  int = 1280
	ScreenHeight int = 720
)

func main() {
	log := logrus.New()
	log.Infof("golife version: %#v", version.Read())
	g := game.New(ScreenWidth, ScreenHeight)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Vector (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
