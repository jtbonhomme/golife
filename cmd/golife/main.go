package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"

	"github.com/jtbonhomme/golife/internal/version"
	"github.com/jtbonhomme/golife/pkg/game"
)

const (
	TileDimension int = 80
	ScreenWidth   int = 16 * TileDimension // 1280
	ScreenHeight  int = 9 * TileDimension  // 720
)

func main() {
	log := logrus.New()
	log.Infof("golife version: %#v", version.Read())
	os.Setenv("EBITEN_SCREENSHOT_KEY", "s")
	g := game.New(ScreenWidth, ScreenHeight, TileDimension)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("golife (jtbonhomme@gmail.com)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
