/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [23/100] Compresión LZ77 | LZ77 Compression
*/
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"
)

// Tamaño máximo de la ventana
const wndSize = 255

// Longitud mínima de la coincidencia para que sea almacenada
const minMatchSize = 6

// Longitud máxima de la coincidencia (0xFF)
const maxMatchSize = 255

// Prefijo de bloque comprimido
const startBlock = "~"

// lz77Compress Utiliza la lógica del algoritmo LZ77 para comprimir una cadena
// de caracteres utilizando bloques de 5 caracteres para codificar el
// desplazamiento y la longitud de cada patrón.
func lz77Compress(content string) string {
	// Los primeros n caracteres del texto original
	var compressed = content[:minMatchSize]

	// Recorremos cada caracter
nxtChar:
	for ix := minMatchSize; ix <= len(content)-1; ix++ {
		cChar := string(content[ix])

		// Recorremos la ventana hacia atrán en búsqueda de un patrón
		for bkIx := 1; bkIx <= wndSize && ix-int(bkIx) > 0; bkIx++ {
			backPos := ix - bkIx

			// Si encontramos una coincidencia
			if cChar == string(content[backPos]) {
				// Continuamos comparando el resto de caracteres
				matchSize := 1
				for backPos+matchSize < ix && ix+matchSize < len(content) && content[ix+matchSize] == content[backPos+matchSize] {
					matchSize++
				}

				// Evaluamos la lingitud del patrón encontrado y si se
				// encuentra dentro del rango funcional escribimos un bloque
				if matchSize > minMatchSize && matchSize < maxMatchSize {
					// Convertimos el offset del patrón y la longitud al string
					// hexadecimal de su valor
					strOffset := strconv.FormatInt(int64(bkIx), 16)
					if len(strOffset) < 2 {
						strOffset = "0" + strOffset
					}
					strSize := strconv.FormatInt(int64(matchSize), 16)
					if len(strSize) < 2 {
						strSize = "0" + strSize
					}

					// Imprime cada bloque
					/*fmt.Printf("%v\t - \t\t%s, %s\n", ix, strOffset, strSize)
					fmt.Printf("\t\t\t%v, %v, %v, %v\n", bkIx, matchSize, ix - bkIx, ix - bkIx + matchSize)
					fmt.Printf("%s\n", content[backPos:backPos+matchSize])*/

					// Añadimos el caracter que señala un incio de bloque
					// seguido del offset y el tamaño
					compressed += startBlock +
						strOffset +
						strSize
					// incrementamos la posición de acuerdo al tamaño del
					// patrón y continuamos con la siguiente posición
					ix += int(matchSize) - 1
					continue nxtChar
				}
			}
		}
		// Si no se encontraran coincidencia o no fueran funcionales, se
		// almacena el texto plano
		compressed += string(content[ix])
	}

	return compressed
}

func lz77Decompress(compressed string) (decompressed string) {
	var offset, size int64
	var err error

	// Recorremos el texto comprimido
	for ix := 0; ix < len(compressed); ix++ {
		// Si encontramos un inicio de bloque
		if string(compressed[ix]) == startBlock {
			// Recuperamos el offset y el tamaño
			bOff := compressed[ix+1 : ix+3]
			bSize := compressed[ix+3 : ix+5]

			// Imprime la información del bloque
			//fmt.Printf("%v\t - \t%s, %s\n", len(decompressed), bOff, bSize)

			// Decodificamos de hex a int
			offset, err = strconv.ParseInt(bOff, 16, 16)
			if err != nil {
				log.Fatal(err)
			}
			size, err = strconv.ParseInt(bSize, 16, 16)
			if err != nil {
				log.Fatal(err)
			}

			// Obtenemos el inicio y fin del patrón
			start := len(decompressed) - int(offset)
			end := start + int(size)

			// Imprimimos la información del patrón
			/*fmt.Printf("\t\t\t%v, %v, %v, %v\n", offset, size, start, end)
			fmt.Printf("%s\n", decompressed[start:end])*/

			// Extraemos el bloque de texto
			strChunck := decompressed[start:end]
			// Y lo concatenamos al resultado
			decompressed += strChunck
			// Continuamos con el siguiente caracter
			ix += 4
		} else {
			// Recuperamos el texto no comprimido
			decompressed += string(compressed[ix])
		}
	}
	return
}

func main() {
	var startTime time.Time
	var secDuration float64
	txtFilePath := "./plaintext.txt"
	//compFilePath := "./plaintext.lz"

	// Leemos el contenido del archivo
	content, err := ioutil.ReadFile(txtFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Codificamos para tener un rango seguro de caracteres.
	/*
	 * NOTA: Esto puede incrementar el tamaño del texto original lo cual es
	 * contraproducente si lo que queremos es reducir el tamaño del contenido.
	 * En este caso es necesario para asegurarnos que el caracter que indicará
	 * el inicio de cada bloque no se encuentre en el texto original y pueda
	 * alterar el comportamiento del descompresor.
	 *
	 * En un caso real, este indicador sería un bit a modo de bandera que
	 * indicaría si la siguiente secuencia de bits codifica un bloque de datos
	 * o bytes planos.
	 */
	content = []byte(url.QueryEscape(string(content)))

	// Medimos el tiempo que tome comprimir
	startTime = time.Now()
	// Shrink it baby!
	compressed := lz77Compress(string(content))
	secDuration = time.Since(startTime).Seconds()

	// Imprimimos resultados
	fmt.Printf("Original: %v chars\n", len(content))
	hashOriginal := sha256.Sum256(content)
	fmt.Printf("Original (hash):\n%s\n\n", base64.RawURLEncoding.EncodeToString(hashOriginal[:]))

	fmt.Printf("Compressed: %v chars\n", len(compressed))
	hashComp := sha256.Sum256([]byte(compressed))
	fmt.Printf("Compressed (hash):\n%s\n", base64.RawURLEncoding.EncodeToString(hashComp[:]))
	fmt.Printf("Compressed proportion: %v percent\n", 100-((len(compressed)*100)/len(content)))
	fmt.Printf("Compression duration: %vsec.\n\n", secDuration)

	// Escribe un archivo de salida, se debe descomentar la línea 146:
	// compFilePath := "./plaintext.lz"
	// para definir el nombre del archivo
	/*err = ioutil.WriteFile(compFilePath, compressed, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}*/

	// Ahora medimos el tiempo que nos toma descomprimir la información
	startTime = time.Now()
	decompressed := lz77Decompress(compressed)
	secDuration = time.Since(startTime).Seconds()

	// Imprimimos los resultados
	fmt.Printf("Decompression duration: %vsec.\n", secDuration)
	fmt.Printf("Decompressed: %v chars\n", len(decompressed))
	// Si el hash del texto original y este coinciden, significa que el
	// cotnenido no ha sido modificado por el proceso de compresión, es decir,
	// que ha sido revertido de manera exitosa.
	hashDecomp := sha256.Sum256([]byte(decompressed))
	fmt.Printf("Decompressed (hash):\n%s\n", base64.RawURLEncoding.EncodeToString(hashDecomp[:]))
	//fmt.Printf("Decompressed:\n%s \n", decompressed)

	// ¿Quieres saber como se ve el texto comprimido?
	//fmt.Printf("Compressed: \n\n%s\n", compressed)
}
