/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 * - Nearest Neighbor 2D
 *             @umarquez
 */
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const mapWidth int = 800  // ancho del mapa
const mapHeight int = 600 // alto del mapa
const animalsNum int = 50 // num de animales
const foodNum int = 70    // num de comida
const animalSize = 9      // tamaño en px
const foodSize = 5        // tamaño en px

// Entity será cualquier elemento dentro del mapa
type Entity struct {
	position image.Point
}

// SetPosition establece las coordenadas de la entidad
func (ntt *Entity) SetPosition(point image.Point) {
	ntt.position = point
}

// Food es una entidad consumible por animales
type Food struct {
	Entity
	Color color.Color
}

// Draw dibuja una entidad Food en una imagen RGBA
func (food Food) Draw(dest *image.RGBA) {
	wSize := foodSize / 2
	for x := food.position.X - wSize; x <= food.position.X+wSize; x++ {
		for y := food.position.Y - wSize; y <= food.position.Y+wSize; y++ {
			dest.Set(x, y, food.Color)
		}
	}
}

// Animal es una entidad que consume Food cercana a él
type Animal struct {
	Entity
	Color color.Color
}

// Draw dibuja una entidad Animal en una imagen RGBA
func (animal Animal) Draw(dest *image.RGBA) {
	wSize := animalSize / 2
	for x := animal.position.X - wSize; x <= animal.position.X+wSize; x++ {
		for y := animal.position.Y - wSize; y <= animal.position.Y+wSize; y++ {
			dest.Set(x, y, animal.Color)
		}
	}
}

// FindNearestFood busca el elemento Food más cercano al Animal actual
func (animal Animal) FindNearestFood(foods []*Food) *Food {
	var nearest *Food
	var distance float64 = 0xFFFFFF // distancia inicial muy grande

	// Itera sobre la lista de Foods disponibles
	for _, f := range foods {
		// primer cateto
		a := float64(animal.position.X) - float64(f.position.X)
		// segundo cateto
		b := float64(animal.position.Y) - float64(f.position.Y)
		// teorema de Pitágoras
		cDist := math.Sqrt((a * a) + (b * b))

		// Si es la distancia menor hasta ahora
		if cDist < distance {
			distance = cDist // Guardamos la distancia
			nearest = f      // Y el elemento más cercano
		}
	}

	// Retornamos el elemento más cercano.
	return nearest
}

// setLinePx dibuja un pixel de una línea
func setLinePx(x, y int, dest *image.RGBA, rgba color.RGBA) {
	dest.Set(x, y, rgba)
	dest.Set(x+1, y, rgba)
	dest.Set(x-1, y, rgba)
	dest.Set(x, y+1, rgba)
	dest.Set(x, y-1, rgba)
}

// DrawLine conecta dos puntos con una línea dentro de una imágen
func DrawLine(posA, posB image.Point, dest *image.RGBA) {
	lineColor := color.RGBA{B: 0xEE, A: 0xFF}
	distanceX := float64(posB.X - posA.X)
	distanceY := float64(posB.Y - posA.Y)

	// Toma como base la distancia más grande entre `x` y `y`
	if math.Abs(distanceX) > math.Abs(distanceY) {
		stepX := int(distanceX / math.Abs(distanceX))
		for pxX := posA.X; pxX != posB.X; pxX += stepX {
			pxY := posA.Y + int(float64(pxX-posA.X)*(distanceY/distanceX))
			setLinePx(pxX, pxY, dest, lineColor)
		}
	} else {
		stepY := int(distanceY / math.Abs(distanceY))
		for pxY := posA.Y; pxY != posB.Y; pxY += stepY {
			pxX := posA.X + int(float64(pxY-posA.Y)*(distanceX/distanceY))
			setLinePx(pxX, pxY, dest, lineColor)
		}
	}
	dest.Set(posB.X, posB.Y, lineColor)
}

// RandomPoint genera una posición aleratoria dentro de los límites
func RandomPoint(width, height int) image.Point {
	return image.Point{
		X: rand.Intn(width - 1),
		Y: rand.Intn(height - 1),
	}
}

func main() {
	// Aleatorio con una semilla diferente cada vez
	rand.Seed(time.Now().Unix())
	//var animals[] *Animal
	var foods []*Food

	// mapa donde con el que trabajaremos, una imagen en blanco
	terrain := image.NewRGBA(image.Rect(0, 0, mapWidth, mapHeight))
	draw.Draw(terrain, terrain.Bounds(), image.White, image.ZP, draw.Src)

	// Generamos cada comida en una posición aleatorio
	for i := foodNum; i > 0; i-- {
		cFood := &Food{
			Color: color.RGBA{
				R: 0xDD,
				A: 0xFF,
			},
		}
		cFood.SetPosition(RandomPoint(mapWidth, mapHeight))
		cFood.Draw(terrain)
		foods = append(foods, cFood)
	}

	// Generamos cada animal en una posición aleatoria
	for i := animalsNum; i > 0; i-- {
		cAnimal := &Animal{
			Color: color.RGBA{
				B: 0x00,
				A: 0xFF,
			},
		}
		cAnimal.SetPosition(RandomPoint(mapWidth, mapHeight))
		//animals = append(animals, cAnimal)

		// Buscamos la comida más cercana a la posición del animal actual
		nFood := cAnimal.FindNearestFood(foods)

		// Dibujamos una línea desde el animal hasta la comida cercana
		DrawLine(cAnimal.position, nFood.position, terrain)

		// Dibujamos el animal en la posición actual
		cAnimal.Draw(terrain)
	}

	// Almacenamos la imagen resultante
	filename := fmt.Sprintf("./%v.png", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, terrain)
	if err != nil {
		log.Fatal(err)
	}
}
