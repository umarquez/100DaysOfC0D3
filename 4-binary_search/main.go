package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ApiUrl Endpoint de la API de donde obtenemos el texto aleatorio
const ApiUrl = "http://www.randomtext.me/api/gibberish/p-%v/%v"

// CharMatch Los caracteres son iguales
const CharMatch = 0

// CharIsBigger El caracter oculto es mayor
const CharIsBigger = 1

// CharIsSmaller El caracter oculto es menos
const CharIsSmaller = -1

// PsCounter Párrafos a obtener
const PsCounter = 1

// WordsCounter Palabras a obtener
const WordsCounter = 50

var secret *SecretText
var errNotFound = errors.New("not found")

type txtWrap struct {
	Type      string `json:"type"`
	Amount    int    `json:"amount"`
	Number    int    `json:"number,string"`
	NumberMax int    `json:"number_max,string"`
	Format    string `json:"format"`
	Time      string `json:"time"`
	TextOut   string `json:"text_out"`
}

// SecretText Es la interface para interactual con el texto secreto
type SecretText struct {
	text    txtWrap // Texto oculto
	guessed bool    // Ha sido adivinado?
}

// NewSecretText Retorna un nuevo *SecretText con un texto aleatorio
// precargado, obtenido de la API.
func NewSecretText() (*SecretText, error) {
	secret := new(SecretText)
	// Consultando la API
	res, err := http.Get(fmt.Sprintf(ApiUrl, PsCounter, WordsCounter))
	if err != nil {
		return nil, err
	}

	// Decodificando resultado
	jDec := json.NewDecoder(res.Body)
	err = jDec.Decode(&secret.text)
	if err != nil {
		return nil, err
	}
	// Necesitamos limpiar las etiquetas HTML
	secret.text.TextOut = strings.NewReplacer("<p>", "", "</p>", "").Replace(secret.text.TextOut)
	return secret, nil
}

// GuessChar Retorna 0 si el caracter en la posición dada es igual al del
// parámetro sample, -1 si el de la posición dada es menor y 1 si es mayor.
func (secText *SecretText) GuessChar(position int, sample rune) int {
	if rune(secText.text.TextOut[position]) > sample {
		return CharIsBigger // +1
	} else if rune(secText.text.TextOut[position]) < sample {
		return CharIsSmaller // -1
	}
	return CharMatch // 0
}

// Length Devuelve la longitud  en caracteres del texto secreto
func (secText *SecretText) Length() int {
	return len(secText.text.TextOut)
}

// TestFullText Prueba si la face dada coincide con el texto secreto.
// Si esto es cierto, desbloquea el texto plano para se consultado.
func (secText *SecretText) TestFullText(txt string) bool {
	secret.guessed = txt == secText.text.TextOut
	return secret.guessed
}

// Si el texto secreto ya ha sido adivinado, devuelve este en texto plano.
func (secText *SecretText) RevealedText() string {
	if secret.guessed {
		return secret.text.TextOut
	}

	return "" // Aún no ha sido adivinado
}

// ============================================================================
// Búsqueda lineal.
func linearSearch() (steps int, err error) {
	var result string

	// iteramos cada caracter del texto secreto
	for position := 0; position < secret.Length(); position++ {
		var iChar = 0x00
		// Probamos cada caracter a partir del 0x00 y sumando 1 hasta encontrar
		// una coincidencia.
		for secret.GuessChar(position, rune(iChar)) != 0 {
			steps++
			iChar++
		}

		// Concatenamos el caracter encontrado a la cadena resultante
		result += string(iChar)
	}

	// Evaluamos si la cadena resultante coincide con el secreto
	if !secret.TestFullText(result) {
		return -1, errNotFound
	}

	return
}

// Búsqueda binaria.
func binarySearch() (steps int, err error) {
	var result string
	var position = 0
	var char rune

	// Repetimos hasta adivinar el texto o igualar las longitudes.
	for !secret.TestFullText(result) && len(result) <= secret.Length() {
		// Asumimos que los caracteres se encuentran entre ASCII(0) y ASCII(255)
		var min = 0
		var max = 0xFF
		// Estado inicial
		var guessResult = CharIsBigger

		// Repetimos mientras no encontremos coincidencias
		for guessResult != 0 {
			steps++

			char = rune(min + ((max - min) / 2))           // caracter de la mitad del rango
			guessResult = secret.GuessChar(position, char) // Coincide?

			if guessResult == 1 {
				// El caracter oculto es mayor, así que ajustámos el límite
				// inferior a la mitad del rango +1.
				min += (max - min) / 2
				min++
			} else if guessResult == -1 {
				//	El caracter oculdo es menor, así que reducimos el límite
				//	inferior a la mitad del rango
				max -= (max - min) / 2
			}

			// Si terminamos sin encontrar coincidencias.
			if min > max {
				return -1, errNotFound
			}
		}

		// Concatenamos eñ caracter al recultado final
		result += string(char)
		// Siguiente posición
		position++
	}
	return
}

func main() {
	var err error
	// Inicializamos el texto secreto
	secret, err = NewSecretText()
	if err != nil {
		panic(err)
	}
	fmt.Printf("La longitud del texto secreto es de %v caracteres.\n", secret.Length())
	fmt.Println("====================\n")

	// Búsqueda lineal
	steps, err := linearSearch()
	if err != nil {
		fmt.Printf("error en búsqueda lineal, %v\n", err)
	}
	fmt.Printf("La búsqueda lineal de texto requirió %v pasos\n", steps)

	// Búsqueda lineal
	steps, err = binarySearch()
	if err != nil {
		fmt.Printf("error en búsqueda lineal, %v\n", err)
	}
	fmt.Printf("La búsqueda binaria de texto requirió %v pasos\n", steps)

	// Imprimimos el resultado
	fmt.Printf("\nTexto revelado: \n%v\n", secret.RevealedText())
}
