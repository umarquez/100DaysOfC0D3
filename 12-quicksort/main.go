package main

import (
	"fmt"
)

// Nuestro contador de rondas
var roundCounter int

func Quicksort(input []int) []int {
	if len(input) <= 1 {
		return input
	} else if len(input) == 2 {
		if input[0] < input[1] {
			return input
		}
		return append([]int{}, input[1], input[0])
	}

	roundCounter++
	fmt.Printf(
		"--------------- [ Ronda %v ] ---------------\n",
		roundCounter,
	)

	midIx := len(input) / 2
	pivot := input[midIx]

	var left, right []int

	for i := 0; i <= midIx; i++ {
		//fmt.Printf("ix: %v, midIx: %v\n", i, midIx)
		if midIx-i >= 0 && midIx > midIx-i {
			if input[midIx-i] < pivot {
				left = append(left, input[midIx-i])
			} else {
				right = append(right, input[midIx-i])
			}
		}

		if midIx+i < len(input) && midIx < midIx+i {
			if input[midIx+i] < pivot {
				left = append(left, input[midIx+i])
			} else {
				right = append(right, input[midIx+i])
			}
		}
	}

	fmt.Printf("%v - %v - %v\n", left, pivot, right)

	if len(left) > 1 {
		left = Quicksort(left)
	}
	input = append(left, pivot)

	if len(right) > 1 {
		right = Quicksort(right)
	}
	input = append(input, right...)

	return input
}

func main() {
	// El slice que deseamos ordenar
	var unsorted = []int{100, 1, 55, 6, 23, 28, 97, 61, 34, 2, 7, 17, 44, 5, 1, 0}
	fmt.Print("===========================================\n")
	fmt.Printf("ENTRADA: %v\n", unsorted)
	unsorted = Quicksort(unsorted)
	// Imprimimos el resultado final
	fmt.Printf("RESULTADO: %v\n", unsorted)
	fmt.Print("===========================================\n")
}
