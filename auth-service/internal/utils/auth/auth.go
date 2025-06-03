package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenMaker interface {
	GenerateTokens(userID string, email string) (accessToken string, refreshToken string, refreshExpiresAt time.Time, err error)
	VerifyToken(token string) (userID string, email string, err error)
}

type jwtMaker struct {
	secretKey       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewTokenMaker(jwtSecret string) TokenMaker {
	return &jwtMaker{
		secretKey:       jwtSecret,
		accessDuration:  24 * time.Hour,
		refreshDuration: 7 * 24 * time.Hour,
	}
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (j *jwtMaker) GenerateTokens(userID string, email string) (string, string, time.Time, error) {
	now := time.Now()

	// Access Token
	accessClaims := jwtCustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh Token
	refreshExpiresAt := now.Add(j.refreshDuration)
	refreshClaims := jwtCustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, refreshExpiresAt, nil
}

func (j *jwtMaker) VerifyToken(token string) (string, string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := parsedToken.Claims.(*jwtCustomClaims)
	if !ok || !parsedToken.Valid {
		return "", "", err
	}

	return claims.UserID, claims.Email, nil
}
