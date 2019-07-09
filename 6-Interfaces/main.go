package main

import (
	"fmt"
	"math"
)

// PuedeSerMartillo Define que: cualquier cosa que pueda Golpear(), entonces
// "PuedeSerMartillo" y será usado como tal! (En caso necesario).
type PuedeSerMartillo interface {
	Golpear(clavo float64)
}

// =========================================
// LlaveDeTuercas Es una herramienta que no necesariamente es un Martillo
type LlaveDeTuercas struct {
	pesoKg float64
}

// NewLlaveDeTuercas Devuelve una nueva llave de tuercas.
func NewLlaveDeTuercas() LlaveDeTuercas {
	return LlaveDeTuercas{
		pesoKg: .6,
	}
}

// Golpear Este método hace que la llave de tuercas pueda ser usada como Martillo
func (llave LlaveDeTuercas) Golpear(clavo float64) {
	var golpes = math.Ceil(clavo /llave.pesoKg)

	fmt.Printf("Se necesitaron %v golpes con la llave de tuercas\n", golpes)
}

// =========================================
// Destornillador Es una herramienta que no necesariamente es un Martillo
type Destornillador struct {
	pesoKg float64
}

// NewDestornillador Devuelve un nuevo destornillador.
func NewDestornillador() Destornillador {
	return Destornillador{
		pesoKg: .3,
	}
}

// Golpear Este método hace que el destornillador pueda ser usado como Martillo
func (desT Destornillador) Golpear(clavo float64) {
	var golpes = math.Ceil(clavo /desT.pesoKg)

	fmt.Printf("Se necesitaron %v golpes con el destornillador\n",golpes)
}

// =========================================
// Alicatas Es una herramienta que no necesariamente es un Martillo
type Alicatas struct {
	pesoKg float64
}

// NewAlicatas Devuelve nuevas alicatas.
func NewAlicatas() Alicatas {
	return Alicatas{
		pesoKg: .5,
	}
}

// Golpear Este método hace que las alicatas puedan ser usadas como Martillo
func (alic Alicatas) Golpear(clavo float64) {
	var golpes = math.Ceil(clavo /alic.pesoKg)

	fmt.Printf("Se necesitaron %v golpes con las alicatas\n",golpes)
}

// =========================================
// Piedra Es un objeto que no es herramienta pero puede Golpear()
type Piedra struct {
	pesoKg float64
}

// NewPiedra Devuelve una nueva piedra.
func NewPiedra() Piedra {
	return Piedra{
		pesoKg: .85,
	}
}

// Golpear Este método hace que las alicatas puedan ser usadas como Martillo
func (piedra Piedra) Golpear(clavo float64) {
	var golpes = math.Ceil(clavo /piedra.pesoKg)

	fmt.Printf("Se necesitaron %v golpes con la piedra\n",golpes)
}

// =========================================
// Pincel es un objeto que no puede Golpear()
type Pincel struct {}

// NewPincel Devuelve un nuevo pincel.
func NewPincel() Pincel {
	return Pincel{}
}

// =========================================
// ClavarClavoCon Usa un objeto que PuedeSerMartillo para clavar un clavo...
func ClavarClavoCon(clavo float64, noMartillo PuedeSerMartillo) {
	noMartillo.Golpear(clavo)
}

func main() {
	// Definimos el tamaño ficticio del clavo
	var clavo float64 = 8//cm
	fmt.Printf("-  Probando un clavo de %vcm\n", clavo)
	fmt.Println()

	// Intentamos golpearlo con cada uno de los objetos; se puede hacer de
	// manera directa de la forma: ClavarClavoCon(clavo, NewAlicatas())
	// pero decidimos demostrar que son tipos diferentes.
	a := NewAlicatas()
	fmt.Printf("Usando %T\n", a)
	ClavarClavoCon(clavo, a)
	fmt.Println()

	d := NewDestornillador()
	fmt.Printf("Usando %T\n", d)
	ClavarClavoCon(clavo, d)
	fmt.Println()

	ll := NewLlaveDeTuercas()
	fmt.Printf("Usando %T\n", ll)
	ClavarClavoCon(clavo, ll)
	fmt.Println()

	p := NewPiedra()
	fmt.Printf("Usando %T\n", p)
	ClavarClavoCon(clavo, p)
	fmt.Println()

	pi := NewPincel()
	fmt.Printf("Usando %T\n", pi)
	// Pincel no puede Golpear()
	//ClavarClavoCon(clavo, pi) // Error, no retorna un PuedeSerMartillo
}