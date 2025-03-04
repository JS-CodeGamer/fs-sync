package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/js-codegamer/fs-sync/config"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func GenerateTokenPair(username string, email string) (*TokenPair, error) {
	cfg := config.GetConfig().Auth

	accessToken, err := generateToken(username, email, AccessToken, cfg.JWTSecret, cfg.AccessTokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(username, "", RefreshToken, cfg.JWTSecret, cfg.RefreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(username string, email string, tokenType TokenType, secret string, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration)
	claims := &Claims{
		Username:  username,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "drive-clone",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	cfg := config.GetConfig().Auth
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func RefreshTokens(refreshToken string, email string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshToken, RefreshToken)
	if err != nil {
		return nil, err
	}

	return GenerateTokenPair(claims.Username, email)
}
