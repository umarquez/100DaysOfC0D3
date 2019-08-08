/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [25/100] Tramado | Dithering
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
	"image/png"
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
const imgDefaultWidth = 800
const imgDefaultHeight = 600

// Image ext
const jpgExt = ".jpg"
const outDir = "out"

const idR = "R"
const idG = "G"
const idB = "B"

// Source files paths
var imgSrc, palettes *string

// Color's table
var colorsList map[color.Color]color.Color
var cErrors map[color.Color]map[string]float64

// Error distribution table
var errorDistribution map[image.Point]map[string]float64

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
// calcColor calcula el vecino más cercano de la paleta al color seleccionado y
// lo retorna además del error de cada valor RGB
func calcColor(colorIn color.Color, palette []color.Color) (colorOut color.Color, errors map[string]float64) {
	errors = make(map[string]float64)
	oR, oG, oB, _ := colorIn.RGBA()
	minDist := math.Sqrt(1E18) // Max distance possible

	for _, c := range palette {
		pR, pG, pB, _ := color.NRGBAModel.Convert(c).RGBA()
		dR, dG, dB := float64(uint8(oR))-float64(uint8(pR)), float64(uint8(oG))-float64(uint8(pG)), float64(uint8(oB))-float64(uint8(pB))
		dR2, dG2, dB2 := dR*dR, dG*dG, dB*dB
		dist := math.Sqrt(dR2 + dG2 + dB2)
		if dist < minDist {
			minDist = dist
			colorOut = c
			errors[idR] = dR
			errors[idG] = dG
			errors[idB] = dB
		}
	}

	return
}

// getNearestColor recupera el valor ya calculado, si no existiera lo calcula y
// almacena dentro de la tabla
func getNearestColor(colorIn color.Color, palette []color.Color) (color.Color, map[string]float64) {
	// Si el color no ha sido calculado
	if colorsList[colorIn] == nil {
		// Lo calculamos y salvamos dentro de la tabla
		colorsList[colorIn], cErrors[colorIn] = calcColor(colorIn, palette)
	}

	// recuperamos el color más cercano disponible
	return colorsList[colorIn], cErrors[colorIn]
}

// processImageWithPalette reemplaza los píxeles de la imagen seleccionada con
// los correspondientes de la paleta.
func processImageWithPalette(palette []color.Color, img image.Image, dither DitherFunc) (image.Image, time.Duration) {
	// Inicializamos la tabla de hash
	errorDistribution = make(map[image.Point]map[string]float64)
	colorsList = make(map[color.Color]color.Color)
	cErrors = make(map[color.Color]map[string]float64)

	// Almacenamos la hora de inicio
	startTime := time.Now()
	// Generamos una imagen vacía
	res := image.NewRGBA(img.Bounds())

	// Por cada pixel de la imagen original
	for y := 0; y <= img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			// Inicializamos el map que almacenará el error correspondientes a
			// cada valor RGB
			colorModifiers := make(map[string]float64)

			// si ya existen valores almacenados para este píxel lo recuperamos
			if errorDistribution[image.Point{X: x, Y: y}] != nil {
				colorModifiers = errorDistribution[image.Point{X: x, Y: y}]
			}

			// Obtenemos los valores RGB del pixel actual
			r, g, b, _ := color.NRGBAModel.Convert(img.At(x, y)).RGBA()
			newR, newG, newB := float64(uint8(r)), float64(uint8(g)), float64(uint8(b))

			// Y sumamos el erros acumulado a cada valor
			newR += colorModifiers[idR]
			newG += colorModifiers[idG]
			newB += colorModifiers[idB]

			// Verificamos que cada valor se encuentre dentro del rango
			if newR > 0xFF {
				newR = 0xFF
			} else if newR < 0 {
				newR = 0
			}
			if newG > 0xFF {
				newG = 0xFF
			} else if newG < 0 {
				newG = 0
			}
			if newB > 0xFF {
				newB = 0xFF
			} else if newB < 0 {
				newB = 0
			}

			// Y generamos el color resultante
			ditheredColor := color.RGBA{
				R: uint8(math.Round(newR)),
				G: uint8(math.Round(newG)),
				B: uint8(math.Round(newB)),
			}

			// Escribimos un pixel con el color más cercano disponible en la
			// paleta, dentro de la imagen nueva y recuperamos el error
			c, e := getNearestColor(ditheredColor, palette)
			//fmt.Printf("%v, %v: %v\n", x, y, c)
			res.Set(x, y, c)

			// Si se ha definido un método para distribuir el error, este será
			// invocado pasando como parámetro la ubicación del pixel, el error
			// a distribuir y el mapa que almacena los valores.
			if dither != nil {
				dither(image.Point{X: x, Y: y}, e, errorDistribution)
			}
		}
	}

	// Calculamos el tiempo que tomó el proceso
	return res, time.Since(startTime)
}

func init() {
	errorDistribution = make(map[image.Point]map[string]float64)
	colorsList = make(map[color.Color]color.Color)
	cErrors = make(map[color.Color]map[string]float64)
	imgSrc = flag.String("img", "", "-img ./example.jpg")
	palettes = flag.String("palettes", "./", "-palettes ./palettes/")

	flag.Parse()
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

			// Generamos una versión para cada método de tramado disponible
			for ditherName, dFilter := range DitherFuncsCat {
				// Procesamos la imagen con la paleta seleccionada
				res, t := processImageWithPalette(cPalette, img, dFilter)

				// Generamos el nombre de la imagen resultante
				strReplace := fmt.Sprintf("-%s-%s.png", strings.Replace(f.Name(), ".hex", "", -1), ditherName)
				outFile := strings.Replace(*imgSrc, ".jpg", strReplace, -1)

				// Creando archivo de salida
				f, err := os.Create(outFile)
				if err != nil {
					log.Fatal(err)
				}

				// Escribimos el resultado.tif
				err = png.Encode(f, res)
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
}
