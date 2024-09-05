package securityfx

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/security"
	"github.com/neutrinocorp/geck/security/encryption"
)

var EncryptorAESModule = fx.Module("encryptor_aes",
	fx.Provide(
		env.ParseAs[security.ConfigEncryptor],
		fx.Annotate(
			encryption.NewEncryptorAES,
			fx.As(new(encryption.Encryptor)),
		),
	),
)

var CognitoModule = fx.Module("amazon_cognito",
	fx.Provide(
		env.ParseAs[security.CognitoConfig],
		fx.Annotate(
			security.NewAmazonCognitoKeysJWK,
			fx.As(new(keyfunc.Keyfunc)),
		),
		fx.Annotate(
			security.NewPrincipalFactoryCognito,
			fx.As(new(security.PrincipalFactory[*jwt.Token])),
		),
	),
)
