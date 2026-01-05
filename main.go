package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linech := make(chan string)

	go func() {
		defer func() {
			f.Close()
			close(linech)
		}()

		buf := make([]byte, 8)
		curline := ""
	outer:
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Print("error reading file")
				}
				break outer
			}

			read := string(buf[:n])
			parts := strings.Split(read, "\n")
			if len(parts) == 0 {
				fmt.Print("error there should more than 0 parts")
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
	file, err := os.Open("./messages.txt")
	if err != nil {
		os.Exit(1)
	}

	linech := getLinesChannel(file)

	for line := range linech {
		fmt.Printf("read: %s\n", line)
	}
}
