package proof_test

import (
	"bytes"
	"fmt"
	"net/url"
	"testing"

	"markhuge.com/proof"
)

func TestGeminiFetcher(t *testing.T) {
	tt := []struct {
		name        string
		response    string
		wantErr     bool
		wantContent string
	}{
		{
			name:        "Successful fetch",
			response:    "20 text/gemini\r\nThis is a test content.",
			wantErr:     false,
			wantContent: "This is a test content.",
		},
		{
			name:        "Error from server",
			response:    "51 Not found\r\n",
			wantErr:     true,
			wantContent: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockConn := &MockConn{
				Buffer: bytes.NewBufferString(tc.response),
			}

			geminiURL, _ := url.Parse("gemini://example.com")

			fetcher := proof.Gemini(*geminiURL, mockConn)

			p := &proof.Proof{}
			err := fetcher(p)

			if (err != nil) != tc.wantErr {
				fmt.Println(err)
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err != nil)
			}

			if string(p.Content) != tc.wantContent {
				t.Errorf("Expected content: %s, got: %s", tc.wantContent, string(p.Content))
			}
		})
	}
}
