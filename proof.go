// package proof provides methods for fetching and verifying signed proofs
// over the Internet.
package proof

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// type Fetcher is a function that is triggered by a Proof's Fetch method.
// It's used to retrieve signed proof content
type Fetcher func(*Proof) error

// type Sigtype represents the type of signature used in a Proof.
// Currently, only PGP is supported.
type Sigtype int

const (
	PGP Sigtype = iota
)

// type Proof represents a signed proof that can be fetched and verified.
type Proof struct {
	URL         url.URL       // URL for proof data
	Kind        Sigtype       // type of signature
	Content     []byte        // signed proof content
	Pubkey      []byte        // public key used to verify the proof
	LastChecked time.Time     // last time the proof was checked
	Verified    bool          // whether the proof is verified
	MaxAge      time.Duration // maximum age since last verification
	FetchFunc   Fetcher
	mu          sync.RWMutex
}

// method Verify checks the signature of a Proof and returns an error if it's invalid.
func (p *Proof) Verify() error {
	switch p.Kind {
	case PGP:
		return p.verifyPGP()
	default:
		return fmt.Errorf("unsupported signature type")
	}
}

func (p *Proof) verifyPGP() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Pubkey) == 0 {
		return fmt.Errorf("missing public key")
	}

	if len(p.Content) == 0 {
		return fmt.Errorf("missing content")
	}

	_, err := helper.VerifyCleartextMessageArmored(
		string(p.Pubkey),
		string(p.Content),
		crypto.GetUnixTime(),
	)
	p.Verified = err == nil
	p.LastChecked = time.Now()
	return err
}

// method IsExpired returns true if a Proof is older than its MaxAge.
func (p *Proof) IsExpired() bool {
	return time.Now().After(p.LastChecked.Add(p.MaxAge))
}

// method Fetch retrieves the signed proof content
func (p *Proof) Fetch() error {
	if p.FetchFunc == nil {
		return fmt.Errorf("missing fetch function")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	return p.FetchFunc(p)
}

// method Check verifies a Proof and fetches it if it's expired.
func (p *Proof) Check() error {
	if p.IsExpired() {
		if err := p.Fetch(); err != nil {
			return err
		}
	}
	return p.Verify()
}
