package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const defaultCapacity = 5000
const apartments = 6
const consumptionProbability = 0.18
const consumptionLiters = 5

// tank es el tipo de datos interno, del que solo necesitamos una instancia,
// no necesitamos exportarla pues no queremos que nadie más pueda instanciarla
// desde fuera del paquete, si se tratara de una librería o plug-in.
type tank struct {
	capacity float64
	level    float64
}

// Fill es el método que nos permitirá alimentar el tanque de agua, este es
// público pues es uno de los que nos permitrá interactuar con los valores
// protejidos de la instancia.
func (t *tank) Fill(liters float64) {
	t.level = liters
	return
}

// Flush es el método que nos permite tomar agua del tanque, este es
// público pues es uno de los que nos permitrá interactuar con los valores
// protejidos de la instancia.
func (t *tank) Flush(liters float64) {
	t.level -= liters
	return
}

// GetCapacity devuelve la capacidad actual del tanque, este es
// público pues es uno de los que nos permitrá interactuar con los valores
// protejidos de la instancia.
func (t *tank) GetCapacity() float64 {
	return t.capacity
}

// GetCurrentLevel devuelve el nivel actual de llenado del tanquem, este es
// público pues es uno de los que nos permitrá interactuar con los valores
// protejidos de la instancia.
func (t *tank) GetCurrentLevel() float64 {
	return t.level
}

// WaterTank funciona como una máscara, solo nos permitirá acceder a los
// métodos definidor por ella, protegiendo el resto de la instancia, asegurando
// que esta no será modificada desde fuera.
type WaterTank interface {
	Fill(liters float64)
	Flush(liters float64)
	GetCapacity() float64
	GetCurrentLevel() float64
}

// tankInstance es la instancia única, esta solo será accesible através de los
// métodos públicos del struct que la define y la interface a la que satisface.
var tankInstance *tank

// GetTank es la función que nos permite obtener la intancia única, si esta no
// está inicializada la misma función se encarga de ello.
func GetTank() WaterTank {
	if tankInstance == nil {
		tankInstance = new(tank)
		tankInstance.capacity = defaultCapacity
	}

	// Retornando la misma instancia cada vez desde la priemra vez...
	return tankInstance
}

func main() {
	rand.Seed(time.Now().Unix())

	// Rellenamos el tanque y de paso inicializamos la instancia
	GetTank().Fill(defaultCapacity)

	var minCounter int
	// Repetimos mientras el tanque tenga agua
	for GetTank().GetCurrentLevel() > 0 {
		// +1 minuto
		minCounter++

		wg := sync.WaitGroup{}
		wg.Add(apartments)
		// Inicia una goroutine por cada departamento
		for i := 0; i < apartments; i++ {
			go func() {
				if rand.Float32() < consumptionProbability {
					// Cada goroutine consume del mismo tanque
					GetTank().Flush(consumptionLiters)
				}
				wg.Done()
			}()
		}
		// Esperamos a que terminen las rutinas
		wg.Wait()
	}

	// Imprimiendo el tiempo total en que tomó consumir toda el agua.
	fmt.Printf("%v hours with %v apartments and a %v liters tank\n",
		math.Round((float64(minCounter)/60)*100)/100,
		apartments,
		defaultCapacity,
	)

	os.Exit(0)
}
