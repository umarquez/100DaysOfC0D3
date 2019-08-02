/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [22/100] Ordenamiento Gnomo | Gnome Sort
*/
package main

import (
	"fmt"
)

// GnomeSort Ordena un []int utilizando ordemaniento Gnome
func GnomeSort(input []int) {
	fmt.Print("\n===========================================\n")
	fmt.Printf("ENTRADA: %v\n", input)

	// Nuestro contador de pasos
	stepCounter := 1

	// Recorremos cada posición desde la segunda posición
	for i := 1; i < len(input); i++ {
		if i < 1 {
			i = 1
		} else if i == len(input) {
			break
		}

		v := input[i]

		fmt.Printf(
			"--------------- [ PASO %v ] ---------------\n",
			stepCounter,
		)
		fmt.Printf("%v < %v = %v\n", v, input[i-1], v < input[i-1])
		stepCounter++

		if v < input[i-1] {
			input[i] = input[i-1]
			input[i-1] = v
			i -= 2
		}

		// Estado actual del slice
		fmt.Printf("%v\n", input)
	}
	fmt.Print("===========================================\n\n")
}

func main() {
	// El slice que deseamos ordenar
	var unsorted = []int{100, 6, 34, 28, 97, 23, 61, 2, 7, 1, 44, 0}
	GnomeSort(unsorted)
	// Imprimimos el resultado final
	fmt.Printf("RESULTADO: %v", unsorted)
}
