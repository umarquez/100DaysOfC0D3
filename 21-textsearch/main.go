/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [21/100] Búsqueda simple de texto | Simple text search
*/
package main

import (
	"fmt"
	"strings"
)

// Diálogos
const slSpeech = `<!-- start slipsum code -->
My money's in that office, right? If she start giving me some bullshit about it ain't there, and we got to go someplace else and get it, I'm gonna shoot you in the head then and there. Then I'm gonna shoot that bitch in the kneecaps, find out where my goddamn money is. She gonna tell me too. Hey, look at me when I'm talking to you, motherfucker. You listen: we go in there, and that nigga Winston or anybody else is in there, you the first motherfucker to get shot. You understand?
Normally, both your asses would be dead as fucking fried chicken, but you happen to pull this shit while I'm in a transitional period so I don't wanna kill you, I wanna help you. But I can't give you this case, it don't belong to me. Besides, I've already been through too much shit this morning over this case to hand it over to your dumb ass.
Look, just because I don't be givin' no man a foot massage don't make it right for Marsellus to throw Antwone into a glass motherfuckin' house, fuckin' up the way the nigger talks. Motherfucker do that shit to me, he better paralyze my ass, 'cause I'll kill the motherfucker, know what I'm sayin'?
<!-- end slipsum code -->`

// Lista de palabras a censurar
const censoredWordsCSV = "fuck, ass, dead, bitch"

// Separador de csv
const csvSeparator = ", "

// Caracter con el que reemplazaremos las palabras censuradas
const replacer = "🚫"

// SimpleSearchAndReplace busca y reemplaza todas las coincidencias de una
// lista de palabras en un texto dado con los caracteres seleccionados,
// retornando la lista de coincidencias y el texto resultante.
func SimpleSearchAndReplace(content string, search ...string) (resultString string, positions map[string][]int) {
	positions = make(map[string][]int)

	// Recorremos cada caracter del texto
	for i := 0; i < len(content); i++ {
		// Evaluamos cada palabra de la lista
		for _, cWrd := range search {
			// Si los caracteres restantes no alcanzan, omitimos la palabra
			if i+len(cWrd) >= len(content) {
				continue
			}

			//fmt.Print(content[i:i+len(cWrd)]+"\n")

			// Obtenemos una sección del tamaño de la palabra y lo comparamos
			// para saber si coincide con esta
			if content[i:i+len(cWrd)] == cWrd {
				// Aquí reemplazamos la coincidencia con los caracteres
				content = content[:i] + strings.Repeat(replacer, len(cWrd)) + content[i+len(cWrd):]
				// Y guardamos la posición
				positions[cWrd] = append(positions[cWrd], i)
			}
		}
	}
	return content, positions
}

func main() {
	// Separamos la lista de palabras en un slice y lo pasamos como parámetro
	// a la función, además del texto ogiginal.
	censoredText, censoredPositions := SimpleSearchAndReplace(slSpeech, strings.Split(censoredWordsCSV, csvSeparator)...)

	// Imprimimos los resultados
	fmt.Printf("%s\n\n%s\n\n", slSpeech, censoredText)
	for k, v := range censoredPositions {
		fmt.Printf("%s: %v\n", k, v)
	}
}
