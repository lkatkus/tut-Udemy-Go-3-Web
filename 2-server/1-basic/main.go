package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Unable to start server", err)
	}
	defer li.Close()

	log.Println("Server running on localhost:8080")

	for {
		conn, err := li.Accept()

		if err != nil {
			log.Println("Error:", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	// if err != nil {
	// 	log.Println("Connection timeout")
	// }

	scanner := bufio.NewScanner(conn)

	fmt.Fprint(conn, "*** Welcome to ZeWeb ***")

	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println("User:", ln)
	}

	defer conn.Close()
	log.Println("Connection closed")
}
