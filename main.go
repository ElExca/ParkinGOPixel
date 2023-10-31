package main

import (
	"modulos/models"
	"modulos/views"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(func() {
		cfg := pixelgl.WindowConfig{
			Title:  "Estacionamiento",
			Bounds: pixel.R(0, 0, 800, 600),
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		e := models.NuevoEstacionamiento(20)
		views.Run(win, e)
	})
}
