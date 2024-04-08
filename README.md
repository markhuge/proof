# proof

package proof provides methods for fetching and verifying signed proofs
over the Internet.

## Types

### type [Fetcher](/proof.go#L17)

`type Fetcher func(*Proof) error`

type Fetcher is a function that is triggered by a Proof's Fetch method.
It's used to retrieve signed proof content

#### func [Gemini](/gemini.go#L14)

`func Gemini(geminiURL url.URL, conn net.Conn) Fetcher`

Gemini creates a Fetcher that retrieves proof data over the Gemini protocol

#### func [Gopher](/gopher.go#L11)

`func Gopher(u url.URL, conn net.Conn) Fetcher`

Gopher creates a fetcher that retrieves proof data over the Gopher protocol

#### func [HTTP](/http.go#L9)

`func HTTP(client *http.Client) Fetcher`

#### func [TCP](/tcp.go#L11)

`func TCP(addr net.Addr, request string) Fetcher`

TCP is a Fetcher that connects to a TCP server and reads the response.
This is more of an example since there is no handling of the response.

### type [Proof](/proof.go#L28)

`type Proof struct { ... }`

type Proof represents a signed proof that can be fetched and verified.

#### func (*Proof) [Check](/proof.go#L90)

`func (p *Proof) Check() error`

method Check verifies a Proof and fetches it if it's expired.

#### func (*Proof) [Fetch](/proof.go#L78)

`func (p *Proof) Fetch() error`

method Fetch retrieves the signed proof content

#### func (*Proof) [IsExpired](/proof.go#L73)

`func (p *Proof) IsExpired() bool`

method IsExpired returns true if a Proof is older than its MaxAge.

#### func (*Proof) [Verify](/proof.go#L41)

`func (p *Proof) Verify() error`

method Verify checks the signature of a Proof and returns an error if it's invalid.

### type [Sigtype](/proof.go#L21)

`type Sigtype int`

type Sigtype represents the type of signature used in a Proof.
Currently, only PGP is supported.
