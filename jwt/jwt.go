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
		GenerateToken(userID, roleID string, permissions []string, duration time.Duration) (string, error)
		GenerateImpersonateToken(userID, roleID, originalUserID string, permissions []string, duration time.Duration) (string, error)
		ValidateToken(tokenString string) (*CustomClaims, error)
	}

	CustomClaims struct {
		UserID          string   `json:"user_id"`
		RoleID          string   `json:"role_id"`
		IsImpersonating bool     `json:"is_impersonating"`
		OriginalUserID  string   `json:"original_user_id,omitempty"`
		Permissions     []string `json:"permissions"`

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

func (j *JWT) GenerateToken(userID, roleID string, permissions []string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) GenerateImpersonateToken(userID, roleID, originalUserID string, permissions []string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID:          userID,
		RoleID:          roleID,
		IsImpersonating: true,
		OriginalUserID:  originalUserID,
		Permissions:     permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			ID:        uuid.NewString(),
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

func (j *JWT) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, j.parseToken)
	if err != nil {
		return nil, dto.ErrValidateToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, dto.ErrTokenInvalid
	}

	return claims, nil
}
