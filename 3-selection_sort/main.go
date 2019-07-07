package main

import (
	"fmt"
)

// SelectionSort Ordena un []int utilizando ordemaniento por inserción
func SelectionSort(input []int) {
	fmt.Print("\n===========================================\n")
	fmt.Printf("ENTRADA: %v\n", input)

	// Nuestro contador de pasos
	stepCounter := 1

	// Iteramos de la primera a la penúltima posición del slice
	for i := 0; i < len(input)-1; i++ {
		menor := i
		// iteramos los elementos restantes del slice, buscándo el número menor
		for j := i + 1; j < len(input); j++ {
			fmt.Printf(
				"--------------- [ PASO %v ] ---------------\n",
				stepCounter,
			)
			stepCounter++

			fmt.Printf("%v > %v = %v\n", input[menor], input[j], input[menor] > input[j])
			if input[menor] > input[j] {
				menor = j
			}
		}

		v := input[i]
		input[i] = input[menor]
		input[menor] = v

		fmt.Printf("%v\n", input)
	}
	fmt.Print("===========================================\n\n")
}

func main() {
	// El slice que deseamos ordenar
	var unsorted = []int{100, 6, 34, 28, 97, 23, 61, 2, 7, 1, 44, 0}
	SelectionSort(unsorted)
	// Imprimimos el resultado final
	fmt.Printf("RESULTADO: %v", unsorted)
}
