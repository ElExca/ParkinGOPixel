package views

import (
	"math/rand"
	"strconv"
	"time"

	"modulos/models"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

const (
	Velocidad           = 1.0
	AnchoAuto           = 20.0
	AltoAuto            = 20.0
	AltoEspacio         = 20.0
	DistanciaEntreAutos = 10.0
)

func DibujarEspacios(imd *imdraw.IMDraw, e *models.Estacionamiento, win *pixelgl.Window, atlas *text.Atlas) {
	for i := 0; i < 20; i++ {
		imd.Color = pixel.RGB(0, 0, 0)
		x := float64(i) * (AnchoAuto + DistanciaEntreAutos)
		y := AltoAuto

		imd.Push(pixel.V(x, y))
		imd.Push(pixel.V(x+AnchoAuto, y+AltoEspacio))
		imd.Rectangle(0)

		//Dibuja

		txt := text.New(pixel.V(x+AnchoAuto/2, y+AltoEspacio/.5), atlas) // Ajusta el valor 1.5 para cambiar la posición vertical del número
		txt.Color = pixel.RGB(255, 255, 255)
		txt.WriteString(strconv.Itoa(i + 1))
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 1))

	}
}

func Run(win *pixelgl.Window, e *models.Estacionamiento) {
	rand.Seed(time.Now().UnixNano())
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	go func() {
		for {
			auto := &models.Auto{PosX: -AnchoAuto - DistanciaEntreAutos, PosY: AltoAuto + AltoEspacio, Dir: 1}
			pos := e.Entrar(auto)

			if pos != -1 {
				go func(p int) {
					// Dormir un tiempo antes de que el auto salga
					time.Sleep(time.Duration(rand.Intn(15)+5) * time.Second)

					// Cambiar la dirección del auto para que salga por el mismo lado
					auto.Dir = -1
					e.Salir(p)
				}(pos)
			}

			// tiempo antes crear el siguiente auto
			time.Sleep(time.Millisecond * 1500)
		}
	}()

	espacioEntreAutos := AnchoAuto + DistanciaEntreAutos

	for !win.Closed() {
		win.Clear(pixel.RGB(1, 0, 1))

		im := imdraw.New(nil)
		DibujarEspacios(im, e, win, atlas)

		e.Mu.Lock()
		for i, auto := range e.Ocupados {
			if auto != nil {
				// Si el auto ha llegado a la posición de estacionamiento
				if auto.PosX >= espacioEntreAutos*float64(i) {
					// Dibujar el auto en su lugar de estacionamiento
					auto.PosX = espacioEntreAutos * float64(i)
					auto.PosY = AltoAuto + AltoEspacio

					im.Color = pixel.RGB(0, 0, 1)
					im.Push(pixel.V(auto.PosX, auto.PosY))
					im.Push(pixel.V(auto.PosX+AnchoAuto, auto.PosY+AltoAuto))
					im.Rectangle(0)
				} else {
					// Dibujar el auto arriba antes de estacionarse
					auto.PosX += Velocidad * auto.Dir
					auto.PosY = AltoAuto + AltoEspacio*2 // Cambiar la altura
					im.Color = pixel.RGB(0, 0, 1)
					im.Push(pixel.V(auto.PosX, auto.PosY))
					im.Push(pixel.V(auto.PosX+AnchoAuto, auto.PosY+AltoAuto))
					im.Rectangle(0)
				}
			}
		}
		e.Mu.Unlock()

		im.Draw(win)
		win.Update()
	}
}
