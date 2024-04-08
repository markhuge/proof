package proof_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"markhuge.com/proof"
)

func TestHTTPSuccess(t *testing.T) {
	pubKey, signedMsg, err := sign("This is a test message.")
	if err != nil {
		t.Fatalf("Failed to generate signed message: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, signedMsg)
	}))
	defer mockServer.Close()

	p := proof.Proof{
		URL:       *mustParseURL(mockServer.URL),
		FetchFunc: proof.HTTP(http.DefaultClient),
		Pubkey:    []byte(pubKey),
	}

	// Fetch the content from the mock server
	if err := p.Fetch(); err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	// Verify the fetched content
	if err := p.Verify(); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}
