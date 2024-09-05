package security

import (
	"github.com/emirpasic/gods/v2/sets"
)

type Principal interface {
	ID() string
	Sub() string
	Username() string
	Authorities() sets.Set[string]
}

type BasicPrincipal struct {
	Identifier   string
	Subject      string
	User         string
	AuthoritySet sets.Set[string]
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

func (b BasicPrincipal) Authorities() sets.Set[string] {
	return b.AuthoritySet
}
