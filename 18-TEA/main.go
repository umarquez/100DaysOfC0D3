/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [18/100] TEA
*/
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

/* The constant delta = 31 ( 5 - 1)*2 = , is derived from the golden h 9E3779B9
 * number ratio to ensure that the sub keys are distinct and its precise value
 * has no cryptographic significance.
 *
 * http://www.csshl.net/sites/default/files/downloadable/crypto/TEA_Cryptanalysis_-_VRAndem.pdf
 */
const deltaCrypt uint32 = 0x9E3779B9

// Cantidad de rondas totales
const cipherRounds = 64

// Clave de cifrado a utilizar
const strKey = "ASDFasdfQWERqwer"

// El algoritmo requiere big-endian pues convertiremos [4]byte en DWORD y
// viseversa, tando con la clave como con los datos planos
var bin = binary.BigEndian

// TEAdecrypt decodifica un mensaje cifrado con TEA
func TEAdecrypt(encContent []byte, key []byte) (res []byte, err error) {
	if len(encContent)%8 != 0 {
		return nil, errors.New("content length is incorrect")
	}
	if len(key) < 16 {
		return nil, errors.New("insufficient key lenght")
	}

	// Dividimos la clave en 4 partes de 32 bits
	k0 := bin.Uint32(key[0:])
	k1 := bin.Uint32(key[4:])
	k2 := bin.Uint32(key[8:])
	k3 := bin.Uint32(key[12:])

	// Por cada bloque a encriptar de 8 bytes
	for i := 0; i < len(encContent); i += 8 {
		/*
		 * Por alguna razón no se puede expresar así:
		 * sum :=  deltaCrypt * uint32(cipherRounds/2)
		 *
		 * o así:
		 * sum :=  uint32(deltaCrypt) * uint32(cipherRounds/2)
		 */
		d := uint32(deltaCrypt)
		sum := d * uint32(cipherRounds/2)

		// Obtenemos las 2 partes
		block := encContent[i : i+8]
		v0 := bin.Uint32(block[0:])
		v1 := bin.Uint32(block[4:])

		// Realizamos las rondas a la inversa
		for r := 0; r < cipherRounds/2; r++ {
			// Si en la encripción comenzamos con v0 aquí deberemos comezar con
			// el opuesto, lo que intentamos es recorrer el proceso al revez
			v1 -= ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
			v0 -= ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
			sum -= deltaCrypt
		}

		// volcamos los resultados en un []byte que luego añadiremos al
		// resultado final
		var partial [8]byte
		bin.PutUint32(partial[0:], v0)
		bin.PutUint32(partial[4:], v1)
		res = append(res, partial[:]...)
	}

	return res, nil
}

// TEAencrypt cifra un mensaje utilizando TEA
func TEAencrypt(clearContent []byte, key []byte) (res []byte, err error) {
	if len(key) < 16 {
		return nil, errors.New("insufficient key lenght")
	}

	// si la información no es suficiente, rellenamos con 0x00
	for len(clearContent)%8 != 0 {
		clearContent = append(clearContent, 0x00)
	}

	// Dividimos la clave en 4 partes de 32 bits
	k0 := bin.Uint32(key[0:])
	k1 := bin.Uint32(key[4:])
	k2 := bin.Uint32(key[8:])
	k3 := bin.Uint32(key[12:])

	// Por cada bloque a encriptar de 8 bytes
	for i := 0; i < len(clearContent); i += 8 {
		block := clearContent[i : i+8]
		sum := uint32(0)
		v0 := bin.Uint32(block[0:])
		v1 := bin.Uint32(block[4:])

		// Realizamos las rondas correspondientes
		for r := 0; r < cipherRounds/2; r++ {
			sum += deltaCrypt
			v0 += ((v1 << 4) + k0) ^ (v1 + sum) ^ ((v1 >> 5) + k1)
			v1 += ((v0 << 4) + k2) ^ (v0 + sum) ^ ((v0 >> 5) + k3)
		}

		// volcamos los resultados en un []byte que luego añadiremos al
		// resultado final
		var partial [8]byte
		bin.PutUint32(partial[0:], v0)
		bin.PutUint32(partial[4:], v1)
		res = append(res, partial[:]...)
	}

	return res, nil
}

/*
 * NOTA IMPORTANTE:
 * No es precisamente una práctica segura ejecutar el mismo proceso sobre cada
 * bloque con la misma clave.
 */
func main() {
	var key = []byte(strKey)

	// Texto a encriptar
	txtClear := []byte("Hola Mundo!")
	fmt.Printf("Mensaje original:\n%s\n\n", txtClear)

	// Encriptando
	enc, err := TEAencrypt(txtClear, key)

	// Desencriptando
	clear, err := TEAdecrypt(enc, key)
	if err != nil {
		fmt.Print(err)
	}

	// Resultados:
	fmt.Printf("Mensaje cifrado: \n%s\n\nMensaje decodificado: \n%s", enc, clear)
}
