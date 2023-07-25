package main

import (
	"log"
	"net"
	"time"
)

func main() {
	// Create a listener
	l, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Listener returned: %s", err)
	}
	defer l.Close()

	for {
		// Accept new connections
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("Unable to accept new connections: %s", err)
		}

		// Create a goroutine that reads and writes-back data
		go func() {
			log.Printf("TCP Session Open")
			// Clean up session when goroutine completes, it's ok to
			// call Close more than once.
			defer c.Close()

			for {
				b := make([]byte, 120)

				// Read from TCP Buffer
				_, err := c.Read(b)
				if err != nil {
					log.Printf("Error reading TCP Session: %s", err)
				}

				// Write-back data to TCP Client
				_, err = c.Write(b)
				if err != nil {
					log.Printf("Error writing TCP Session: %s", err)
				}
				time.Sleep(1 * time.Second)
			}
		}()

		// Create a goroutine that closes a session after 15 seconds
		go func() {
			<-time.After(time.Duration(150) * time.Second)
			defer c.Close()
		}()
	}
}
