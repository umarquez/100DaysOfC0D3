package main

import (
	"fmt"
	"log"
	"time"
)

const ChocolateChipsCookieType CookieType = "Chocolate Chips"

// ChocolateChipsCookie es el objeto que vamos a ensamblar mediente el builder,
// siguiendo la receta correspondiente, este implementa los métodos de la
// interfaz AnyCookie y podrá ser utilizado como tal.
type ChocolateChipsCookie struct {
	ingredients    map[string]float64
	weigth         float64
	flavor         string
	expirationDate time.Time
}

// GetDescription devuelve la descripción de la galleta
func (cookie *ChocolateChipsCookie) GetDescription() string {
	return fmt.Sprintf("%v cookie.", cookie.flavor)
}

// GetIngredients devuelve el nombre de los ingredientes con los que se fabricó la galleta
func (cookie *ChocolateChipsCookie) GetIngredients() []string {
	var ingredients []string
	for k := range cookie.ingredients {
		ingredients = append(ingredients, k)
	}
	return ingredients
}

// GetWeightGr devuelve el peso de la galleta
func (cookie *ChocolateChipsCookie) GetWeightGr() float64 {
	return cookie.weigth
}

// GetExpirationDate devuelve la fecha de caducidad
func (cookie *ChocolateChipsCookie) GetExpirationDate() time.Time {
	return cookie.expirationDate
}

/*
 * Galletas de chispas de chocolate
 * https://peopleenespanol.com/recetas/829-galletas-de-chispas-de-chocolate-facil-simas/
 */
// ChocolateChipsCookiesBuilder encapsula los pasos para fabricar una galleta
//  de chispas de chocolate utilizando la interfaz CookieBuilder.
type ChocolateChipsCookiesBuilder struct {
	currentCookie *ChocolateChipsCookie
	processStatus internalCookieStatus
}

// MixIngredients es el primer paso del preparado, inicializa la galleta con
// la que trabajaremos y agrega los ingredientes de esta.
func (builder *ChocolateChipsCookiesBuilder) MixIngredients() CookieBuilder {
	builder.currentCookie = new(ChocolateChipsCookie)

	// - Bate el huevo y el azúcar juntos hasta lograr una consistencia
	// espumosa y espesa.
	builder.currentCookie.ingredients = map[string]float64{
		"huevo":            1,
		"azúcar mascabado": 1, // taza
	}

	// - Derrite la mantequilla y añade a la mezcla de huevo.
	builder.currentCookie.ingredients["mantequilla"] = 125 // gramos

	// - Agrega el harina, sal y chispas de chocolate. Mezcla bien.
	builder.currentCookie.ingredients["harina leudante"] = 1.5    // tazas
	builder.currentCookie.ingredients["sal"] = 1                  // pizca
	builder.currentCookie.ingredients["chispas de chocolate"] = 1 // taza

	builder.processStatus = mixed
	return builder
}

// ServeTheDough debe ser el segundo paso, este divide la masa en las porciones
// individuales, en este caso 6, para poder generar una galleta.
func (builder *ChocolateChipsCookiesBuilder) ServeTheDough() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != mixed {
		log.Printf("you need to mix the ingredients before serve")
		return nil
	}

	// - Coloca la masa a cucharadas (una cucharada colmada por galleta) sobre
	// una charola para hornear forrada con papel encerado.
	for name := range builder.currentCookie.ingredients {
		builder.currentCookie.ingredients[name] /= 6 // Porciones
	}
	builder.processStatus = served
	return builder
}

// OvenIt Hornea la charola con la masa, este proceso podría tomar algún tiempo
func (builder *ChocolateChipsCookiesBuilder) OvenIt() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != served {
		log.Printf("you need to serve the mix on a tray before bake it")
		return nil
	}

	// - Hornea durante 15 minutos.

	// caducidad 30 días después, en el empaque
	builder.currentCookie.expirationDate = time.Now().Add(30 * 24 * time.Hour)
	builder.processStatus = baked
	return builder
}

// ChillAndPackage Deja enfriar y coloca cada galleta en un empaque individual
func (builder *ChocolateChipsCookiesBuilder) ChillAndPackage() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != baked {
		log.Printf("you need to bake before package the cookie")
		return nil
	}

	// Deja que las galletas se enfríen en la charola.
	builder.currentCookie.flavor = "chocolate chips"
	builder.currentCookie.weigth = 50

	builder.processStatus = packaged
	return builder
}

// GetTheCookie Valida y recupera la galleta fabricada
func (builder *ChocolateChipsCookiesBuilder) GetTheCookie() AnyCookie {
	cCookie := builder.currentCookie
	if !(cCookie.weigth > 0 &&
		cCookie.flavor != "" &&
		cCookie.expirationDate.After(time.Now()) &&
		len(cCookie.ingredients) > 0) ||
		cCookie == nil ||
		builder.processStatus != packaged {
		log.Printf("this cookie is not ready yet, complete all the steps or start with a new cookie")
		return nil
	}

	builder.currentCookie = nil
	return cCookie
}
