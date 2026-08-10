package middleware

import (
	//entity "cspirt/internal/domain/auth"
	//cacheRepo "cspirt/internal/domain/cache/repo"
	//"context"
	"fmt"
	"net/http"
	"strings"
	//"time"

	entityClaims "ditdah/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, errMessage := accessTokenFromRequest(c)
		if errMessage != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMessage})
			c.Abort()
			return
		}

		claims := &entityClaims.Claims{}

		tok, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !tok.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			message := "invalid or expired token"
			if err != nil {
				message += ": " + err.Error()
			}

			c.Abort()
			return
		}

		if claims.Username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// if isTokenBlacklisted(tokenString) {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
		// 	c.Abort()
		// 	return
		// }

		c.Set("Id", claims.UserID)
		c.Set("Username", claims.Username)
		c.Next()

	}
}

// isTokenBlacklisted reports whether tokenString was revoked via logout.
// Fails open (returns false) on a nil cache or a Redis error so an outage
// never locks every user out - see internal/adapter/redis/README.md.
// func isTokenBlacklisted(tokenString string) bool {
// 	if cache == nil {
// 		return false
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	exists, err := cache.Exists(ctx, entity.BlacklistTokenKey(tokenString))
// 	if err != nil {
// 		logger.WriteSafe(logger.LogEntry{
// 			Level:   "error",
// 			Action:  "auth_middleware",
// 			Message: "redis unavailable, failing open: " + err.Error(),
// 		})
// 		return false
// 	}

// 	return exists
// }

func accessTokenFromRequest(c *gin.Context) (string, string) {
	auth := c.GetHeader("Authorization")
	if auth != "" {
		parts := strings.Fields(auth)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", "Invalid Authorization header"
		}

		return parts[1], ""
	}

	tokenString, err := c.Cookie(AccessTokenCookieName)
	if err != nil || strings.TrimSpace(tokenString) == "" {
		return "", "Authorization header or access token cookie missing"
	}

	return tokenString, ""
}