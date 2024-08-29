package security

import (
	"context"

	"github.com/neutrinocorp/nolan/collection/list"

	"github.com/neutrinocorp/geck/systemerror"
)

func HasAnyAuthorities(ctx context.Context, authorities []string) error {
	principal, err := GetPrincipalFromContext(ctx)
	if err != nil {
		return err
	}

	isAuthz := false
	for _, authority := range authorities {
		if ok := principal.Authorities().Contains(authority); ok {
			isAuthz = true
			break
		}
	}
	if !isAuthz {
		return systemerror.NewPermissionDeniedAuthorities(principal.ID(), principal.Authorities().ToSlice(), authorities)
	}
	return nil
}

func HasAuthorities(ctx context.Context, authorities []string) error {
	principal, err := GetPrincipalFromContext(ctx)
	if err != nil {
		return err
	}

	isAuthz := principal.Authorities().ContainsAll(list.NewSliceList(authorities))
	if !isAuthz {
		return systemerror.NewPermissionDeniedAuthorities(principal.ID(), principal.Authorities().ToSlice(), authorities)
	}
	return nil
}

func IsResourceOwner(ctx context.Context, resourceReqOwner string) error {
	principal, err := GetPrincipalFromContext(ctx)
	if err != nil {
		return err
	}

	if isAuthz := principal.ID() == resourceReqOwner; !isAuthz {
		return systemerror.NewPermissionDeniedInvalidOwner(principal.Username(), resourceReqOwner)
	}
	return nil
}

func IsResourceOwnerOrHasAnyAuthorities(ctx context.Context, resourceReqOwner string, authorities []string) error {
	principal, err := GetPrincipalFromContext(ctx)
	if err != nil {
		return err
	}

	isAuthz := false
	for _, authority := range authorities {
		if ok := principal.Authorities().Contains(authority); ok {
			isAuthz = true
			break
		}
	}
	if isAuthz {
		return nil
	}

	if isAuthz = principal.ID() == resourceReqOwner; !isAuthz {
		return systemerror.NewPermissionDeniedInvalidOwner(principal.Username(), resourceReqOwner)
	}
	return nil
}
