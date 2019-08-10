/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [26/100] WS Shell - Server
 */
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/crypto/sha3"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const shellPath = "/26-ws-shell"
const loginPath = "/login"
const winCmd = "cmd"
const lnxCmd = "/bin/bash"
const password = "0xC0D3C0D3"

var sessions map[string]string

func decodeId(reqUrl *url.URL) string {
	return strings.Replace(reqUrl.Path, "/", "", -1)
}

// wsCmd detecta el sistema operativo e inicializa una 26-ws-shell que adjunta al WS
func wsCmd(ws *websocket.Conn) {
	var err error

	// Recuperamos la ruta de la URL para obtener el ID
	refUrl := ws.Config().Origin
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	id := decodeId(refUrl)

	// verificamos la contraseña almacenada, en caso de error salimos de la
	// función actual
	if sessions[id] != password {
		_ = ws.Close()
		return
	}

	// Imprimimos el origen de la conexión
	fmt.Printf("Connecting from %s\n", ws.Request().RemoteAddr)

	var cmdAttr []string

	// Detección del OS
	strCmd := ""
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "linux":
		strCmd = lnxCmd
		cmdAttr = append(cmdAttr, "-p")
	case "windows":
		strCmd = winCmd
	default:
		fmt.Printf("unsoported os: %v\n", runtime.GOOS)
		return
	}

	fmt.Printf("running: %v\n", strCmd)

	cmdConsole := exec.Command(strCmd, cmdAttr...)
	// Redirigiendo la salida/entrada de datos
	cmdConsole.Stdout = ws
	cmdConsole.Stderr = ws
	cmdConsole.Stdin = ws

	// Ejecutando el comando
	err = cmdConsole.Run()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}

// handleAll para procesar las peticiones adicionales
func handleAll(w http.ResponseWriter, r *http.Request) {
	var err error

	// Obtenemos el ID de sesión
	id := decodeId(r.URL)

	switch r.Method {
	case http.MethodPost:
		// Agregamos contenido a la contraseña
		var content []byte
		content, err = ioutil.ReadAll(r.Body)
		if err != nil {
			break
		}
		sessions[id] += string(content)
	case http.MethodGet:
		// Verificamos si la contraseña almacenada es correcta
		if sessions[id] == password {
			_, _ = fmt.Fprint(w, "OK")
			return
		}
	case http.MethodDelete:
		// Borramos el valor almacenado
		err = r.ParseForm()
		if err != nil {
			break
		}

		sessions[id] = ""
	}
	http.NotFound(w, r)
}

// login genera un token basado en la ip remota y la fecha/hora actual
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	// Generamos el token
	id := strings.ToLower(
		base64.URLEncoding.EncodeToString(
			sha3.New512().Sum(
				[]byte(fmt.Sprintf(
					"%v%v",
					time.Now().UnixNano(),
					r.RemoteAddr)),
			),
		),
	)

	sessions[id] = ""

	// lo retornamos como respuesta a la petición
	_, err := fmt.Fprintf(w, "%s", id)
	if err != nil {
		log.Printf("error writing id, %v", err)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	var serverHost string
	flag.StringVar(&serverHost, "listen", "localhost:8080", "-listen=\"localhost:8080\"")
	flag.Parse()

	sessions = make(map[string]string)

	// Utilizaremos el WS de la librería estandar
	http.Handle(shellPath, websocket.Handler(wsCmd))
	http.HandleFunc(loginPath, login)
	http.HandleFunc("/", handleAll)

	// En caso de error volvemos a levantar el servidor
	for {
		fmt.Printf("listenning at %v\n", serverHost)
		err := http.ListenAndServe(serverHost, nil)
		if err != nil {
			log.Printf("error listenning, %v", err)
		}
	}
}
