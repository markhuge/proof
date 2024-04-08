package proof

import (
	"fmt"
	"io"
	"net/http"
)

func HTTP(client *http.Client) Fetcher {
	return func(p *Proof) error {
		res, err := client.Get(p.URL.String())
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return fmt.Errorf("http request failed with status code: %d", res.StatusCode)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		p.Content = b
		return nil
	}
}
