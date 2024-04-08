package proof

import (
	"io"
	"net"
	"time"
)

// TCP is a Fetcher that connects to a TCP server and reads the response.
// This is more of an example since there is no handling of the response.
func TCP(addr net.Addr, request string) Fetcher {
	return func(p *Proof) error {
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			return err
		}
		defer conn.Close()

		// Set a deadline for the operation
		conn.SetDeadline(time.Now().Add(5 * time.Second))

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
