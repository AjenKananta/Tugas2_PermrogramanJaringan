package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Menangani penerimaan pesan dari server
func receiveMessages(conn net.Conn, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()

	reader := bufio.NewReader(conn)
	for {
		select {
		case <-stopChan:
			return
		default:
			// Menerima pesan dari server
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server.")
				stopChan <- true
				return
			}
			fmt.Print("\r" + message)
			fmt.Print("Enter message: ")
		}
	}
}

func main() {
	// Koneksi ke server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Sinkronisasi untuk menerima pesan
	var wg sync.WaitGroup
	stopChan := make(chan bool)

	wg.Add(1)
	go receiveMessages(conn, &wg, stopChan)

	// Kirim nama pengguna
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	conn.Write([]byte(name + "\n"))

	fmt.Println("Welcome,", name)
	fmt.Println("You can now start sending messages.")

	// Kirim pesan dari client ke server
	for {
		fmt.Print(name, ": ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		conn.Write([]byte(msg + "\n"))

		if msg == "exit" {
			stopChan <- true
			break
		}
	}

	wg.Wait()
}
