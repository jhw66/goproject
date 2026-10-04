// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type AuthConfig struct {
	AccessSecret    string
	AccessExpire    int64
	ChallengeSecret string
	ChallengeExpire int64
}

type TotpConfig struct {
	Issuer          string
	SecretCipherKey string
	QRCodeSize      int
}

type Config struct {
	rest.RestConf
	DataSource string
	Auth       AuthConfig
	Totp       TotpConfig
}
