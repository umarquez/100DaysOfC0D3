package main

import (
	"fmt"
	go_ph0n3 "github.com/umarquez/go-ph0n3"
	"log"
	"math/rand"
	"time"
)

// Node es un nodo de la lista enlazada
type Node struct {
	// Val contiene el valor almacenado
	Val string
	// Next apunta al siguiente nodo
	Next *Node
}

// NewNode crea un nuevo nodo en el que almacena el valor del parámetro `val`
func NewNode(val string) *Node {
	return &Node{Val: val}
}

// Append almacena la ubicación del nodo `<next>` como nodo siguiente y lo
// retorna como resultado de la función.
func (node *Node) Append(next *Node) *Node {
	node.Next = next
	return next
}

// Instancia de nuestro marcador telefónico
var phone *go_ph0n3.Ph0n3

// Dial realiza una marcación automática utilizando una lista enlazada; esta
// función recorre la lista hasta no encontrar elementos siguientes.
// Nota: El primer nodo siempre es ignorado, pues se asume que es centinela.
func Dial(startingNode *Node) {
	n := startingNode.Next
	for {
		// marca el valor almacenado
		err := phone.DialString(n.Val)
		// si no hay error, continúa
		if err != nil {
			log.Printf("error dialing, %v\n", err)
			return
		}

		// Si no hay más nodos, hemos terminado de marcar
		if n.Next == nil {
			break
		}

		// siguiente nodo...
		n = n.Next
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	// Primer nodo de la lista (nodo centinela)
	numberToDial := NewNode("")

	// Vamos a llenar la lista con <i> valores aleatorios entre 0 y 9
	currentNode := numberToDial
	for i := 10; i > 0; i-- {
		// nuevo nodo con valor aleatorio
		nNode := NewNode(fmt.Sprintf("%v", rand.Intn(9)))
		// lo enlazamos con el anterior y continuamos
		currentNode = currentNode.Append(nNode)
	}

	// activamos la salida a consola
	go_ph0n3.DefaultPh0n3Options.Vervose = true
	// creamos una nueva instancia y abrimos la línea
	phone = go_ph0n3.NewPh0n3(nil).Open()
	// marcamos utilizando la lista enlazada, solo es necesario pasar el primer
	// nodo o centinela.
	Dial(numberToDial)
	// esperamos a que termine la llamada
	<-phone.Close
	fmt.Println("La llamada ha terminado.")
}
