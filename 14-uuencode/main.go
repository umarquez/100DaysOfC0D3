package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func UUEncode(content []byte, filename string, outfile string) (encoded []byte) {
	var tempByte byte
	var position int

	// Recorremos todos los elementor del slice
	//loop :
	for i, octet := range content {
		// Solo imprimimos los primeros bytes de ejemplo
		if i < 12 {
			fmt.Printf("%b ", octet)
		}
		// Evaluamos que grupo estamos procesando
		position = (i + 3) % 3
		switch position {
		// Primer octeto
		case 0:
			b := octet & 0xFC
			b = b >> 2
			encoded = append(encoded, b+32)

			b = octet & 0x3
			tempByte = b << 4
		// Segundo octeto
		case 1:
			b := octet & 0xF0
			b = b >> 4
			b = b | tempByte
			encoded = append(encoded, b+32)

			tempByte = (octet & 0x0F) << 2
		// Tercer octeto
		case 2:
			b := octet & 0xC0
			b = b >> 6
			b = b | tempByte
			encoded = append(encoded, b+32)

			b = octet & 0x3F
			encoded = append(encoded, b+32)
			tempByte = 0x00

			// Solo imprimimos los primeros bytes de ejemplo
			if i < 12 {
				fmt.Printf(" --> %b\n", encoded[len(encoded)-4:])
			}
			//break loop
		}
	}

	if position < 2 {
		// Agregamos el último dato obtenido en caso necesario
		encoded = append(encoded, tempByte+32)

	}

	// Preparando la salida
	var strOut string
	for {
		// Es la última linea?
		if len(encoded) < 60 { // SI
			// Obtenemos la cantidad de bytes codificados en la línea actual
			lnSize := (len(encoded) / 4) * 3
			if len(encoded)%4 > 0 {
				lnSize += (len(encoded) % 4) - 1
			}

			/*
			 * If the source is not divisible by 3 then the last 4-byte section
			 * will contain padding bytes to make it cleanly divisible. These
			 * bytes are subtracted from the line's <length character> so that
			 * the decoder does not append unwanted characters to the file.
			 *
			 * https://en.wikipedia.org/wiki/Uuencoding#Formatting_mechanism
			 */

			extraChars := ""
			// Si la longitud no es divisible entre 4 agregamos padding bytes
			if len(encoded)%4 > 0 {
				extraChars = strings.Repeat(" ", 4-len(encoded)%4)
			}

			// Almacenamos la linea final
			strOut += string(lnSize+32) + string(encoded) + extraChars + "\n"
			break
		} else { // NO
			// Entonces todas las líneas tienen un tamaño fijo:
			// 45 caracteres planos + 32 que debemos sumar a todas las salidas
			strOut += string(45+32) + string(encoded[:60]) + "\n"
			encoded = encoded[60:]
		}
	}

	/*
	 * Note that 96 ("`" grave accent) is a character that is seen in uuencoded
	 * files but is typically only used to signify a 0-length line, usually at
	 * the end of a file. It will never naturally occur in the actual
	 * converted data since it is outside the range of 32 to 95. The sole
	 * exception to this is that some uuencoding programs use the grave accent
	 * to signify padding bytes instead of a space. However, the character used
	 * for the padding byte is not standardized, so either is a possibility.
	 *
	 * https://en.wikipedia.org/wiki/Uuencoding#Uuencode_table
	 */
	// Nosotros lo agregaremos para poder decodificar nuestros archivos con
	// las aplicaciones típicas como WinRar, 7Zip, etc...
	strOut = strings.Replace(strOut, " ", "`", -1)

	// Añadimos la cabecera del archivo
	strOut = fmt.Sprintf("begin 644 %v\n%v", filename, strOut)
	// Fin del archivo
	strOut += "`\nend"

	// Escribiendo archivo
	err := ioutil.WriteFile("./"+outfile, []byte(strOut), os.ModePerm)
	if err != nil {
		log.Fatalf("error creating UUE file, %v", err)
	}

	return []byte(strOut)
}

func UUDecode(content []byte) (decoded []byte) {
	lines := strings.Split(string(content), "\n")

	// Solo decodificamos las líneas con información, asumimos que la primera
	// es el header y la última y penúltima el cierre del archivo
	for _, uueLine := range lines[1 : len(lines)-2] {
		// Obtenemos el tamaño final de bytes decodificados
		size := byte(uueLine[0]) - 32
		// Slice donde almacenaremos los datos decodificados
		var decLine []byte

		// Restauramos los espacios sustituidos por acentos
		// y omitimos el primer caracter pues es el que indica el tamaño.
		uueLine = strings.Replace(uueLine, "`", " ", -1)

		var cByte byte
		for i, char := range []byte(uueLine[1:]) {
			char = char - 32
			switch (i + 4) % 4 {
			case 0:
				cByte = char << 2
			case 1:
				cByte = cByte | (char >> 4)
				decLine = append(decLine, cByte)
				cByte = char << 4
			case 2:
				cByte = cByte | (char >> 2)
				decLine = append(decLine, cByte)
				cByte = char << 6
			case 3:
				cByte = cByte | char
				decLine = append(decLine, cByte)
			}
		}
		decoded = append(decoded, decLine[:size]...)
	}
	return
}

func main() {
	// - Codificando
	filename := "message.txt"
	content, err := ioutil.ReadFile("./" + filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Primeros octetos:\n")
	fmt.Print("=============================================================\n")
	outfile := "message.uue"
	encoded := UUEncode(content, filename, outfile)

	fmt.Print("\nArchivo UUE:\n")
	fmt.Print("=============================================================\n")
	fmt.Printf("%s\n", encoded)

	// - Decodificando
	decoded := UUDecode(encoded)
	fmt.Print("\nDecodificado:\n")
	fmt.Print("=============================================================\n")
	fmt.Printf("%s\n", decoded)
}
