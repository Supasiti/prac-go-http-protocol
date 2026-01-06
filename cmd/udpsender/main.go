package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error: unable to resolve address; %s\n", err)
		os.Exit(1)
	}

	// dial to address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error: unable to connect with address; %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())

	// Read user input
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		read, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Printf("Error: reading user input; %s\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(read))
		if err != nil {
			fmt.Printf("Error: writing user input to UDP; %s\n", err)
			os.Exit(1)
		}
	}
}
