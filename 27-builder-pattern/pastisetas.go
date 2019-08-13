package main

import (
	"fmt"
	"log"
	"time"
)

const PastisetaCookieType CookieType = "Pastiseta"

// PastisetaCookie es el objeto que vamos a ensamblar mediente el builder,
// siguiendo la receta correspondiente, este implementa los métodos de la
// interfaz AnyCookie y podrá ser utilizado como tal.
type PastisetaCookie struct {
	ingredients    map[string]float64
	weigth         float64
	flavor         string
	expirationDate time.Time
}

// GetDescription devuelve la descripción de la galleta
func (cookie *PastisetaCookie) GetDescription() string {
	return fmt.Sprintf("%v cookie.", cookie.flavor)
}

// GetIngredients devuelve el nombre de los ingredientes con los que se fabricó
// la galleta
func (cookie *PastisetaCookie) GetIngredients() []string {
	var ingredients []string
	for k := range cookie.ingredients {
		ingredients = append(ingredients, k)
	}
	return ingredients
}

// GetWeightGr devuelve el peso de la galleta
func (cookie *PastisetaCookie) GetWeightGr() float64 {
	return cookie.weigth
}

// GetExpirationDate devuelve la fecha de caducidad
func (cookie *PastisetaCookie) GetExpirationDate() time.Time {
	return cookie.expirationDate
}

/*
 * Pastisetas
 * https://www.kiwilimon.com/receta/postres/pastisetas
 */
// PastisetaCookieBuilder encapsula los pasos para fabricar una galleta
//  de chispas de chocolate utilizando la interfaz CookieBuilder.
type PastisetaCookieBuilder struct {
	currentCookie *PastisetaCookie
	processStatus internalCookieStatus
}

// MixIngredients es el primer paso del preparado, inicializa la galleta con
// la que trabajaremos y agrega los ingredientes de esta.
func (builder *PastisetaCookieBuilder) MixIngredients() CookieBuilder {
	builder.currentCookie = new(PastisetaCookie)
	builder.currentCookie.ingredients = make(map[string]float64)

	// - En una batidora , acrema la mantequilla con el azúcar por 3 minutos.
	builder.currentCookie.ingredients["mantequilla"] = 1.25 // tazas
	builder.currentCookie.ingredients["azúcar glass"] = 1   // taza

	// - Agrega la vainilla y continúa batiendo.
	builder.currentCookie.ingredients["vainilla"] = 1 // cucharadita

	// - Alterna los huevos con la harina, empezando y terminando con la harina
	// Bate hasta integrar.
	builder.currentCookie.ingredients["huevo"] = 2  // piezas
	builder.currentCookie.ingredients["harina"] = 3 // tazas

	builder.processStatus = mixed

	return builder
}

// ServeTheDough debe ser el segundo paso, este divide la masa en las porciones.
func (builder *PastisetaCookieBuilder) ServeTheDough() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != mixed {
		log.Printf("you need to mix the ingredients before serve")
	}

	// - Coloca la mezcla en una manga con dulla rizada y forma las galletas.
	for name := range builder.currentCookie.ingredients {
		builder.currentCookie.ingredients[name] /= 12 // Porciones
	}
	// - Refrigera 15 minutos

	builder.processStatus = served
	return builder
}

// OvenIt Hornea la charola con la masa, este proceso podría tomar algún tiempo
func (builder *PastisetaCookieBuilder) OvenIt() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != served {
		log.Printf("you need to serve the mix on a tray before bake it")
		return nil
	}

	// - Hornea 15 minutos o hasta que los bordes estén dorados. Enfría.

	// caducidad 10 días después, en el empaque
	builder.currentCookie.expirationDate = time.Now().Add(10 * 24 * time.Hour)
	builder.processStatus = baked
	return builder
}

// ChillAndPackage Deja enfriar y coloca cada galleta en un empaque individual
func (builder *PastisetaCookieBuilder) ChillAndPackage() CookieBuilder {
	if builder.currentCookie == nil || builder.processStatus != baked {
		log.Printf("you need to bake before package the cookie")
		return nil
	}

	// Deja que las galletas se enfríen en la charola.
	builder.currentCookie.flavor = "pastiseta"
	builder.currentCookie.weigth = 15

	builder.processStatus = packaged
	return builder
}

// GetTheCookie Valida y recupera la galleta fabricada
func (builder *PastisetaCookieBuilder) GetTheCookie() AnyCookie {
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
