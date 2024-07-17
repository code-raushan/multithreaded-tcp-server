package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	WAIT_FOR_MSG = iota
	IN_MSG
)

func serveConnection(conn net.Conn) {
	defer conn.Close()

	// sending acknowledgment to the client
	if _, err := conn.Write([]byte("*")); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to client: %v\n", err)
		return
	}

	state := WAIT_FOR_MSG

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // Client closed the connection
			}
			fmt.Fprintf(os.Stderr, "Error reading from client: %v\n", err)
			return
		}

		for i := 0; i < n; i++ {
			switch state {
			case WAIT_FOR_MSG:
				if buf[i] == '^' {
					state = IN_MSG
					fmt.Fprintf(os.Stdout, "In-Message State\n")
				}
			case IN_MSG:
				if buf[i] == '$' {
					state = WAIT_FOR_MSG
					fmt.Fprintf(os.Stdout, "Wait-For-Message State\n")

				}else{
					buf[i]++
					if _, err := conn.Write(buf[i : i+1]); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing to client: %v\n", err)
						return
					}
				}
			}
		}
	}
}

func main(){
	
}