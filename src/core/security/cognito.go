package security

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neutrinocorp/nolan/collection/set"
	"github.com/samber/lo"
)

func NewAmazonCognitoKeysJWK(cfg CognitoConfig) (keyfunc.Keyfunc, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	return keyfunc.NewDefaultCtx(ctx,
		[]string{
			"https://cognito-idp." + cfg.Region + ".amazonaws.com/" + cfg.UserPoolID + "/.well-known/jwks.json",
		})
}

func NewPrincipalFactoryCognito() PrincipalFactoryTemplate[*jwt.Token] {
	return PrincipalFactoryTemplate[*jwt.Token]{
		PrincipalConverterFunc: ConvertCognitoJWTToPrincipal,
	}
}

// converters

var _ PrincipalConverterFunc[*jwt.Token] = ConvertCognitoJWTToPrincipal

func ConvertCognitoJWTToPrincipal(token *jwt.Token) (Principal, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot cast jwt claims")
	}
	sub, _ := claims.GetSubject()
	username, _ := claims["username"].(string)
	scopesRaw, _ := claims["scope"].(string)
	scopes := strings.Split(scopesRaw, " ")

	groupsRaw, _ := claims["cognito:groups"].([]any)
	groups := lo.Map(groupsRaw, func(src any, index int) string {
		return fmt.Sprintf("%v", src)
	})

	authoritySet := set.HashSet[string]{}
	authoritySet.AddSlice(scopes...)
	authoritySet.AddSlice(groups...)
	return BasicPrincipal{
		Identifier:   username,
		Subject:      sub,
		User:         username,
		AuthoritySet: authoritySet,
	}, nil
}
