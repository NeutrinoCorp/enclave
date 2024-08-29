package security

type PrincipalConverterFunc[T any] func(args T) (Principal, error)
