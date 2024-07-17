package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func handleClient(conn net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	serveConnection(conn)
}

func main(){
	port := 9090

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %d: %v\n", port, err)
		return
	}
	defer listener.Close()

	fmt.Fprintf(os.Stdout, "TCP server running on port %d\n", port)

	var wg sync.WaitGroup

	stopch := make(chan os.Signal, 1)
	signal.Notify(stopch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func(){
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-done:
					os.Exit(1)
					return 

				default:
					fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
					continue
				}
			}
	
			fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())
			
			wg.Add(1)
			go handleClient(conn, &wg)
		}
	}()

	<-stopch
	fmt.Println("Server stopped")
	close(done)
	listener.Close()
	wg.Wait()
}