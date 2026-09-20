package middleware

import (
	"github.com/golang-jwt/jwt"
)

// JWT describes what a jwt impl is capable of
type JWT interface {
	Encrypt(claims jwt.Claims) (string, error)
	Decrypt(tokenStr string, claims jwt.Claims, skipClaimsValidation bool) error
}

type jwtGo struct {
	jwtSecret []byte
}

// NewJWT returns JWT instance for encrypt/decrypt jwt tokens
func NewJWT(jwtSecret string) JWT {
	return &jwtGo{
		jwtSecret: []byte(jwtSecret),
	}
}

func (uc *jwtGo) Encrypt(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uc.jwtSecret)
}

// Decrypt decrypts and verifys jwt claims
func (uc *jwtGo) Decrypt(tokenStr string, claims jwt.Claims, skipClaimsValidation bool) error {
	parser := jwt.Parser{
		SkipClaimsValidation: skipClaimsValidation,
	}

	if _, err := parser.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return uc.jwtSecret, nil
		},
	); err != nil {
		return err
	}

	return nil
}
