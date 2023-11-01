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

		txt := text.New(pixel.V(x+AnchoAuto/2, y+AltoEspacio/.5), atlas)
		txt.WriteString(strconv.Itoa(i + 1))
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 1))

	}
}

func Run(win *pixelgl.Window, e *models.Estacionamiento) {
	rand.Seed(time.Now().UnixNano())
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	go func() {
		for {
			auto := &models.Auto{PosX: -AnchoAuto - DistanciaEntreAutos, PosY: AltoAuto + AltoEspacio, Dir: 1, State: models.StateEntering}
			pos := e.Entrar(auto)

			if pos != -1 {
				go func(p int) {
					time.Sleep(time.Duration(rand.Intn(15)+5) * time.Second)
					auto.State = models.StateExiting
				}(pos)
			}

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
				if auto.State == models.StateEntering {
					auto.PosX += Velocidad * auto.Dir
					auto.PosY = AltoAuto + AltoEspacio*2
					im.Color = pixel.RGB(0, 0, 1)
					im.Push(pixel.V(auto.PosX, auto.PosY))
					im.Push(pixel.V(auto.PosX+AnchoAuto, auto.PosY+AltoAuto))
					im.Rectangle(0)

					if auto.PosX >= espacioEntreAutos*float64(i) {
						auto.State = models.StateParked
					}
				} else if auto.State == models.StateParked {
					auto.PosX = espacioEntreAutos * float64(i)
					auto.PosY = AltoAuto + AltoEspacio
					im.Color = pixel.RGB(0, 0, 1)
					im.Push(pixel.V(auto.PosX, auto.PosY))
					im.Push(pixel.V(auto.PosX+AnchoAuto, auto.PosY+AltoAuto))
					im.Rectangle(0)
				} else if auto.State == models.StateExiting {

					auto.PosX += Velocidad * auto.Dir
					auto.PosY = AltoAuto + AltoEspacio*2
					im.Color = pixel.RGB(0, 0, 1)
					im.Push(pixel.V(auto.PosX, auto.PosY))
					im.Push(pixel.V(auto.PosX+AnchoAuto, auto.PosY+AltoAuto))
					im.Rectangle(0)

					if auto.PosX <= -AnchoAuto-DistanciaEntreAutos {

						e.Ocupados[i] = nil
					}
				}
			}
		}
		e.Mu.Unlock()

		im.Draw(win)
		win.Update()
	}
}
