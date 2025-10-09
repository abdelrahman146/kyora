package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const (
	jwtCookieName = "jwt"
	bearerPrefix  = "Bearer "
)

type jwtHelper struct{}

type CustomClaims struct {
	UserID string `json:"userId"`
	OrgID  string `json:"orgId"`
	jwt.RegisteredClaims
}

func (jwtHelper) GenerateToken(userID string, orgID string) (string, error) {
	expiry := viper.GetInt("jwt.expiry") // in seconds
	jwtExpiry := time.Hour * 24
	if expiry > 0 {
		jwtExpiry = time.Duration(expiry) * time.Second
	}
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    viper.GetString("app.name"),
			Subject:   userID,
		},
	})
	return token.SignedString([]byte(secret))
}

func (jwtHelper) ParseToken(tokenString string) (*CustomClaims, error) {
	secret := viper.GetString("jwt.secret")
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

func (jwtHelper) SetJwtCookie(c *gin.Context, token string) {
	jwtExpiry := viper.GetInt("jwt.expiry") // in seconds
	domain := viper.GetString("app.domain")
	if jwtExpiry <= 0 {
		jwtExpiry = 3600 // default to 1 hour
	}
	c.SetCookie(jwtCookieName, token, jwtExpiry, "/", domain, false, true)
}

func (jwtHelper) ClearJwtCookie(c *gin.Context) {
	domain := viper.GetString("app.domain")
	c.SetCookie(jwtCookieName, "", -1, "/", domain, false, true)
}

func (jwtHelper) GetJwtFromContext(c *gin.Context) string {
	jwtToken, err := c.Cookie(jwtCookieName)
	if err != nil {
		return ""
	}
	return jwtToken
}

var JWT = jwtHelper{}
