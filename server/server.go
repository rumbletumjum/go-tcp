package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var port = flag.Int("port", 6969, "Listen port; default is 6969")
var addr = flag.String("addr", "", "Listen address; default is \"\" (all interfaces)")

func main() {
	flag.Parse()

	log.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatalf("Error binding to port: %d on %s\n", port, *addr)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting incoming connection")
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Client connected from %s\n", remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		handleMessage(scanner.Text(), conn)
	}
	log.Println("Client at" + remoteAddr + " Disconnected")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := fmt.Sprintf("The time is now: %s\n", time.Now())
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		default:
			conn.Write([]byte("Unrecognized command\n"))

		}
	}
}
