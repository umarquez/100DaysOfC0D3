/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [20/100] Segundo proyecto: Repintado
 */
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Random img API
const imgApiUrl = "https://source.unsplash.com/%vx%v"

// Image size
const imgDefaultWidth = 1920
const imgDefaultHeight = 1080

// Image ext
const jpgExt = ".jpg"
const outDir = "out"

// Source files paths
var imgSrc, palettes *string

// Color's table
var colorsList map[color.Color]color.Color

func downloadRandomImage() (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf(imgApiUrl, imgDefaultWidth, imgDefaultHeight))
	if err != nil {
		err = fmt.Errorf("error tratando de descargar la imágen aleatorioa, %v\n", err)
		return nil, err
	}

	dir := "./" + outDir
	*imgSrc = fmt.Sprintf("%v/%v%v", dir, time.Now().Unix(), jpgExt)

	if info, err := os.Stat(dir); os.IsNotExist(err) || !info.IsDir() {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("error tratando de crear el directorio de trabajo [%v], %v\n", dir, err)
			return nil, err
		}
	}

	f, err := os.Create(*imgSrc)
	if err != nil {
		err = fmt.Errorf("error tratando de crear el archivo origen [%v], %v\n", imgSrc, err)
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error leyendo el contenido de la imagen descargada, %v\n", err)
		return nil, err
	}

	_, err = f.Write(content)
	if err != nil {
		err = fmt.Errorf("error escribiendo el archivo origen, %v\n", err)
		return nil, err
	}

	return bytes.NewBuffer(content), nil
}

func fetchImage() (image.Image, error) {
	var imgReader io.Reader
	var err error

	if *imgSrc == "" {
		imgReader, err = downloadRandomImage()
		if err != nil {
			return nil, err
		}
	} else {
		content, err := ioutil.ReadFile(*imgSrc)
		if err != nil {
			err = fmt.Errorf("error leyendo el contenido del archivo, %v", err)
			return nil, err
		}
		imgReader = bytes.NewBuffer(content)
	}

	img, err := jpeg.Decode(imgReader)
	if err != nil {
		err = fmt.Errorf("error leyendo el contenido del archivo, %v", err)
		return nil, err
	}

	return img, nil
}

// loadPalette convierte un arvhivo .hex en un slice de colores
func loadPalette(paletteFile string) (palette []color.Color, err error) {
	// leemos el contenido del archivo
	paletteContent, err := ioutil.ReadFile(paletteFile)
	if err != nil {
		return nil, err
	}

	// separamo cada línea del archivo
	for _, line := range strings.Split(string(paletteContent), "\r\n") {
		// Si la línea está vacía la omitimos
		if strings.TrimSpace(line) == "" {
			continue
		}

		// obtenemos el valor de la línea y lo convertimos a entero
		n, err := strconv.ParseInt("0x"+line, 0, 32)
		if err != nil {
			return nil, err
		}

		// Obtenemos el valor de cada color
		rgb := make([]byte, 4)
		binary.BigEndian.PutUint32(rgb, uint32(n))

		palette = append(palette, color.NRGBA{R: uint8(rgb[1]), G: uint8(rgb[2]), B: uint8(rgb[3]), A: 0xFF})
	}

	// Si no se obtuvieron valores válidos se genera un error.
	if len(palette) == 0 {
		return nil, errors.New("empty palette file")
	}
	return
}

/*
 * Pitágoras de nuevo:
 * Suponemos que cada color es un punto en un espacio 3D determinado por los
 * los valores RGB, vamos a calcular la distancia entre el color dado y todos
 * los colores de la paleta hasta en contrar el más parecido, que en este caso
 * será la menor distancia.
 */
// calcColor calcula el vecino más cercano de la paleta al color seleccionado
func calcColor(colorIn color.Color, palette []color.Color) (colorOut color.Color) {
	oR, oG, oB, _ := colorIn.RGBA()
	minDist := math.Sqrt(1E12) // Max distance possible
	for _, c := range palette {
		pR, pG, pB, _ := c.RGBA()
		dR, dG, dB := oR-pR, oG-pG, oB-pB
		dist := math.Sqrt(float64(dR*dR) + float64(dG*dG) + float64(dB*dB))
		if dist < minDist {
			minDist = dist
			colorOut = c
		}
	}

	return
}

// getNearestColor recupera el valor ya calculado, si no existiera lo calcula y
// almacena dentro de la tabla
func getNearestColor(colorIn color.Color, palette []color.Color) color.Color {
	// Si el color no ha sido calculado
	if colorsList[colorIn] == nil {
		// Lo calculamos y salvamos dentro de la tabla
		colorsList[colorIn] = calcColor(colorIn, palette)
	}

	// recuperamos el color más cercano disponible
	return colorsList[colorIn]
}

// processImageWithPalette reemplaza los píxeles de la imagen seleccionada con
// los correspondientes de la paleta.
func processImageWithPalette(palette []color.Color, img image.Image) (image.Image, time.Duration) {
	// Inicializamos la tabla de hash
	colorsList = make(map[color.Color]color.Color)

	// Almacenamos la hora de inicio
	startTime := time.Now()
	// Generamos una imagen vacía
	res := image.NewRGBA(img.Bounds())

	// Por cada pixel de la imagen original
	for y := img.Bounds().Dy(); y >= 0; y-- {
		for x := img.Bounds().Dx(); x >= 0; x-- {
			// Escribimos un pixel con el color más cercano disponible en la
			// paleta, dentro de la imagen nueva.
			c := getNearestColor(img.At(x, y), palette)
			//fmt.Printf("%v, %v: %v\n", x, y, c)
			res.Set(x, y, c)
		}
	}

	// Calculamos el tiempo que tomó el proceso
	return res, time.Since(startTime)
}

func init() {
	imgSrc = flag.String("img", "", "-img ./example.jpg")
	palettes = flag.String("palettes", "./", "-palettes ./palettes/")

	flag.Parse()
	colorsList = make(map[color.Color]color.Color)
}

func main() {
	log.Println("- Obteniendo imágen")
	img, err := fetchImage()
	if err != nil {
		log.Fatal(err)
	}

	// Obtenemos los archivos del directorio de paletas
	d, err := ioutil.ReadDir(*palettes)
	if err != nil {
		log.Fatal(err)
	}

	if len(d) < 1 {
		log.Fatal(errors.New("no files in palettes dir"))
	}

	for _, f := range d {
		if f.IsDir() {
			continue
		}

		// Solo procesamos los archivos con extención .hex
		if f.Name()[len(f.Name())-4:] == ".hex" {

			// Cargamos la paleta de colores
			paletteFile := path.Join(*palettes, f.Name())
			cPalette, err := loadPalette(paletteFile)
			if err != nil {
				log.Printf("error loading palette file [%v], %v", f.Name(), err)
			}

			// Procesamos la imagen con la paleta seleccionada
			res, t := processImageWithPalette(cPalette, img)

			// Generamos el nombre de la imagen resultante
			strReplace := fmt.Sprintf("-%s-out.jpg", strings.Replace(f.Name(), ".hex", "", -1))
			outFile := strings.Replace(*imgSrc, ".jpg", strReplace, -1)

			// Creando archivo de salida
			f, err := os.Create(outFile)
			if err != nil {
				log.Fatal(err)
			}

			// Escribimos el resultado.jpg
			err = jpeg.Encode(f, res, &jpeg.Options{Quality: 100})
			if err != nil {
				log.Fatal(err)
			}

			// Cerramos el archivo resultante
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}

			//Imprimimos los resultados
			log.Printf("- Finalizando (%v seg.): %s\n", t, outFile)
		}
	}
}
