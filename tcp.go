package proof

import (
	"io"
	"net"
	"time"
)

// TCP creates a fetcher that retrieves data over a TCP connection.
// `address` is the TCP address of the server in the form "host:port".
// `request` is an optional string sent to the server upon connection.
func TCP(addr net.Addr, request string) Fetcher {
	return func(p *Proof) error {
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			return err
		}
		defer conn.Close()

		// Set a deadline for the operation
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		// If a request is specified, send it.
		if request != "" {
			if _, err := conn.Write([]byte(request)); err != nil {
				return err
			}
		}

		// Read the response
		response, err := io.ReadAll(conn)
		if err != nil {
			return err
		}

		// Assign the response to the Proof's Content
		p.Content = response
		return nil
	}
}
