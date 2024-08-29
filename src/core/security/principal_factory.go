package security

import (
	"context"
)

type PrincipalFactory[T any] interface {
	NewContextWithPrincipal(parent context.Context, args T) (context.Context, error)
}

// PrincipalFactoryTemplate is the default implementation of PrincipalFactory interface.
type PrincipalFactoryTemplate[T any] struct {
	PrincipalConverterFunc PrincipalConverterFunc[T]
}

var _ PrincipalFactory[string] = (*PrincipalFactoryTemplate[string])(nil)

func (p PrincipalFactoryTemplate[T]) NewContextWithPrincipal(parent context.Context, args T) (context.Context, error) {
	principal, err := p.PrincipalConverterFunc(args)
	if err != nil {
		return nil, err
	}
	return context.WithValue(parent, PrincipalContextKey, principal), nil
}
