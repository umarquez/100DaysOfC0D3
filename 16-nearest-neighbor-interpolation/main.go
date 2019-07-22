package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
)

// Tamaño original
const imgW = 200
const imgH = 150

// Factor de escalado
const scaleFactor = 5

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
	// Ya sabes, la imagen aleatoria...
	img, err := RandomImage(image.Point{imgW, imgH})
	if err != nil {
		log.Fatal(err)
	}

	// El nombre aleatorio para la imagen aleatoria del mundo aleatorio...
	filename := time.Now().Unix()
	orgnImage := fmt.Sprintf("./%v_%vx%v.jpg", filename, imgW, imgH)
	f, err := os.Create(orgnImage)
	if err != nil {
		log.Fatal(err)
	}

	// Guardamos la imagen original
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}

	// Cerramos el archivo
	err = f.Close()
	if err != nil {
		// Y si ocurre un error lo mostramos y seguimos, un tropiezo nada más...
		log.Printf("error closing originam image file, %v", err)
	}

	/************************************
	Aquí ocurre la magia:               */
	imgOut := ScaleImgNearestNeighbor(img, scaleFactor) // Qué imagen y cuánto la vamos a escalar
	/************************************/

	// Guardamos el resultado
	outImage := fmt.Sprintf("./%v_%vx%v.jpg", filename, imgW*scaleFactor, imgH*scaleFactor)
	f, err = os.Create(outImage)
	if err != nil {
		log.Fatal(err)
	}

	err = jpeg.Encode(f, imgOut, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}

	// Y terminamos cerrando el archivo
	err = f.Close()
	if err != nil {
		log.Printf("error closing originam image file, %v", err)
	}
}
