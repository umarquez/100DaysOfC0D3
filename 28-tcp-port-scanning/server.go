package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func startServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Printf("[SERVER] error starting listener, %v\n", err)
		return
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Printf("[SERVER] error closing listener, %v\n", err)
		}
	}()

	log.Printf("- [SERVER] Starting echo server at %v\n", port)

	for {
		cnn, err := listener.Accept()
		if err != nil {
			log.Printf("[SERVER] error managing new connection, %v\n", err)
			continue
		}

		go func(c net.Conn) {
			_, err := io.Copy(c, c)
			if err != nil {
				log.Printf("[SERVER] error writing data, %v\n", err)
				return
			}

			_ = c.Close()
		}(cnn)
	}
}
