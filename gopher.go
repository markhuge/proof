package proof

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

// Gopher creates a fetcher that retrieves proof data over the Gopher protocol
func Gopher(u url.URL, conn net.Conn) Fetcher {
	return func(p *Proof) error {
		var err error

		if u.Scheme != "gopher" {
			return fmt.Errorf("invalid URL scheme: %s", u.Scheme)
		}

		if conn == nil {
			conn, err = net.Dial("tcp", u.Host)
			if err != nil {
				return fmt.Errorf("failed to connect to Gopher server: %w", err)
			}
		}
		defer conn.Close()

		// Send the selector; Gopher protocol specifies a selector string followed by CRLF
		// If the Path is empty, use "/" as the default selector
		selector := u.Path
		if selector == "" {
			selector = "/"
		}
		_, err = conn.Write([]byte(selector + "\r\n"))
		if err != nil {
			return fmt.Errorf("failed to send selector: %w", err)
		}

		response, err := io.ReadAll(conn)
		if err != nil {
			return fmt.Errorf("failed to read response from Gopher server: %w", err)
		}

		p.Content = response
		return nil
	}
}
