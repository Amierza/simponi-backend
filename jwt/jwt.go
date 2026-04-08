package jwt

import (
	"os"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	IJWT interface {
		GenerateToken(userID, roleID string, duration time.Duration) (string, error)
		ValidateToken(tokenString string) (*jwtCustomClaim, error)
	}

	jwtCustomClaim struct {
		UserID string `json:"user_id"`
		RoleID string `json:"role_id"`
		jwt.RegisteredClaims
	}

	JWT struct {
		secretKey string
		issuer    string
	}
)

func NewJWT() *JWT {
	return &JWT{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		panic("JWT_SECRET is required")
	}

	return secretKey
}

func (j *JWT) GenerateToken(userID, roleID string, duration time.Duration) (string, error) {
	claims := jwtCustomClaim{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration))),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			ID:        uuid.NewString(), // jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, dto.ErrUnexpectedSigningMethod
	}

	return []byte(j.secretKey), nil
}

func (j *JWT) ValidateToken(tokenString string) (*jwtCustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, j.parseToken)
	if err != nil {
		return nil, dto.ErrValidateToken
	}

	claims, ok := token.Claims.(*jwtCustomClaim)
	if !ok || !token.Valid {
		return nil, dto.ErrTokenInvalid
	}

	return claims, nil
}
