package main

import "time"

// AnyCookie define los métodos que derá tener una galleta para ser considerada
// como tal, en este caso los métodos retornan la información del empaque
type AnyCookie interface {
	GetDescription() string       // Descripción del producto
	GetIngredients() []string     // Ingredientes
	GetWeightGr() float64         // Peso en gr.
	GetExpirationDate() time.Time // Fecha de caducidad
}

// CookieBuilder es cualquier receta para cocinar una galleta que cuente con
// los siguiente pasos (métodos), mismo que al ser invocados en orden, deberán
// producir una nueva galleta (AnyCookie) como resultado.
type CookieBuilder interface {
	MixIngredients() CookieBuilder  // Mezclar los ingredientes
	ServeTheDough() CookieBuilder   // Distribuir y servir en porciones
	OvenIt() CookieBuilder          // Hornear
	ChillAndPackage() CookieBuilder // Enfriar y empacar
	GetTheCookie() AnyCookie        // Recuperar la galleta fabricada.
}
