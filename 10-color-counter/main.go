/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [10/100] Primer proyecto: Cuenta colores
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
)

const imgW, imgH = 300, 300

var imgSrc string
var outFile string

func init() {
	flag.StringVar(&imgSrc, "img", "", "-img ./example.jpg")
	flag.StringVar(&outFile, "out", "", "-out ./example.json")

	flag.Parse()
}

func main() {
	log.Println("- Obteniendo imágen")
	var img image.Image
	var err error

	// Si el parámetro img no se encuentra
	if imgSrc == "" {
		// Obtenemos una imagen aleatorio de <imgW>px x <imgH>px
		img, err = RandomImage(image.Point{X: imgW, Y: imgH})
	} else {
		// De lo contrario, intentamos leer la imagen
		img, err = LocalImage(imgSrc)
	}

	// Si no pudimos obtener alguna de las dos opciones no podremos continuar
	if err != nil {
		log.Fatal(err)
	}

	// - En este punto debemos contar con una imagen con la que podemos seguir
	// trabajando, iniciamos a medir el tiempo.
	startingTime := time.Now()
	log.Println("- Iniciando conteo de colores | Hash Table")
	pxCounter := NewPixelCounter(img)

	log.Println("- Iniciamos ordenamiento por frecuencia | Binary Tree Sort")
	sortedKeys := pxCounter.Sort()

	log.Println("- Generando reporte y formateando resultados")
	result := new(processResults)
	for _, key := range sortedKeys {
		result.Colors = append(result.Colors, colorTotal{
			RGB:   key,
			Total: pxCounter.Get(key),
		})
	}

	result.Width = img.Bounds().Dx()
	result.Height = img.Bounds().Dy()
	result.TotalTime = time.Since(startingTime).Seconds()

	//fmt.Printf("%#v\n", result)

	log.Println("- Terminando, escribiendo archivos.")
	if imgSrc == "" {
		imgSrc = fmt.Sprintf("./%v.jpg", time.Now().Unix())
		f, err := os.Create(imgSrc)
		if err != nil {
			log.Printf("[Error] - error almacenando imagen, %v", err)
		}

		defer func() {
			_ = f.Close()
		}()

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Printf("[Error] - error almacenando imagen, %v", err)
		}

		outFile = imgSrc + ".json"
	}

	var jEnc *json.Encoder
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal(fmt.Errorf("[Error] - error creado archivo (%v), %v", err))
		}

		jEnc = json.NewEncoder(f)
	} else {
		jEnc = json.NewEncoder(os.Stdout)
	}

	jEnc.SetIndent("", "  ")

	err = jEnc.Encode(result)
	if err != nil {
		log.Printf("[Error] - error codificando salida de datos., %v", err)
	}

	os.Exit(0) // Terminamos sin errores.
}
