package main

import (
	"fmt"
	"golang.org/x/image/tiff"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

// Utilizaremos una lista enlazada para almacenar la secuencia de bits a codificar
// BitNode será el nodo de nuestra lista
type BitNode struct {
	value bool
	next  *BitNode
}

// BytesToBits Almacena los bits de un []byte dentro de una lista enlazada,
// retorna la lista y la longitud de esta.
func BytesToBits(content []byte) (result *BitNode, length int) {
	var firstNode = new(BitNode)

	if len(content) < 1 {
		return nil, 0
	}
	current := firstNode
	for _, b := range content {
		current.value = (b & 0x01) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x02) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x04) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x08) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x10) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x20) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x40) != 0
		current.next = new(BitNode)
		current = current.next
		current.value = (b & 0x80) != 0
		current.next = new(BitNode)
		current = current.next
		length += 8
	}

	return firstNode, length
}

// Encode almacena el contenido del mensaje dado dentro de una imagen cuadrada
// del tamaño mínimo necesario (o de 64x64px si fuera menor a este tamaño), esta es obtenida de manera aleatorio a partir
// de una llamada a un API.
func Encode(msj string, outFile string, saveOriginal bool) {
	// Tamaño minimo necesario
	minSize := int(math.Ceil(math.Sqrt(float64(len(msj)*8) / 3)))

	// El tamaño mínimo de la imagen será de 64x64px
	if minSize < 64 {
		minSize = 64
	}

	// Obtenemos una imagen jpg aleatorio
	img, err := RandomImage(image.Pt(minSize, minSize))
	if err != nil {
		panic(err)
	}

	// Agregamos 0x000000 al final del mensaje para saber cuando termina (EOF)
	content := append([]byte(msj), 0x00, 0x00, 0x00)
	// Almacenamos cada bit en una lista enlazada
	cNode, _ := BytesToBits(content)

	// Recorremos todos los pixeles de la imagen
	outImg := image.NewRGBA(image.Rect(0, 0, minSize, minSize))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// convertimos el color de cada pixen a RGB
			pxcolor := color.NRGBAModel.Convert(img.At(x, y))
			r, g, b, a := pxcolor.RGBA()
			var nr, ng, nb uint8
			nr = uint8(r)
			ng = uint8(g)
			nb = uint8(b)

			// Establecemos el bit menos significativo de cada octeto
			/*  _____
			 * |  _  |
			 * |    _|
			 * | |\ \
			 * |_| \_\
			 */
			if cNode != nil {
				if cNode.value {
					nr = nr | 0x01
				} else {
					nr = nr & 0xFE
				}
				cNode = cNode.next
			}

			/*  _____
			 * |  ___|
			 * | |  __
			 * | |_|_ |
			 * |______|
			 */
			if cNode != nil {
				if cNode.value {
					ng = ng | 0x01
				} else {
					ng = ng & 0xFE
				}
				cNode = cNode.next
			}

			/*  _____
			 * |  _  |
			 * |    _|
			 * |  _  |
			 * |_____|
			 */
			if cNode != nil {
				if cNode.value {
					nb = nb | 0x01
				} else {
					nb = nb & 0xFE
				}
				cNode = cNode.next
			}

			// Conformamos el pixel con los nuevos valores y lo guardamos en
			// la imagen de salida.
			pxcolor = color.NRGBA{R: nr, G: ng, B: nb, A: uint8(a)}
			outImg.Set(x, y, pxcolor)
		}
	}

	// Generamos el archivo
	f, err := os.Create(outFile + ".tiff")
	if err != nil {
		log.Fatal(err)
	}

	// Lo cerramos al terminar
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("error closing tiff file, %v", err)
		}
	}()

	// Escribimos la imagen resultante en el archivo creado
	err = tiff.Encode(f, outImg, &tiff.Options{
		// Deflate es un algoritmo de compresión sin pérdidas.
		Compression: tiff.Deflate,
	})
	if err != nil {
		log.Fatal(err)
	}

	// ¿almacenamos la imagen original?
	if saveOriginal {
		// Generamos el archivo
		fo, err := os.Create(outFile + "_original.tiff")
		if err != nil {
			log.Fatal(err)
		}

		// Lo cerramos al terminar
		defer func() {
			err := fo.Close()
			if err != nil {
				log.Printf("error closing tiff file, %v", err)
			}
		}()

		// Escribimos la imagen resultante en el archivo creado
		err = tiff.Encode(fo, outImg, &tiff.Options{
			// Deflate es un algoritmo de compresión sin pérdidas.
			Compression: tiff.Deflate,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Decode lee el contenido de cada pixel de una imagen, intentando recuperar
// un mensaje oculto dentro de ella.
func Decode(img image.Image) []byte {
	// slices de bits
	var vals []bool

	// recorremos los pixeles de la imagen
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			pxcolor := color.NRGBAModel.Convert(img.At(x, y))
			r, g, b, _ := pxcolor.RGBA()
			// recuperando bits menos significativos de cada octeto
			vals = append(vals, uint8(r)&0x01 != 0)
			vals = append(vals, uint8(g)&0x01 != 0)
			vals = append(vals, uint8(b)&0x01 != 0)
		}
	}

	// contador de valores de final de mensaje
	var eof int
	// bytes resultantes
	var res []byte
	// byte actualmente procesado
	var cByte byte

	// iteramos cada bit recuperado
	for i, v := range vals {
		// posición actual
		var pos = uint(i+8) % 8

		// si es el inicio del byte
		if pos == 0 {
			// empezamos con un octeto vacío
			cByte = byte(0x00)
		}

		// establecemos la posición actual en 1 si es necesario
		if v {
			cByte = cByte | (1 << pos)
		}
		// si es el final del byte
		if pos == 7 {
			// y el byte es un 0x00, incrementamos el contador de final de
			// mensaje, de lo contrario lo establecemos en cero
			if cByte == 0x00 {
				eof++
			} else {
				eof = 0
			}

			// almacenamos el byte obtenido
			res = append(res, cByte)
		}

		// si hemos obtenido la secuencia necesario para considerar el final
		// del mensaje salimos del loop.
		if eof > 2 {
			break
		}
	}

	if eof > 2 {
		res = res[:len(res)-eof]
	}
	return res
}

func main() {
	// Leemos el mensaje original desde el archivo
	content, err := ioutil.ReadFile("./message.txt.uue")
	if err != nil {
		log.Fatal(err)
	}
	var secretMessage = string(content)

	// Nombre aleatorio de acuerdo a la fecha/hora
	imgFileName := fmt.Sprintf("./%v", time.Now().Unix)
	// Codificando mensaje dentro de una imagen aleatoria
	Encode(secretMessage, imgFileName, true)

	// Leyendo contnido de la imagen TIFF
	img, err := LocalTiffImage(imgFileName + ".tiff")
	if err != nil {
		log.Fatal(err)
	}

	// Decodificando imagen para recuperar mensaje
	r := Decode(img)

	// ¿Mensaje original?
	fmt.Printf("%v\n", string(r))
}
