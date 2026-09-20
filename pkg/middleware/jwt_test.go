package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestJWTSignAndVerify(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mustMarshalPKCS8(t, privateKey)})
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mustMarshalPublic(t, &privateKey.PublicKey)})

	service, err := NewJWTFromPEM(string(privatePEM), strings.ReplaceAll(string(publicPEM), "\n", `\n`))
	if err != nil {
		t.Fatal(err)
	}
	claims := &jwt.StandardClaims{Subject: "user-1", ExpiresAt: time.Now().Add(time.Hour).Unix()}
	token, err := service.Sign(claims)
	if err != nil {
		t.Fatal(err)
	}
	verified := &jwt.StandardClaims{}
	if err := service.Verify(token, verified); err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if verified.Subject != claims.Subject {
		t.Fatalf("subject = %q, want %q", verified.Subject, claims.Subject)
	}
}

func TestJWTRejectsNonRS256Token(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	service := &jwtGo{privateKey: privateKey, publicKey: &privateKey.PublicKey}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour).Unix()})
	signed, err := token.SignedString([]byte("not-an-rsa-key"))
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Verify(signed, &jwt.StandardClaims{}); err == nil {
		t.Fatal("expected non-RS256 token to be rejected")
	}
}

func TestNewJWTRejectsMismatchedKeys(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	otherKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	publicDER, _ := x509.MarshalPKIXPublicKey(&otherKey.PublicKey)
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})
	_, err := NewJWTFromPEM(string(privatePEM), string(publicPEM))
	if err == nil {
		t.Fatal("expected mismatched keys to be rejected")
	}
}

func mustMarshalPKCS8(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()
	data, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustMarshalPublic(t *testing.T, key *rsa.PublicKey) []byte {
	t.Helper()
	data, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
