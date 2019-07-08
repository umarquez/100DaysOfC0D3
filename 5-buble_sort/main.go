package main

import (
	"fmt"
)

// InsertionSort Ordena un []int utilizando ordemaniento por inserción
func BubleSort(input []int) {
	fmt.Print("\n===========================================\n")
	fmt.Printf("ENTRADA: %v\n", input)

	// Nuestro contador de pasos
	stepCounter := 1

	// Nuestro indicador de estado, asumimos que está desordenado
	var unsorted = true

	// Mientras se encuentre desordenado...
	for unsorted {
		unsorted = false // si no hay cambios, salimos en la siguiente vuelta.
		// Recorremos cada posición (exepto la última).
		// Nota: el siguiente for es equivalente a:
		/*for i := range input {
			if i >= len(input) - 1 { break }
		 	...
		}
		*/
		// Y tambien a:
		/*for i := 0; i < len(input)- 2; i++ {
			if i >= len(input) - 1 { break }
		 	...
		}*/
		for i, v := range input[:len(input)-1] {
			fmt.Printf(
				"--------------- [ PASO %v ] ---------------\n",
				stepCounter,
			)
			fmt.Printf("%v > %v = %v\n", v, input[i+1], v > input[i+1])
			stepCounter++

			if v > input[i+1] {
				input[i] = input[i+1]
				input[i+1] = v
				unsorted = true
			}

			// Estado actual del slice
			fmt.Printf("%v\n", input)
		}
	}
	fmt.Print("===========================================\n\n")
}

func main() {
	// El slice que deseamos ordenar
	var unsorted = []int{100, 6, 34, 28, 97, 23, 61, 2, 7, 1, 44, 0}
	BubleSort(unsorted)
	// Imprimimos el resultado final
	fmt.Printf("RESULTADO: %v", unsorted)
}
