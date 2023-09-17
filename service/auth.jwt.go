package service

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SawitProRecruitment/UserService/repository"
)

type JWTAuthService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type JWTAuthServiceOptions struct {
	PrivateKey string
	PublicKey  string
}

func NewJWTAuthService(opts JWTAuthServiceOptions) *JWTAuthService {
	pem, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(opts.PrivateKey))
	if err != nil {
		panic(err)
	}

	cert, err := jwt.ParseRSAPublicKeyFromPEM([]byte(opts.PublicKey))
	if err != nil {
		panic(err)
	}

	return &JWTAuthService{
		privateKey: pem,
		publicKey:  cert,
	}
}

func (s *JWTAuthService) CreateToken(usr *repository.User, now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(usr.ID, 10),
		"name": usr.FullName,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JWTAuthService) ValidateToken(token string) (int64, error) {
	jot, err := s.getJWT(token)
	if err != nil {
		return 0, err
	}

	if jot == nil {
		return 0, fmt.Errorf("token is not valid")
	}

	if !jot.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	return s.getClaimedID(jot)
}

func (s *JWTAuthService) getJWT(token string) (*jwt.Token, error) {
	jot, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return jot, nil
}

func (s *JWTAuthService) getClaimedID(jot *jwt.Token) (int64, error) {
	var id int64 = 0
	claims, ok := jot.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("token is not valid")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("token is not valid")
	}

	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
