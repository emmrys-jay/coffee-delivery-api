package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func IPWhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		whitelistedIPs := strings.Split(os.Getenv("WHITELISTED_IPS"), ",")

		// Check if IP is whitelisted
		isWhitelisted := false
		for _, ip := range whitelistedIPs {
			if ip == clientIP {
				isWhitelisted = true
				break
			}
		}

		if !isWhitelisted {
			logrus.Warnf("Unauthorized IP attempt from: %s", clientIP)
			c.JSON(http.StatusForbidden, Response{
				Status:  false,
				Message: "Unknown IP address sending to webhook",
				Data:    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getTokenFromHeader(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "Authorization token not provided", Data: nil})
			c.Abort()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			logrus.Error("Error parsing token: ", err)
			c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "Invalid token", Data: nil})
			c.Abort()
			return
		}

		if role, ok := claims["role"].(string); !ok || (role != "user" && role != "admin") {
			c.JSON(http.StatusForbidden, Response{Status: false, Message: "Insufficient permissions", Data: nil})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getTokenFromHeader(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "Authorization token not provided", Data: nil})
			c.Abort()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			logrus.Error("Error parsing token: ", err)
			c.JSON(http.StatusUnauthorized, Response{Status: false, Message: "Invalid token", Data: nil})
			c.Abort()
			return
		}

		if role, ok := claims["role"].(string); !ok || role != "admin" {
			c.JSON(http.StatusForbidden, Response{Status: false, Message: "Insufficient permissions", Data: nil})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Replace with your secret key
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
