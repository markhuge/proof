package proof_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"markhuge.com/proof"
)

func TestProofCheck(t *testing.T) {
	pubkey, content, createFetcher, err := InitMock("This is a test message.")
	if err != nil {
		t.Fatalf("Failed to initialize mock fetcher: %v", err)
	}

	tt := []struct {
		name            string
		setupFetcher    func() proof.Fetcher // Setup function to initialize fetcher
		modifyProof     func(*proof.Proof)   // Function to modify the Proof object before test
		expectedError   bool
		expectedFetched bool
	}{
		{
			name: "Valid and not expired",
			setupFetcher: func() proof.Fetcher {
				return createFetcher(nil) // Fetcher without errors
			},
			modifyProof: func(p *proof.Proof) {
				p.Pubkey = []byte(pubkey)
				p.MaxAge = 10 * time.Minute
				p.Content = []byte(content)
			},
			expectedError:   false,
			expectedFetched: false,
		},
		{
			name: "Valid but expired, requires fetch",
			setupFetcher: func() proof.Fetcher {
				return createFetcher(nil) // Fetcher without errors
			},
			modifyProof: func(p *proof.Proof) {
				p.Pubkey = []byte(pubkey)
				p.LastChecked = time.Now().Add(-2 * time.Hour) // Ensure it's expired
				p.MaxAge = 1 * time.Second
			},
			expectedError:   false,
			expectedFetched: true,
		},
		{
			name: "Fetch error",
			setupFetcher: func() proof.Fetcher {
				return createFetcher(errors.New("simulated fetch error"))
			},
			modifyProof:     func(p *proof.Proof) {}, // No modifications needed
			expectedError:   true,
			expectedFetched: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			p := proof.Proof{
				URL:       *mustParseURL("http://example.com"),
				FetchFunc: tc.setupFetcher(),
				MaxAge:    1 * time.Hour,
			}

			tc.modifyProof(&p)

			err := p.Check()

			if (err != nil) != tc.expectedError {
				fmt.Printf("%+v\n", err)
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err != nil)
			}

		})
	}
}
