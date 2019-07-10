package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
)

// ApiUrl Endpoint de la API de donde obtenemos el texto aleatorio
const ApiUrl = "http://www.randomtext.me/api/gibberish/p-%v/%v"
const PsCounter = 5
const WordsCounter = 1000

type RandomText struct {
	Type      string `json:"type"`
	Amount    int    `json:"amount"`
	Number    int    `json:"number,string"`
	NumberMax int    `json:"number_max,string"`
	Format    string `json:"format"`
	Time      string `json:"time"`
	TextOut   string `json:"text_out"`
}

func NewRandomText() *RandomText {
	txt := new(RandomText)

	// Consultando la API
	res, err := http.Get(fmt.Sprintf(ApiUrl, PsCounter, WordsCounter))
	if err != nil {
		log.Fatal(err)
		//return nil
	}

	// Decodificando resultado
	jDec := json.NewDecoder(res.Body)
	err = jDec.Decode(&txt)
	if err != nil {
		log.Fatal(err)
		//return nil
	}

	// Necesitamos limpiar las etiquetas HTML
	txt.TextOut = strings.NewReplacer("<p>", "", "</p>", "").Replace(txt.TextOut)

	return txt
}

// Hash Es la función que generará las claves de la tabla, en este caso, debido
// a que nuestro objetivo es medir la frecuencia de las palabras, esta deberá
// retornar una versión optimizada del texto de entrada.
func Hash(input string) string {
	// Limpiamos el texto y lo conventimos a minúsculas
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	// Tabla Hash:
	var wordsCounter = make(map[string]int)
	var keys = []string{}

	// Texto a analizar
	rndText := NewRandomText()

	// Partiendo texto en "palabras" y procesando cada una
	words := strings.Split(rndText.TextOut, " ")
	for _, word := range words {
		// Obteniendo hash
		hWord := Hash(word)
		wordsCounter[hWord]++
	}

	// Obtenemos un slice de las claves para ordenarlas posteriormente
	for k := range wordsCounter {
		keys = append(keys, k)
	}

	fmt.Printf("- Texto:\n%v\n\n\n", rndText.TextOut)
	fmt.Printf("- Palabras totales:\n%v\n\n", len(words))

	// Vamos a ordenar los resultados de mayor a menor
	sort.Slice(keys, func(i, j int) bool {
		return wordsCounter[keys[i]] > wordsCounter[keys[j]]
	})

	// Imprimimos el TopN de resultados
	for i, k := range keys {
		if i >= 20 {
			break
		}

		pcent := math.Floor((float64(wordsCounter[k]*100)/float64(len(words)))*100) / 100
		fmt.Printf("%v\t %v [%v] = %v%%\n", strings.Repeat("█", int(pcent*10)/3), k, wordsCounter[k], pcent)
	}
}
