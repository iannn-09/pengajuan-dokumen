package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims defines the custom claims structure for JWT tokens
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// jwtSecret holds the resolved secret key (cached after first call)
var jwtSecret []byte

// getJWTSecret resolves the JWT secret from environment or generates an ephemeral one.
// Resolution order: Environment variable → Ephemeral random generation + warning log.
func getJWTSecret() []byte {
	if jwtSecret != nil {
		return jwtSecret
	}

	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecret = []byte(secret)
		return jwtSecret
	}

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Fatalf("Failed to generate ephemeral JWT secret: %v", err)
	}
	jwtSecret = []byte(hex.EncodeToString(randomBytes))
	log.Println("WARNING: JWT_SECRET not set. Generated ephemeral secret. Tokens will not persist across restarts. This is NOT suitable for production.")

	return jwtSecret
}

// GenerateToken creates a signed JWT token for a given user
func GenerateToken(userID uint, email string, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pengajuan-dokumen-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseToken validates and parses a JWT token string
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// JWTAuthMiddleware validates the JWT token from the Authorization header or query parameter 'token'
// and sets user context (user_id, email, role) for downstream handlers
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// Fallback to query parameter 'token' if Authorization header is absent (e.g. for browser target="_blank" links)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header or token query parameter is required"})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if the authenticated user has one of the allowed roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role information not found"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
		c.Abort()
	}
}

// GetUserID extracts the user ID from the Gin context (set by JWTAuthMiddleware)
func GetUserID(c *gin.Context) (uint, error) {
	id, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user_id not found in context")
	}
	userID, ok := id.(uint)
	if !ok {
		return 0, errors.New("user_id is not a valid uint")
	}
	return userID, nil
}

// GetUserRole extracts the user role from the Gin context
func GetUserRole(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("role is not a valid string")
	}
	return roleStr, nil
}

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a plaintext password against a bcrypt hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
