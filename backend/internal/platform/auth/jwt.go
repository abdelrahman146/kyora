package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	bearerPrefix = "Bearer "
)

type CustomClaims struct {
	UserID      string `json:"userId"`
	WorkspaceID string `json:"workspaceId"`
	AuthVersion int    `json:"authVersion"`
	jwt.RegisteredClaims
}

func NewJwtToken(userID string, workspaceID string, authVersion int) (string, error) {
	expiry := viper.GetInt(config.JWTExpirySeconds) // in seconds
	jwtExpiry := time.Hour * 24
	if expiry > 0 {
		jwtExpiry = time.Duration(expiry) * time.Second
	}
	secret := viper.GetString(config.JWTSecret)
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:      userID,
		WorkspaceID: workspaceID,
		AuthVersion: authVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.KsuidWithPrefix("jwt"),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    viper.GetString(config.JWTIssuer),
			Audience:  jwt.ClaimStrings{viper.GetString(config.JWTAudience)},
			Subject:   userID,
		},
	})
	return token.SignedString([]byte(secret))
}

func ParseJwtToken(tokenString string) (*CustomClaims, error) {
	secret := viper.GetString(config.JWTSecret)
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is not configured")
	}

	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, bearerPrefix), &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// JwtFromContext extracts the JWT token from the Authorization header.
// The Authorization header must be in the format: "Bearer <token>"
func JwtFromContext(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader
	}
	return ""
}
