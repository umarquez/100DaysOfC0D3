/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [29/100] Box blur
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

const imgW = 1024
const imgH = 768

// avgPixel Devuelve el valor promedio de la posición actual y los vecinos
// dentro de la ventana
func avgPixel(windowSize, x, y int, original image.Image) color.Color {
	// distancia del pixel al borde de la ventana
	half := windowSize / 2

	rSum := 0
	gSum := 0
	bSum := 0
	pCount := 0

	// Recorremos la ventana
	for ny := y - half; ny < y+half; ny++ {
		for nx := x - half; nx < x+half; nx++ {
			// si no es una posición válida, la omitimos
			if ny < 0 || ny > original.Bounds().Dy()-1 || nx < 0 || nx > original.Bounds().Dx()-1 {
				continue
			}

			// Obtenemos el valor RGB del pixel y lo añadimos a la suma total
			r, g, b, _ := color.RGBAModel.Convert(original.At(nx, ny)).RGBA()
			rSum += int(uint8(r))
			gSum += int(uint8(g))
			bSum += int(uint8(b))
			pCount++
		}
	}

	// si caemos en este caso, hemos procesado una posición fuera de la imagen
	// o inválida...
	if pCount == 0 {
		fmt.Printf("%v, %v, %v\n", x, y, pCount)
	}

	// Retornamos el color promedio resultante
	return color.RGBA{
		R: uint8(rSum / pCount),
		G: uint8(gSum / pCount),
		B: uint8(bSum / pCount),
		A: 0xFF,
	}
}

// averageBlur genera una imagen desenfocada usando el método Box Blur
func averageBlur(original image.Image, windowSize int) image.Image {
	// Genenamos una imagen vacía
	var outImg = image.NewRGBA(image.Rectangle{
		Max: image.Pt(
			imgW,
			imgH,
		),
	})

	// Iteramos sobre cada pixel de la nueva imagen
	for y := 0; y < outImg.Bounds().Dy(); y++ {
		for x := 0; x < outImg.Bounds().Dx(); x++ {
			// Obtenemos el color RGB promedio de la posición actual
			//fmt.Printf("%v, %v\n", x, y)
			newColor := avgPixel(windowSize, x, y, original)
			outImg.Set(x, y, newColor)
		}
	}

	// Retornamos la imagen resultante
	return outImg
}

func main() {
	windows := []int{
		3,
		9,
		27,
		50,
		100,
	}

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

	for _, size := range windows {
		start := time.Now()
		// Procesamos la imagen con el tamaño de la ventana actual
		imgOut := averageBlur(img, int(size))
		elapsed := time.Since(start)

		// Almacenamos la imagen resultante em formato png
		outImage := fmt.Sprintf("./%v_%vx%v_blur_w%v.png", filename, imgW, imgH, size)
		fmt.Printf("Blur process ends [%v seconds], writing file: %v\n", elapsed.Seconds(), outImage)

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
}
