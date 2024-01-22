package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	startServer("8082")

}

func startServer(port string) {
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "Listening for connection to %s:%s...\n", "127.0.0.1", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
		fmt.Println()
	}
}

func handleConnection(conn net.Conn) {
	fmt.Fprintf(os.Stdout, "Received request from %s\n", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	var path string
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if s == "\r\n" {
			break
		}
		token := strings.Split(s, " ")
		if token[0] == "GET" {
			path = token[1]
		}

		fmt.Fprint(os.Stdout, s)
	}

	buf := handleRoute(path)
	conn.Write(buf.Bytes())
	conn.Close()
}

func handleRoute(path string) *bytes.Buffer {
	buf := bytes.Buffer{}
	switch path {

	case "/health":
		buf.Write([]byte("HTTP/1.1 204 No Content\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	case "/":
		buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("Content-Length: 27\r\n"))
		buf.Write([]byte("\r\n"))
		buf.Write([]byte("Hello From Backend Server\r\n"))
	default:
		buf.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	}
	return &buf
}
