package proof

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

// Gemini creates a Fetcher that retrieves proof data over the Gemini protocol
func Gemini(geminiURL url.URL, conn net.Conn) Fetcher {
	return func(p *Proof) error {
		var err error

		if geminiURL.Scheme != "gemini" {
			return fmt.Errorf("invalid URL scheme: %s", geminiURL.Scheme)
		}

		if conn == nil {

			conn, err = tls.Dial("tcp", geminiURL.Host+":1965", &tls.Config{})
			if err != nil {
				return fmt.Errorf("failed to establish TLS connection: %w", err)
			}
		}
		defer conn.Close()

		_, err = conn.Write([]byte(geminiURL.String() + "\r\n"))
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		reader := bufio.NewReader(conn)
		header, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response header: %w", err)
		}

		statusCode := header[:2]
		if statusCode < "20" || statusCode > "29" {
			return fmt.Errorf("unsuccessful Gemini response: %s", strings.TrimSpace(header))
		}

		response, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		p.Content = response
		return nil
	}
}
