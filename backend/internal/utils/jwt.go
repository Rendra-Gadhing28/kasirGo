package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"kasirpro/internal/config"
)

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	OutletID  uint   `json:"outlet_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	TokenType string `json:"token_type"` // access or refresh
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID, outletID uint, role, email, name string) (string, string, error) {
	cfg := config.AppConfig
	secret := []byte(cfg.JWTSecret)

	// Access Token
	accessClaims := JWTClaims{
		UserID:    userID,
		OutletID:  outletID,
		Role:      role,
		Email:     email,
		Name:      name,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTAccessExpireMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.AppName,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := JWTClaims{
		UserID:    userID,
		OutletID:  outletID,
		Role:      role,
		Email:     email,
		Name:      name,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTRefreshExpireDays*24) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.AppName,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.AppConfig
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
