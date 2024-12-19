package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"menu-server/src/config/env"
	"menu-server/src/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// fetchEnvVar is a helper function to fetch and validate environment variables.
func fetchEnvVar(key string) (string, error) {
	value, exists := env.AppVar[key]
	if !exists || value == "" {
		return "", errors.New("environment variable not found or empty: " + key)
	}
	return value, nil
}

// GenerateToken generates a JWT token with expiration.
func GenerateToken(user models.UserJwt, jwtType string) (string, error) {
	// Fetch the secret key for the provided JWT type.
	secretKey, err := fetchEnvVar(jwtType + "_TOKEN_SECRET")
	if err != nil {
		return "", err
	}

	// Fetch token age from configuration.
	tokenAgeStr, err := fetchEnvVar(jwtType + "_TOKEN_AGE")
	if err != nil {
		return "", err
	}

	// Parse token age duration.
	tokenAge, err := time.ParseDuration(tokenAgeStr)
	if err != nil {
		return "", errors.New("invalid token age configuration: " + err.Error())
	}

	// Current time.
	now := time.Now().UTC()

	// Create a new token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"iat":     now.Unix(),               // Issued At time.
		"exp":     now.Add(tokenAge).Unix(), // Expiration time.
	})

	// Sign the token with the secret key.
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("failed to sign token: " + err.Error())
	}

	return signedToken, nil
}

// ValidateAndExtractToken validates a JWT token and checks expiration.
func ValidateAndExtractToken(tokenString string, jwtType string) (string, string, error) {
	// Fetch the secret key for the provided JWT type.
	secretKey, err := fetchEnvVar(jwtType + "_TOKEN_SECRET")
	if err != nil {
		return "", "", err
	}

	// Parse the token and validate its signature.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Return the secret key.
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	// Check parsing errors.
	if err != nil {
		return "", "", errors.New("failed to parse token: " + err.Error())
	}

	// Validate token claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	// Check expiration explicitly.
	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", "", errors.New("token is expired")
	}

	// Extract user ID and role.
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("user_id claim is missing or invalid")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("role claim is missing or invalid")
	}

	return userID, role, nil
}

func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// bcrypt generates a salt and hashes the password with it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain password
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
