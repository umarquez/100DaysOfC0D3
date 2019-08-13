package main

/*
 * - Singleton Machine -
 * Solo contamos con una máquina y debemos asegurar una sola instancia de esta.
 */

// CookiesMachine protege la instancia de la máquina
type CookiesMachine interface {
	SetBuilder(CookieBuilder)
	MakeCookie() AnyCookie
}

var machineInstance CookiesMachine

// GetMachineInstance devuelve la única instancia de la máquina
func GetMachineInstance() CookiesMachine {
	if machineInstance == nil {
		machineInstance = new(machine)
	}
	return machineInstance
}

// machine es el componente DIRECTOR, encargado de ejecutar el proceso de
// creación de cada componente en el orde adecuado y devolver el resutlado.
type machine struct {
	currentBuilder CookieBuilder
}

// SetBuilder establece el tipo de galleta (objeto) a cocinar (instanciar)
func (m *machine) SetBuilder(builder CookieBuilder) {
	m.currentBuilder = builder
}

// MakeCookie fabrica una galleta y devuelve el resultado.
func (m *machine) MakeCookie() AnyCookie {
	builder := m.currentBuilder
	builder.MixIngredients().ServeTheDough().OvenIt().ChillAndPackage()
	return builder.GetTheCookie()
}
