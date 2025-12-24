package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	bearerPrefix  = "Bearer "
	jwtCookieName = "jwt"
)

type CustomClaims struct {
	UserID      string `json:"userId"`
	WorkspaceID string `json:"workspaceId"`
	jwt.RegisteredClaims
}

func NewJwtToken(userID string, workspaceID string) (string, error) {
	expiry := viper.GetInt(config.JWTExpirySeconds) // in seconds
	jwtExpiry := time.Hour * 24
	if expiry > 0 {
		jwtExpiry = time.Duration(expiry) * time.Second
	}
	secret := viper.GetString(config.JWTSecret)
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:      userID,
		WorkspaceID: workspaceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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

func SetJwtCookie(c *gin.Context, token string) {
	jwtExpiry := viper.GetInt(config.JWTExpirySeconds) // in seconds
	domain := viper.GetString(config.AppDomain)
	if jwtExpiry <= 0 {
		jwtExpiry = 3600 // default to 1 hour
	}
	c.SetCookie(jwtCookieName, token, jwtExpiry, "/", domain, false, true)
}

func ClearJwtCookie(c *gin.Context) {
	domain := viper.GetString(config.AppDomain)
	c.SetCookie(jwtCookieName, "", -1, "/", domain, false, true)
}

func JwtFromContext(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader
	}

	jwtToken, err := c.Cookie(jwtCookieName)
	if err == nil && jwtToken != "" {
		return jwtToken
	}
	return ""
}
