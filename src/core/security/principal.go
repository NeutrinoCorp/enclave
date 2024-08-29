package security

import "github.com/neutrinocorp/nolan/collection/set"

type Principal interface {
	ID() string
	Sub() string
	Username() string
	Authorities() set.Set[string]
}

type BasicPrincipal struct {
	Identifier   string
	Subject      string
	User         string
	AuthoritySet set.Set[string]
}

var _ Principal = (*BasicPrincipal)(nil)

func (b BasicPrincipal) ID() string {
	return b.Identifier
}

func (b BasicPrincipal) Sub() string {
	return b.Subject
}

func (b BasicPrincipal) Username() string {
	return b.User
}

func (b BasicPrincipal) Authorities() set.Set[string] {
	return b.AuthoritySet
}
