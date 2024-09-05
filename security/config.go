package security

type ConfigEncryptor struct {
	SecretKey string `env:"ENCRYPTOR_SECRET_KEY,unset" envDefault:"Some_Page_Token_Key_1927_!@#$*~<"`
}

type CognitoConfig struct {
	Region     string `env:"COGNITO_REGION" envDefault:"us-east-1"`
	UserPoolID string `env:"COGNITO_USER_POOL_ID"`
}
