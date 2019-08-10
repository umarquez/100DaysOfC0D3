/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [26/100] WS Shell - Client
 */
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

const shellPath = "/26-ws-shell"
const password = "0xC0D3C0D3"

// wsRead Recupera los datos enviados a través del socket y lo escribe en la
// salida de la aplicación
func wsRead(ws *websocket.Conn, wait *sync.WaitGroup) {
	// cierra el socket al salir
	defer wait.Done()

	// definimos un buffer de lectura
	msj := make([]byte, 1024)
	for {
		// Si el socket no ha sido inicializados salimos
		if ws == nil {
			return
		}

		// leemos los datos de entrada y los colocamos en el buffer
		n, err := ws.Read(msj)
		if err != nil {
			log.Printf("error reading, %v\n", err)
			break
		}

		// Si el buffer no se encuentra vacío
		if n > 0 {
			// recuperamos los datos del buffer
			readed := msj[:n]
			// y los imprimimos en la salida
			fmt.Printf("%s", readed)
		}
	}
}

// wsWrite escribe en el socket el texto ingresado cada vez que se completa una
// línea, es decir, cada que se preciona la tecla ENTER
func wsWrite(ws *websocket.Conn, wait *sync.WaitGroup) {
	defer func() {
		_ = ws.Close()
		wait.Done()
	}()

	// Escanearemos la entrada estandar para obtener una línea cada vez
	inScanner := bufio.NewScanner(os.Stdin)
	for inScanner.Scan() {
		// Agregamos un salto de línea al final del texto obtenido
		input := []byte(inScanner.Text() + "\n")
		// Escribimos el texto en el WS
		_, err := ws.Write(input)
		if err != nil {
			log.Printf("error writting, %v\n", err)
			break
		}
	}
}

// Nos ayudará a realizar el proceso de autenticación
func auth(host string) (string, error) {
	// Obtenemos el token
	res, err := http.Post(host+"/login", "text/plain", nil)
	if err != nil {
		return "", fmt.Errorf("error dialing, %v\n", err)
	}

	id, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("login error, %v\n", err)
	}

	// ID de sesión
	strid := string(id)
	fmt.Printf("session: %s\n", strid)

	// Enviamos la contraseña al endpoint
	content := bytes.NewBufferString(password)
	url := fmt.Sprintf("%s/%s", host, strid)
	res, err = http.Post(url, "text/plain", content)
	if err != nil {
		return "", fmt.Errorf("error sending data, %v\n", err)
	}

	// Verificamos la contraseña
	res, err = http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error validating, %v\n", err)
	}

	valContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading validation body, %v\n", err)
	}

	if string(valContent) != "OK" {
		return "", fmt.Errorf("invalid credentials\n")
	}

	return strid, nil
}

func main() {
	// Obtenemos el nombre:puerto del servidor
	var serverHostPort string
	var ref string
	var err error
	flag.StringVar(&serverHostPort, "server", "localhost:8080", "-server=\"localhost:8080\"")
	flag.Parse()

	if ref, err = auth("http://" + serverHostPort); err != nil {
		log.Fatal(errors.New("authentication error|"))
	}

	// Intentamos obtener un websocket accediendo a la ruta correspondiente
	ws, err := websocket.Dial("ws://"+serverHostPort+shellPath, "", "http://"+serverHostPort+"/"+ref)
	if err != nil {
		log.Printf("error dialing, %v\n", err)
	}

	// Una vez establecida la conexión deberemos lanzamos las rutinas de
	// lectura/escritura, el WaitGroup nos ayuda a mantener el programa en
	// ejecucíón mientras el WS esté activo
	wait := new(sync.WaitGroup)
	wait.Add(1)
	go wsRead(ws, wait)
	go wsWrite(ws, wait)
	wait.Wait()
}
