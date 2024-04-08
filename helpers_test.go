package proof_test

import (
	"bytes"
	"net"
	"net/url"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/pkg/errors"
	"markhuge.com/proof"
)

var (
	name    = "Test User"
	email   = "test@user.com"
	pass    = []byte("pass")
	rsaBits = 2048
)

func sign(msg string) (pubkey, signed string, err error) {
	key, err := crypto.GenerateKey(name, email, "rsa", rsaBits)
	if err != nil {
		return "", "", err
	}

	pubkey, err = key.GetArmoredPublicKey()
	if err != nil {
		return "", "", err
	}

	defer key.ClearPrivateParams()

	locked, err := key.Lock(pass)
	if err != nil {
		return "", "", errors.Wrap(err, "gopenpgp: unable to lock new key")
	}

	privkey, err := locked.Armor()
	if err != nil {
		return "", "", errors.Wrap(err, "gopenpgp: unable to armor new key")
	}

	signed, err = helper.SignCleartextMessageArmored(privkey, pass, msg)
	if err != nil {
		return "", "", errors.Wrap(err, "gopenpgp: unable to armor new key")
	}

	return pubkey, signed, nil

}

// Helper function to parse URLs safely in tests
func mustParseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic("Failed to parse URL: " + err.Error())
	}
	return u
}

// InitMock initializes a mock Fetcher function with pre-signed content.
func InitMock(msg string) (string, string, func(error) proof.Fetcher, error) {
	pubKey, signedMsg, err := sign(msg)
	if err != nil {
		return "", "", nil, err
	}

	// Return a function that creates a new Fetcher.
	return pubKey, signedMsg, func(err error) proof.Fetcher {
		return func(p *proof.Proof) error {
			if err != nil {
				return err
			}
			p.Content = []byte(signedMsg)
			p.Pubkey = []byte(pubKey)
			return nil
		}
	}, nil
}

// MockConn implements net.Conn interface for testing.
type MockConn struct {
	net.Conn
	Buffer *bytes.Buffer
}

func (mc *MockConn) Read(b []byte) (n int, err error) {
	return mc.Buffer.Read(b)
}

func (mc *MockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (mc *MockConn) Close() error {
	return nil
}
