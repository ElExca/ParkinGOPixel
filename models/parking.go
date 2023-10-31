package models

import "sync"

const (
	AnchoAuto = 40.0
	AltoAuto  = 80.0
)

type Auto struct {
	PosX float64
	PosY float64
	Dir  float64
}

type Estacionamiento struct {
	Espacios    int
	Mu          sync.Mutex
	Ocupados    []*Auto
	EnEspera    []*Auto
	EntradaPosX float64 // Coordenada X de la entrada
	EntradaPosY float64 // Coordenada Y de la entrada
	SalidaPosX  float64 // Coordenada X de la salida
	SalidaPosY  float64 // Coordenada Y de la salida
}

func NuevoEstacionamiento(capacidad int, entradaX, entradaY, salidaX, salidaY float64) *Estacionamiento {
	return &Estacionamiento{
		Espacios:    capacidad,
		Ocupados:    make([]*Auto, capacidad),
		EntradaPosX: entradaX,
		EntradaPosY: entradaY,
		SalidaPosX:  salidaX,
		SalidaPosY:  salidaY,
	}
}

func (e *Estacionamiento) Entrar(auto *Auto) int {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	for i, lugar := range e.Ocupados {
		if lugar == nil {
			e.Ocupados[i] = auto
			return i
		}
	}

	e.EnEspera = append(e.EnEspera, auto)
	return -1
}

func (e *Estacionamiento) Salir(i int) {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	e.Ocupados[i] = nil
	if len(e.EnEspera) > 0 {
		e.Ocupados[i] = e.EnEspera[0]
		e.EnEspera = e.EnEspera[1:]
	}
}
