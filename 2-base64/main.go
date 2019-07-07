package main

import (
	"errors"
	"fmt"
	"strings"
)

const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func base64Encode(content []byte) (encoded []byte) {
	var tempByte byte
	var position int
	var ixs []uint

	// Recorremos todos los elementor del slice
	for i, octet := range content {
		// Evaluamos que grupo estamos procesando
		position = (i + 3) % 3
		switch position {
		// Primer octeto
		case 0:
			b := octet & 0xFC
			b = b >> 2
			encoded = append(encoded, table[uint(b)])
			ixs = append(ixs, uint(b))

			b = octet & 0x3
			b = b << 4
			tempByte = b
		// Segundo octeto
		case 1:
			b := octet & 0xF0
			b = b >> 4
			b = b | tempByte
			encoded = append(encoded, table[uint(b)])
			ixs = append(ixs, uint(b))

			b = octet & 0x0F
			b = b << 2
			tempByte = b
		// Tercer octeto
		case 2:
			b := octet & 0xC0
			b = b >> 6
			b = b | tempByte
			encoded = append(encoded, table[uint(b)])
			ixs = append(ixs, uint(b))

			b = octet & 0x3F
			encoded = append(encoded, table[uint(b)])
			ixs = append(ixs, uint(b))
			tempByte = 0
		}
	}

	if position < 2 {
		// Agregamos el último dato obtenido en caso necesario
		encoded = append(encoded, table[uint(tempByte)])
		ixs = append(ixs, uint(tempByte))
		// Esta condición también define el primer "="
		encoded = append(encoded, []byte("=")[0])
	}

	// Agregamos el segundo "=" si fuera el caso
	if position == 0 {
		encoded = append(encoded, []byte("=")[0])
	}

	//fmt.Printf("índices: %v\n", ixs)
	//fmt.Printf("binario: %b\n", ixs)
	return encoded
}

func base64Decode(content []byte) (decoded []byte) {
	var tempByte byte
	var position int
	var ixs []uint

	// Iteramos sobre el texto codificado
	for i, item := range content {
		// Ignoramos los "=" adicionales
		if string(item) == "=" {
			break
		}

		// recuperamos la posición del caracter actual
		ix := strings.Index(table, string(item))
		if ix < 0 {
			panic(errors.New("caracter no permitido"))
		}

		b64 := byte(ix)
		ixs = append(ixs, uint(b64))

		// obtenemos la variante del octeto a procesar
		position = (i + 4) % 4

		switch position {
		case 0:
			// - 6 bits +significativos del primer octeto
			tempByte = b64 << 2
		case 1:
			// - Últimos 2 bits del primer octeto
			// - Primero 4 bits del segundo octeto
			b1 := b64 >> 4
			tempByte = tempByte | b1
			decoded = append(decoded, tempByte)

			tempByte = (b64 & 0x0F) << 4
		case 2:
			// - Últimos 4 bits del tercer octeto
			// - Primero 2 bits del cuarto octeto
			b1 := (b64 >> 2) & 0x0F
			tempByte = tempByte | b1

			decoded = append(decoded, tempByte)
			b2 := b64 & 0x03
			tempByte = b2 << 6
		case 3:
			// - Últimos 6 bits del cuarto octeto
			tempByte = tempByte | b64
			decoded = append(decoded, tempByte)
		}
	}

	// Agregamos el último valor en caso necesario
	if position < 3 {
		decoded = append(decoded, tempByte)
	}
	return
}

func main() {
	// - Texto plano
	testContent := []byte(`Hola Mundo!`)

	// - Codificando
	encoded := base64Encode(testContent)
	fmt.Printf("Base64:\n%s\n\n", encoded)

	// - Decodificando
	decoded := base64Decode(encoded)
	fmt.Printf("Decodificado:\n%s\n", decoded)

	fmt.Print("\n\n--------------------------------------------------\n\n")

	// - Texto plano
	testContent = []byte(`	   _ _
	 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
	|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
	|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
	  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
	                                      |___|`)

	// - Codificando
	encoded = base64Encode(testContent)
	fmt.Printf("Base64:\n%s\n\n", encoded)

	// - Decodificando
	decoded = base64Decode(encoded)
	fmt.Printf("Decodificado:\n%s\n", decoded)
}
