package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Node es un nodo del árbol
type Node struct {
	// Val contiene el valor almacenado
	Val int
	// Izquierdo y derecho
	Left, Right *Node
}

// Tree es el árbol binario en el que almacenaremos y procesaremos los datos
type Tree struct {
	firstNode *Node // Nodo inicial
	sorted    []int // Resultado del ordeamiento
	// Steps contiene la cantidad de pasos que tomó el último ordenamiento.
	Steps int
}

// inOrderWalk recorre el árbol en inorden, agregando los valores en el slice
// ordenado, además incrementa el contador de pasos.
func (t *Tree) inOrderWalk(node *Node) {
	t.Steps++
	// 1. rama izquierda
	if node.Left != nil {
		t.inOrderWalk(node.Left)
	}
	// 2. Valor actual
	t.sorted = append(t.sorted, node.Val)

	// 3. Rama derecha
	if node.Right != nil {
		t.inOrderWalk(node.Right)
	}
}

// Add agrega un valor al árbol binario.
func (t *Tree) Add(val int) {
	// Si el árbol está vacío, este será el primer nodo.
	if t.firstNode == nil {
		t.firstNode = new(Node)
		t.firstNode.Val = val
		return
	}

	// Si el árbol contiene nodos, los recorremos hasta encontrar uno vacío
	current := t.firstNode
	for {
		// ¿Corresponde a la rama izquierda?
		if val < current.Val {
			// ¿Y la rama izquierda está vacía?
			if current.Left == nil {
				// Almacenamos el valor
				current.Left = new(Node)
				current.Left.Val = val
				break
			}
			// Continuamos con la rama izquierda
			current = current.Left
		} else { // Corresponde a la derecha
			// ¿Y la rama Derecha está vacía?
			if current.Right == nil {
				// Almacenamos el valor
				current.Right = new(Node)
				current.Right.Val = val
				break
			}
			// Continuamos con la rama derecha
			current = current.Right
		}
	}
}

// Sort devuelve un []int con los valores ordenado.
func (t *Tree) Sort() []int {
	t.Steps = 0
	t.sorted = []int{}
	t.inOrderWalk(t.firstNode)
	return t.sorted
}

func main() {
	// Aleatorio cada vez
	rand.Seed(time.Now().Unix())

	fmt.Println("========================================")
	// Mi madre y yo lo plantamos
	// en el límite del patio
	// donde termina la casa...
	//          - Alberto Cortez.
	t := new(Tree)

	// Agregamos <i> valores aleatorios 0-200
	for i := 10; i > 0; i-- {
		newVal := rand.Intn(200)
		fmt.Printf("Agregando [%v]\t a la lista\n", newVal)
		t.Add(newVal)
	}
	fmt.Println("========================================")
	// Ordenamos e imprimimos los resultados.
	fmt.Printf("- RESULTADO: %v\n", t.Sort())
	fmt.Printf("- PASOS: %v\n", t.Steps)
	fmt.Print("- ÁRBOL:\n")
	jEnc := json.NewEncoder(os.Stdout)
	jEnc.SetIndent("", "  ")
	_ = jEnc.Encode(t.firstNode)
}
