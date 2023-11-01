package models

import "sync"

const (
	AnchoAuto               = 40.0
	AltoAuto                = 80.0
	StateEntering AutoState = iota
	StateParked
	StateExiting
)

type AutoState int

type Auto struct {
	PosX  float64
	PosY  float64
	Dir   float64
	State AutoState
}

type Estacionamiento struct {
	Espacios int
	Mu       sync.Mutex
	Ocupados []*Auto
	EnEspera []*Auto
}

func NuevoEstacionamiento(capacidad int, entradaX, entradaY, salidaX, salidaY float64) *Estacionamiento {
	return &Estacionamiento{
		Espacios: capacidad,
		Ocupados: make([]*Auto, capacidad),
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

func (e *Estacionamiento) Salir(auto *Auto) int {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	for i, a := range e.Ocupados {
		if a == auto {
			e.Ocupados[i] = nil
			return i
		}
	}
	return -1
}
