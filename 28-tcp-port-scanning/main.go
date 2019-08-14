/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [28/100] TCP Port Scanning
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// TestPort intenta realizar una conexión al host/puerto seleccionado, cerrando
// la conexión inmediatamente si esta ha sido satisfactoria y devolviendo el
// estado del puerto como resultado
func TestPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%v", host, port)

	cnn, err := net.Dial("tcp", address)
	if err != nil {
		// Evaluamos el tipo de error
		switch err.(type) {
		case *net.OpError: // Asumimos que el puerto está cerrado
		default: // Error desconocido...
			log.Printf("[%T] error dialing %s, %v", err, address, err)
		}
		return false
	}

	// Cerramos la conexión.
	err = cnn.Close()
	if err != nil {
		log.Printf("[%T] error closing cnn %s, %v", err, address, err)
	}
	return true
}

var targetHost string
var startPort int
var endPort int
var testServer int

// PrintStatus
func PrintStatus(currentPort *int) {
	progress := (float64(*currentPort-startPort) * 100) / float64(endPort-startPort)
	log.Printf("- [port %v] scanning process: %v%%\n", *currentPort, progress)
}

// ForcePrintStatus imprime el estado del escaneo y ayuda a salir de la app
func ForcePrintStatus(currentPort *int) {
	catcher := ""
	// Si no es necesario terminar la aplicación, continuamos
	for strings.ToLower(catcher) != "q" {
		_, _ = fmt.Scanln(&catcher)
		PrintStatus(currentPort)
	}

	log.Print("- ENDING SCAN PROCESS -\n")
	os.Exit(0)
}

func init() {
	flag.StringVar(&targetHost, "host", "localhost", "-host 127.0.0.1 host's ip to scan")
	flag.IntVar(&startPort, "start", 31300, "-start 80 starting port")
	flag.IntVar(&endPort, "end", 31400, "-end 443 ending port")
	flag.IntVar(&testServer, "testing-server-port", 31337, "-testing-server-port 1234 (0 = disabled) testing echo server port")
	flag.Parse()

	if startPort > endPort {
		log.Printf("starting port must be lower than the ending port")
	}
}

func main() {
	// Debemos iniciar el servidor de pruebas?
	if testServer > 0 {
		go func() {
			startServer(testServer)
		}()
		time.Sleep(time.Second)
	}

	log.Printf("- Scanning host: %v ports: %v-%v\n", targetHost, startPort, endPort)

	// Si solo es un puerto, mostramos el resultado sin entrar al loop
	if startPort == endPort {
		if TestPort(targetHost, startPort) {
			log.Printf("- [port %v] is open.\n", startPort)
		}
		return
	}

	// inicializamos la variable que nos ayudará a detectar cuando el
	// porcentaje de progreso cambie.
	last := 0

	// puerto actual
	var p = startPort

	// Inicializamos la rutina que detectará cada vez que presionemos ENTER
	go ForcePrintStatus(&p)
	PrintStatus(&p)

	// Recorremos el rango de puertos
	for ; p <= endPort; p++ {
		// Verificamos e imprimimos el estado de puerto en caso de contar con
		// una conexión exitosa
		if TestPort(targetHost, p) {
			log.Printf("- [port %v] is open.\n", p)
		}

		// porcentaje actual
		progress := (float64(p-startPort) * 100) / float64(endPort-startPort)
		// ¿Ha cambiado y ha sido en un 10% por lo menos?
		if int(progress) > last && (int(progress)+10)%10 == 0 {
			// Almacenamos el nuevo valor
			last = int(progress)
			// Imprimimos el progreso
			PrintStatus(&p)
		}
	}
}
