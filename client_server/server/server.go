package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mutex   sync.Mutex
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Meminta nama dari client
	conn.Write([]byte("Enter your name: "))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	name = strings.TrimSpace(name)

	// Menambahkan klien ke peta
	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	fmt.Printf("%s connected.\n", name)

	// Membaca pesan dari client dan mengirimkan ke client lain
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("%s disconnected.\n", name)
			removeClient(conn)
			return
		}

		trimmedMsg := strings.TrimSpace(msg)
		if trimmedMsg != "" {
			fmt.Printf("Received from %s: %s\n", name, trimmedMsg)
			sendMessageToOtherClients(conn, fmt.Sprintf("%s: %s\n", name, trimmedMsg))
		}
	}
}

// Mengirim pesan ke semua klien lain, kecuali pengirim
func sendMessageToOtherClients(sender net.Conn, msg string) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range clients {
		if conn != sender {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error sending message to", clients[conn], ":", err)
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}

// Hapus client dari peta jika terputus
func removeClient(conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(clients, conn)
}

func main() {
	// Mulai server TCP pada port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on port 8080")

	for {
		// Menerima koneksi baru dari klien
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
