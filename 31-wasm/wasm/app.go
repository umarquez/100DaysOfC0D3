// +build js,wasm
package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

const linesNum = 40
const colsNum = 80
const br = "\n"
const fontSize = 10      // px
const lineHeight = 8     // px
const padding = 5        // px
const letterSpacing = -2 // px
const bgColor = "#333"
const fontColor = "#DDD"
const fontFamily = "monospace"
const fontWeight = "bolder"
const textWidth = "min-content"
const containerWidth = "min-content"

var signal = make(chan int)

// SetCSS modifica el estilo del objeto intence de acuerdo con los valores
// almacenados en el map[string]string values, en donde el nombre de la clave es el
// atributo a modificar y el valor es el que se establecerá para dicho atributo
func SetCSS(instance js.Value, values map[string]string) {
	style := instance.Get("style")
	for k, v := range values {
		style.Set(k, v)
	}
}

// RandomContent Genera caracteres '\' y '/' de manera aleatoria, para esto
// divide el espacio disponible (lineas/columnas) en 4 cuadrantes, dentro de los que
// coloca cada caracter a modo de mosaico, reflejando su posición/dirección.
// Retorna el texto y la semilla que lo generó.
func RandomContent(height, width int) (content string, seed int64) {
	// Una semilla diferente en cada ejecución
	seed = rand.Int63()
	altRand := rand.New(rand.NewSource(seed))
	for lines := height / 2; lines > 0; lines-- {
		lineContent := ""
		lineContentMirror := ""
		for cols := width / 2; cols > 0; cols-- {
			if altRand.Float64() > .5 {
				lineContent = `\` + lineContent + `/`
				lineContentMirror = `/` + lineContentMirror + `\`
			} else {
				lineContent = `/` + lineContent + `\`
				lineContentMirror = `\` + lineContentMirror + `/`
			}
		}
		content = lineContent + br + content + lineContentMirror + br
	}

	return
}

func Generate(cols, rows int, target, seed js.Value) {
	content, seedVal := RandomContent(cols, rows)
	target.Set("innerHTML", content)
	seed.Set("value", fmt.Sprintf("%v", seedVal))
}

func main() {
	seed := time.Now().Unix()
	rand.Seed(seed)
	/*
	 * ####################################
	 * # Inicializando el contenedor      #
	 * ####################################
	 */
	document := js.Global().Get("document")
	container := document.Call("getElementById", "wasm_slashes")
	slashes := document.Call("createElement", "p")
	txt := document.Call("createTextNode", "")

	txtInstructions := document.Call("createElement", "p")
	txtInstructions.Set("innerHTML", "Da click en el área para generar un nuevo patrón")

	seedTag := document.Call("createElement", "strong")
	seedTag.Set("innerHTML", "Semilla:")

	txtSeed := document.Call("createElement", "input")
	txtSeed.Set("type", "text")

	slashes.Call("appendChild", txt)
	container.Call("appendChild", slashes)
	container.Call("appendChild", txtInstructions)
	container.Call("appendChild", seedTag)
	container.Call("appendChild", txtSeed)

	/*
	 * ####################################
	 * # Modificando estilos              #
	 * ####################################
	 */
	SetCSS(slashes, map[string]string{
		"line-height":      fmt.Sprintf("%vpx", lineHeight),
		"background-color": bgColor,
		"font-size":        fmt.Sprintf("%vpx", fontSize),
		"color":            fontColor,
		"font-family":      fontFamily,
		"width":            textWidth,
		"padding":          fmt.Sprintf("%vpx", padding),
		"letter-spacing":   fmt.Sprintf("%vpx", letterSpacing),
		"font-weight":      fontWeight,
	})

	SetCSS(container, map[string]string{
		"width":     containerWidth,
		"font-size": fmt.Sprintf("%vpx", 12),
	})

	click := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Generate(colsNum, linesNum, slashes, txtSeed)
		return nil
	})
	slashes.Call("addEventListener", "click", click)

	for {
		<-signal
	}
}
