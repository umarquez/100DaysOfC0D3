package main

import (
	"fmt"
)

// InsertionSort Ordena un []int utilizando ordemaniento por inserción
func InsertionSort(input []int) {
	fmt.Print("\n===========================================\n")
	fmt.Printf("ENTRADA: %v\n", input)

	// Nuestro contador de pasos
	stepCounter := 1

	// Vamos a evaluar cada valor del slice
	for i, v := range input {
		// El primero lo omitimos pues no hay nada con qué compararlo.
		if i < 1 {
			continue
		}

		// Recorremos cada posición anterior hasta encontrar un valor menor al
		// actual
		for j := i - 1; j >= 0; j-- {
			fmt.Printf(
				"--------------- [ PASO %v ] ---------------\n",
				stepCounter,
			)
			fmt.Printf("%v > %v = %v\n", input[j], v, input[j] > v)
			stepCounter++

			// Si el valor anterior es menor o igual al que estamos comparando
			// hemos encontrado su posición final
			if input[j] <= v {
				fmt.Printf("%v\n", input)
				break
			}

			// De lo contrario, intercambiamos los valores de ambas posiciones
			input[i] = input[j]
			input[j] = v
			// actualizamos la posición del número que estamos evaluando
			i--

			// Estado actual del slice
			fmt.Printf("%v\n", input)
		}
	}
	fmt.Print("===========================================\n\n")
}

func main() {
	// El slice que deseamos ordenar
	var unsorted = []int{100, 6, 34, 28, 97, 23, 61, 2, 7, 1, 44, 0}
	InsertionSort(unsorted)
	// Imprimimos el resultado final
	fmt.Printf("RESULTADO: %v", unsorted)
}
