package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linech := make(chan string)

	go func() {
		defer func() {
			f.Close()
			close(linech)
			fmt.Println("Connection is closed successfully")
		}()

		buf := make([]byte, 8)
		curline := ""
	outer:
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error: reading from io")
				}
				break outer
			}

			read := string(buf[:n])
			parts := strings.Split(read, "\n")
			if len(parts) == 0 {
				fmt.Println("Error: expect more than 0 parts")
				break outer
			}

			// haven't found the newline
			if len(parts) == 1 {
				curline += parts[0]
				continue outer
			}

			// There are multiple parts
			// send the current line with the first part
			numparts := len(parts) - 1
			curline += parts[0]
			linech <- curline

			for i := range numparts {
				if i == numparts-1 {
					curline = parts[i+1]
					continue outer
				}

				linech <- parts[i+1]
			}
		}

		if curline != "" {
			linech <- curline
		}
	}()

	return linech
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Listening on address: %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("Connection %s has been accepted\n", conn.LocalAddr())

		linech := getLinesChannel(conn)

		for line := range linech {
			fmt.Println(line)
		}
	}
}
