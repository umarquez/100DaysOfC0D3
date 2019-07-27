/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [19/100] Ping | ICMP Echo
*/

package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Preparamos la escucha de paquetes ICMP en todas las direcciones
	cnn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		panic(err)
	}
	// Cerramos la conexión al finalizar
	defer func() {
		err := cnn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	/*
	 * Ping operates by sending Internet Control Message Protocol (ICMP) echo
	 * request packets to the target host and waiting for an ICMP echo reply.
	 * https://en.wikipedia.org/wiki/Ping_(networking_utility)
	 */
	echoMsj := icmp.Message{
		Type: ipv4.ICMPTypeEcho, // Mensaje de tipo echo
		Code: 0,                 // No existen múltiples códigos para los mensajes echo
		/*
		 * The Identifier and Sequence Number can be used by the client to
		 * match the reply with the request that caused the reply. In
		 * practice, most Linux systems use a unique identifier for every
		 * ping process, and sequence number is an increasing number within
		 * that process. Windows uses a fixed identifier, which varies
		 * between Windows versions, and a sequence number that is only
		 * reset at boot time.
		 * https://en.wikipedia.org/wiki/Ping_(networking_utility)#Echo_request
		 */
		Body: &icmp.Echo{
			ID:  os.Getpid() & 0xffff, // primeros 32 bits del PID
			Seq: 1,                    // Podríamos utilizarlo como consecutivo del mensaje
			/* The payload of the packet is generally filled with
			 * ASCII characters
			 * https://en.wikipedia.org/wiki/Ping_(networking_utility)#Payload
			 */
			Data: []byte("#100DaysOfC0D3"),
		},
	}

	// Obtenemos la versión en binario del mensaje
	bMsj, err := echoMsj.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Haremos ping a host local (127.0.0.1), pero podría ser a cualquier IP
	target := net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
	//target :=  net.IPAddr{IP: net.IPv4(4,2,2,2)} // Google

	// Para esta prueba enviaremos 5 paquetes
	for i := 0; i < 5; i++ {
		t := time.Now() // Vamos a contar el tiempo que toma en responder

		// Enviamos el paquete a travéz de la conexión establecida
		_, err := cnn.WriteTo(bMsj, &target)
		if err != nil {
			log.Fatal(err)
		}

		// Preparamos el buffer que recibirá la respuesta
		bResp := make([]byte, 1500)
		// Y esperamos los datos
		n, origin, err := cnn.ReadFrom(bResp)
		if err != nil {
			log.Fatal(err)
		}

		// Obtenemos el tiempo de respuesta
		delta := time.Since(t).Seconds()

		// Decodificamos la respuesta
		// 1 = ICMP IPv4 Protocol
		resp, err := icmp.ParseMessage(1, bResp[:n])
		if err != nil {
			log.Fatal(err)
		}

		switch resp.Type {
		case ipv4.ICMPTypeEchoReply:
			// recuperamos el mesaje
			r, _ := resp.Body.Marshal(0)

			// decodificamos el contenido
			msj, err := icmp.ParseMessage(1, r)
			if err != nil {
				log.Fatal(err)
			}

			// Imprimimos los resultados
			fmt.Printf("%s (%v seg.): %s\n", origin, delta, msj.Body)
		default:
			fmt.Printf("%t\n", resp)
		}

		// Nos tomamos un segundo
		time.Sleep(time.Second)
	}
}
