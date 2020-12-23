package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	sendMsg(conn, reader)
}

func sendMsg(conn net.Conn, reader *bufio.Reader) {
	fmt.Print("Say something: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	fmt.Fprintln(conn, text)

	// receiveMsg(conn)
	sendMsg(conn, reader)
}

func receiveMsg(conn net.Conn) {
	bs, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(bs))
}
