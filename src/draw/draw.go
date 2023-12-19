package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DrawPlugin struct {
}


func(d *DrawPlugin) Draw(screen *ebiten.Image){
	ebitenutil.DebugPrint(screen, "Try to change this string in src/draw/draw.go")
}

var Export DrawPlugin