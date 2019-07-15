package main

import (
	"fmt"
	"image"
	"image/color"
)

// Nodos del árbol binario
type treeNode struct {
	name        string
	value       int
	left, right *treeNode
}

// Coloca un nodo nuevo bajo el nivel correspondiente, recorriendo cada nodo
// hijo y comparando valores hasta encontrar la posición final.
func (node *treeNode) Append(colorName string, counter int) {
	// Izquierda?
	if counter < node.value {
		if node.left == nil {
			node.left = new(treeNode)
			node.left.value = counter
			node.left.name = colorName
			return
		}
		node.left.Append(colorName, counter)
	} else { // Derecha?
		if node.right == nil {
			node.right = new(treeNode)
			node.right.value = counter
			node.right.name = colorName
			return
		}
		node.right.Append(colorName, counter)
	}
}

// PixelCounter es nuestro contador de pixeles/colores
type PixelCounter struct {
	img        image.Image    // Imagen a procesar
	htColors   map[string]int // Hash Table para almacenar los colores obtenidos
	treeColors *treeNode      // Primer nodo del árbol
}

// NewPixelCounter retorna una nueva instancia del contador
func NewPixelCounter(img image.Image) *PixelCounter {
	pxc := new(PixelCounter)
	pxc.img = img
	pxc.htColors = make(map[string]int) //  Inicializamos la tabla

	pxc.count() // Iniciamos el conteo
	return pxc
}

// hashColor es la función hash, en este caso el RGB en hexadecimal del color.
func (pxc *PixelCounter) hashColor(c color.Color) string {
	rgba := c.(color.RGBA)
	return fmt.Sprintf("%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

// count recorre los pixeles de la imagen del final al inicio, llenando la
// tabla con los colores y la frecuencia de cada uno.
func (pxc *PixelCounter) count() {
	for x := pxc.img.Bounds().Dx() - 1; x >= 0; x-- {
		for y := pxc.img.Bounds().Dy() - 1; y >= 0; y-- {
			// Convertimos a RGBA
			pxColor := color.RGBAModel.Convert(pxc.img.At(x, y))

			// Incrementamos el valor en la tabla
			// Nota: No es necesario inicializar el valor en la tabla pues los
			// valores iniciales son ceros "0"
			key := pxc.hashColor(pxColor)

			// Incrementamos en 1
			pxc.htColors[key]++
		}
	}

	//fmt.Printf("%#v\n", pxc.htColors)
}

// inorderedWalk recorre el árbol binario, recuperando cada clave en orden.
func inorderedWalk(node *treeNode, keys *[]string) {
	if node.left != nil {
		inorderedWalk(node.left, keys)
	}

	*keys = append(*keys, node.name)

	if node.right != nil {
		inorderedWalk(node.right, keys)
	}
}

// Sort es nuestra función de ordenamiento
func (pxc *PixelCounter) Sort() []string {
	// Si el arbol se encuentra vacío, lo conformamos
	if pxc.treeColors == nil {
		// Recorremos la tabla recuperando cada valor.
		for k, v := range pxc.htColors {
			// Solo la primera vuelta
			if pxc.treeColors == nil {
				pxc.treeColors = new(treeNode)
				pxc.treeColors.name = k
				pxc.treeColors.value = v
				continue
			}

			// A partir del primer nodo, podemos pasar los valores al primero,
			// el cual automáticamenete lo ordenará en cascada
			pxc.treeColors.Append(k, v)
		}
	}

	// Si árbol ya cuenta con valores almacenados, solo lo recorremos para
	// obtener las claves en orden
	var sortedKeys []string
	inorderedWalk(pxc.treeColors, &sortedKeys)
	return sortedKeys
}

// Get retorna la frecuencia de un color
func (pxc *PixelCounter) Get(colorName string) int {
	return pxc.htColors[colorName]
}
