package main

import "image"

// DitherFunc es una función que distribuye el error de un pixel, producto de
// la sustitución de color dentro de la matriz definida
type DitherFunc func(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64)

// DifuseWithFilter distribuye el error correspondiente a cada valor dentro de
// una matriz utilizando el filtro seleccionado
func DifuseWithFilter(values map[string]float64, filter map[image.Point]float64, position image.Point, matrix map[image.Point]map[string]float64) {
	for offset, factor := range filter {
		if matrix[position.Add(offset)] == nil {
			matrix[position.Add(offset)] = make(map[string]float64)
		}
		matrix[position.Add(offset)][idR] += factor * values[idR]
		matrix[position.Add(offset)][idG] += factor * values[idG]
		matrix[position.Add(offset)][idB] += factor * values[idB]
	}
}

// FloydSteinbergFilter divide el error en 16 partes y lo distribuye de acuerdo
// al algoritmo
func FloydSteinbergFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 16
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  7 / sectionsNum,
		image.Point{X: 1, Y: 1}:  1 / sectionsNum,
		image.Point{X: 0, Y: 1}:  5 / sectionsNum,
		image.Point{X: -1, Y: 1}: 3 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// JarvisJudiceNinkeFilter divide el error en 48 partes y lo distribuye de
// acuerdo al algoritmo
func JarvisJudiceNinkeFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 48
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  7 / sectionsNum,
		image.Point{X: 2, Y: 0}:  5 / sectionsNum,
		image.Point{X: 1, Y: 1}:  5 / sectionsNum,
		image.Point{X: 2, Y: 1}:  3 / sectionsNum,
		image.Point{X: 2, Y: 2}:  1 / sectionsNum,
		image.Point{X: 1, Y: 2}:  3 / sectionsNum,
		image.Point{X: 0, Y: 1}:  7 / sectionsNum,
		image.Point{X: 0, Y: 2}:  5 / sectionsNum,
		image.Point{X: -1, Y: 2}: 3 / sectionsNum,
		image.Point{X: -1, Y: 1}: 5 / sectionsNum,
		image.Point{X: -2, Y: 2}: 1 / sectionsNum,
		image.Point{X: -2, Y: 1}: 3 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// StuckiFilter divide el error en 42 partes y lo distribuye de acuerdo al
// algoritmo
func StuckiFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 42
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  8 / sectionsNum,
		image.Point{X: 2, Y: 0}:  4 / sectionsNum,
		image.Point{X: 1, Y: 1}:  4 / sectionsNum,
		image.Point{X: 2, Y: 1}:  2 / sectionsNum,
		image.Point{X: 2, Y: 2}:  2 / sectionsNum,
		image.Point{X: 1, Y: 2}:  1 / sectionsNum,
		image.Point{X: 0, Y: 1}:  8 / sectionsNum,
		image.Point{X: 0, Y: 2}:  4 / sectionsNum,
		image.Point{X: -1, Y: 2}: 2 / sectionsNum,
		image.Point{X: -1, Y: 1}: 4 / sectionsNum,
		image.Point{X: -2, Y: 2}: 1 / sectionsNum,
		image.Point{X: -2, Y: 1}: 2 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// AtkinsonFilter divide el error en 8 partes y lo distribuye de acuerdo al
// algoritmo
func AtkinsonFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 8
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  1 / sectionsNum,
		image.Point{X: 2, Y: 0}:  1 / sectionsNum,
		image.Point{X: 1, Y: 1}:  1 / sectionsNum,
		image.Point{X: 0, Y: 1}:  1 / sectionsNum,
		image.Point{X: 0, Y: 2}:  1 / sectionsNum,
		image.Point{X: -1, Y: 1}: 1 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// BurkesFilter divide el error en 32 partes y lo distribuye de acuerdo al
// algoritmo
func BurkesFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 32
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  8 / sectionsNum,
		image.Point{X: 2, Y: 0}:  4 / sectionsNum,
		image.Point{X: 1, Y: 1}:  4 / sectionsNum,
		image.Point{X: 2, Y: 1}:  2 / sectionsNum,
		image.Point{X: 0, Y: 1}:  8 / sectionsNum,
		image.Point{X: -1, Y: 1}: 4 / sectionsNum,
		image.Point{X: -2, Y: 1}: 2 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// SierraFilter divide el error en 32 partes y lo distribuye de acuerdo al
// algoritmo
func SierraFilter(point image.Point, values map[string]float64, errorMatrix map[image.Point]map[string]float64) {
	var sectionsNum float64 = 32
	filter := map[image.Point]float64{
		image.Point{X: 1, Y: 0}:  5 / sectionsNum,
		image.Point{X: 2, Y: 0}:  3 / sectionsNum,
		image.Point{X: 1, Y: 1}:  4 / sectionsNum,
		image.Point{X: 2, Y: 1}:  2 / sectionsNum,
		image.Point{X: 1, Y: 2}:  2 / sectionsNum,
		image.Point{X: 0, Y: 1}:  5 / sectionsNum,
		image.Point{X: 0, Y: 2}:  3 / sectionsNum,
		image.Point{X: -1, Y: 2}: 2 / sectionsNum,
		image.Point{X: -1, Y: 1}: 4 / sectionsNum,
		image.Point{X: -2, Y: 1}: 2 / sectionsNum,
	}

	DifuseWithFilter(values, filter, point, errorMatrix)
}

// DitherFuncsCat es un catálogo de métodos de distribución de errores
var DitherFuncsCat = map[string]DitherFunc{
	"plain":           nil,
	"floyd-steinberg": FloydSteinbergFilter,
	"jjn":             JarvisJudiceNinkeFilter,
	"stucki":          StuckiFilter,
	"atkinson":        AtkinsonFilter,
	"burkes":          BurkesFilter,
	"sierra":          SierraFilter,
}
