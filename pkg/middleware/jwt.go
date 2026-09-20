package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

// JWT describes what a jwt impl is capable of
type JWT interface {
	Sign(claims jwt.Claims) (string, error)
	Verify(tokenStr string, claims jwt.Claims) error
}

type jwtGo struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWTFromPEM creates an RS256 JWT service from PEM key strings. Values may
// contain real newlines or escaped \n sequences from environment variables.
func NewJWTFromPEM(privateKeyValue, publicKeyValue string) (JWT, error) {
	if privateKeyValue == "" || publicKeyValue == "" {
		return nil, fmt.Errorf("JWT private and public keys are required")
	}
	privatePEM := normalizePEM(privateKeyValue)
	privateKey, err := parsePrivateKey(privatePEM)
	if err != nil {
		return nil, fmt.Errorf("parse JWT private key: %w", err)
	}

	publicPEM := normalizePEM(publicKeyValue)
	publicKey, err := parsePublicKey(publicPEM)
	if err != nil {
		return nil, fmt.Errorf("parse JWT public key: %w", err)
	}
	if privateKey.N.BitLen() < 2048 {
		return nil, fmt.Errorf("JWT private key must be at least 2048 bits")
	}
	if privateKey.PublicKey.N.Cmp(publicKey.N) != 0 || privateKey.PublicKey.E != publicKey.E {
		return nil, fmt.Errorf("JWT private and public keys do not match")
	}

	return &jwtGo{privateKey: privateKey, publicKey: publicKey}, nil
}

func normalizePEM(value string) []byte {
	return []byte(strings.ReplaceAll(strings.TrimSpace(value), `\n`, "\n"))
}

func (uc *jwtGo) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(uc.privateKey)
}

// Verify verifies an RS256 JWT signature and claims.
func (uc *jwtGo) Verify(tokenStr string, claims jwt.Claims) error {
	parser := jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodRS256.Alg()},
	}

	token, err := parser.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return uc.publicKey, nil
		},
	)
	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA")
	}
	return privateKey, nil
}

func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		publicKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA")
		}
		return publicKey, nil
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
