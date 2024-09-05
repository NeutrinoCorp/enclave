package security

import "github.com/neutrinocorp/geck/internal/converter"

type PrincipalConverterFunc[T any] converter.ConvertSafeFunc[T, Principal]
