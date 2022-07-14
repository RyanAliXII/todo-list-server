package pasetoutils

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/vk-rv/pvx"
)

type PasetoClaims struct {
	Issuer     string
	Subject    string
	Expiration time.Time
	NotBefore  time.Time
	IssuedAt   time.Time
	tokenId    string
}

func (p *PasetoClaims) NewToken() (string, string, error) {
	key := []byte(os.Getenv("PASETO_SECRET_KEY"))
	symmetricKey := pvx.NewSymmetricKey(key, pvx.Version4)
	tokenId := uuid.NewString()
	claims := pvx.RegisteredClaims{
		Issuer:     p.Issuer,
		Subject:    p.Subject,
		Expiration: &p.Expiration,
		NotBefore:  &p.NotBefore,
		IssuedAt:   &p.IssuedAt,
		TokenID:    tokenId,
	}
	pasetoV4 := pvx.NewPV4Local()
	token, err := pasetoV4.Encrypt(symmetricKey, &claims, pvx.WithAssert([]byte("task_app_server")))
	return token, tokenId, err
}

func NewToken() (string, error) {
	key := []byte(os.Getenv("PASETO_SECRET_KEY"))

	symmetricKey := pvx.NewSymmetricKey(key, pvx.Version4)
	expi := time.Now().Add(1 * time.Minute)
	now := time.Now()
	claims := pvx.RegisteredClaims{
		Issuer:     "Task App Server",
		Expiration: &expi,
		NotBefore:  &now,
		IssuedAt:   &now,
		TokenID:    uuid.NewString(),
	}
	pasetoV4 := pvx.NewPV4Local()
	token, err := pasetoV4.Encrypt(symmetricKey, &claims, pvx.WithAssert([]byte("task_app_server")))
	return token, err
}

func ValidateToken(token string) (pvx.RegisteredClaims, error) {
	pasetoV4 := pvx.NewPV4Local()
	key := []byte(os.Getenv("PASETO_SECRET_KEY"))
	symmetricKey := pvx.NewSymmetricKey(key, pvx.Version4)
	claims := pvx.RegisteredClaims{}
	err := pasetoV4.Decrypt(token, symmetricKey, pvx.WithAssert([]byte("task_app_server"))).ScanClaims(&claims)

	return claims, err
}
