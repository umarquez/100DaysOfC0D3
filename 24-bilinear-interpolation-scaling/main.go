/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [24/100] Escalado de imágenes mediante interpolación bilineal | Image scaling applying bilinear interpolation
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

const imgW = 200
const imgH = 150
const scaleFactor = 5

// LERP Linear Interpolation: https://es.wikipedia.org/wiki/Interpolaci%C3%B3n_lineal
func LERP(from, to float64, slices, step int) float64 {
	if from == to || step == 0 {
		return from
	}
	total := from
	slice := to - from
	slice /= float64(slices)
	amount := slice * float64(step)
	total += math.Round(amount)
	return total
}

func scaleImgBilinear(original image.Image, factor int) (scaled image.Image) {
	var outImg = image.NewRGBA(image.Rectangle{Max: image.Pt(original.Bounds().Dx()*factor, original.Bounds().Dy()*factor)})

	// Colocamos los pixeles originales de la imagen
	for y := original.Bounds().Dy() - 1; y >= 0; y-- {
		for x := original.Bounds().Dx() - 1; x >= 0; x-- {
			outImg.Set(x*factor, y*factor, original.At(x, y))
		}
	}

	// Ahora llenamos las columnas
	for x := 0; x < original.Bounds().Dx(); x++ {
		for y := 0; y < original.Bounds().Dy()-1; y++ {
			//fmt.Printf("Processing block %v, %v\n", x, y)

			colorA := original.At(x, y)
			colorB := original.At(x, y+1)

			aR32, aG32, aB32, _ := color.RGBAModel.Convert(colorA).RGBA()
			bR32, bG32, bB32, _ := color.RGBAModel.Convert(colorB).RGBA()

			aR, aG, aB := float64(uint8(aR32)), float64(uint8(aG32)), float64(uint8(aB32))
			bR, bG, bB := float64(uint8(bR32)), float64(uint8(bG32)), float64(uint8(bB32))

			//fmt.Printf("From %v, %v, %v, to %v, %v, %v\n", aR, aG, aB, bR, bG, bB)

			// Generamos la columna interpolando los valores existentes
			for i := 0; i < factor; i++ {
				RR := LERP(aR, bR, scaleFactor, i)
				GG := LERP(aG, bG, scaleFactor, i)
				BB := LERP(aB, bB, scaleFactor, i)

				//fmt.Printf("%v - Result: %v, %v, %v\n", i, RR, GG, BB)

				px, py := x*factor, y*factor+i
				tColor := color.RGBA{
					R: uint8(RR),
					G: uint8(GG),
					B: uint8(BB),
					A: 0xFF,
				}
				outImg.SetRGBA(px, py, tColor)

				//	Y si ya estamos en la segunda columna, podemos comenzar a
				//	interpolar las filas internas.
				if x > 0 {
					// Unimos cada columna anterior interpolando hasta la columna actual
					for j := 0; j < factor; j++ {
						npx := (x - 1) * factor
						colorX := outImg.At(npx, py)
						xR32, xG32, xB32, _ := color.RGBAModel.Convert(colorX).RGBA()
						xR, xG, xB := float64(uint8(xR32)), float64(uint8(xG32)), float64(uint8(xB32))

						xRR := LERP(xR, RR, scaleFactor, j)
						xGG := LERP(xG, GG, scaleFactor, j)
						xBB := LERP(xB, BB, scaleFactor, j)

						//fmt.Printf("%v - Result: %v, %v, %v\n", j, xRR, xGG, xBB)
						outImg.SetRGBA(npx+j, py, color.RGBA{
							R: uint8(xRR),
							G: uint8(xGG),
							B: uint8(xBB),
							A: 0xFF,
						})
					}
				}
			}
		}
	}

	// Rellenamos los espacios restantes pues no tenemos información para
	// realizar la interpolación por lo que repetiremos, opcionalmente ser
	// podría interpolar al pixel del extremo contrario (toroide), pero aquí
	// simplemente repetimos los últimos valores.
	x := 0
	for x <= (original.Bounds().Dx()-1)*factor {
		for y := (original.Bounds().Dy() - 1) * factor; y < outImg.Bounds().Dy(); y++ {
			outImg.Set(x, y, outImg.At(x, y-1))
		}
		x++
	}

	y := 0
	for y <= (original.Bounds().Dy())*factor {
		for x := (original.Bounds().Dx() - 1) * factor; x < outImg.Bounds().Dx(); x++ {
			outImg.Set(x, y, outImg.At(x-1, y))
		}
		y++
	}

	return outImg
}

// scaleImg Toma la imagen dada y la escala utilizando una versión de
// interpolación por vecino cercano.
func ScaleImgNearestNeighbor(original image.Image, factor int) (scaled image.Image) {

	// Genenamos una imagen vacía
	var outImg = image.NewRGBA(image.Rectangle{
		Max: image.Pt(
			original.Bounds().Dx()*factor,
			original.Bounds().Dy()*factor,
		),
	})

	// Iteramos sobre cada pixel de la nueva imagen
	for x := 0; x < outImg.Bounds().Dx(); x++ {
		for y := 0; y < outImg.Bounds().Dy(); y++ {
			// Obtenemos la posición del pixel original de acuerdo a la
			// posición actual
			fx, fy := x/factor, y/factor
			// Copiamos el color del pixel original en la posición actual
			outImg.Set(x, y, original.At(fx, fy))
		}
	}

	// Retornamos la imagen escalada resultante
	return outImg
}

func main() {
	// Obtenemos una imagen aleatoria
	img, err := RandomImage(image.Point{imgW, imgH})
	if err != nil {
		log.Fatal(err)
	}

	// Nombre de la imagen basado en la fecha/hora
	filename := time.Now().Unix()
	orgnImage := fmt.Sprintf("./%v_%vx%v.jpg", filename, imgW, imgH)
	f, err := os.Create(orgnImage)
	if err != nil {
		log.Fatal(err)
	}

	// Almacenamos la imagen obtenida
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Printf("error closing originam image file, %v", err)
	}

	// Para utilizar una imagen local:
	/*filename := "source"
	content, err := ioutil.ReadFile("./source.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(bytes.NewBuffer(content))*/

	// Escalamos la imagen utilizando interpolación bilineal
	imgOut := scaleImgBilinear(img, scaleFactor)

	// Almacenamos la imagen resultante em formato png
	outImage := fmt.Sprintf("./%v_%vx%v_bli.png", filename, imgW*scaleFactor, imgH*scaleFactor)
	f, err = os.Create(outImage)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(f, imgOut)
	//err = jpeg.Encode(f, imgOut, &jpeg.Options{Quality:100})
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Printf("error closing originam image file, %v", err)
	}

	// Escalamos la imagen utilizando interpolación bilineal
	imgOut = ScaleImgNearestNeighbor(img, scaleFactor)

	// Almacenamos la imagen resultante em formato png sin perdidas para poder
	// contar el mayor detalle posible
	outImage = fmt.Sprintf("./%v_%vx%v_nni.png", filename, imgW*scaleFactor, imgH*scaleFactor)
	f, err = os.Create(outImage)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(f, imgOut)
	//err = jpeg.Encode(f, imgOut, &jpeg.Options{Quality:100})
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Printf("error closing originam image file, %v", err)
	}
}
