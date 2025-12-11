package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// CustomClaims هي الهيكل المخصص لـ JWT Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates access and refresh tokens
func GenerateJWT(cfg *config.Config, userID, email, role string) (string, string, error) {
	if cfg == nil {
		return "", "", errors.New("config required")
	}
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return "", "", errors.New("JWT_SECRET not set")
	}

	now := time.Now()
	accessClaims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(cfg.Auth.JWTExpiration))),
			Issuer:    "nawthtech",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(cfg.Auth.RefreshExpiration))),
		Issuer:    "nawthtech",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenStr, nil
}

// ValidateJWT validates a JWT token
func ValidateJWT(cfg *config.Config, tokenString string) (*CustomClaims, error) {
	if cfg == nil {
		return nil, errors.New("config required")
	}
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	return nil, errors.New("invalid token")
}

// RefreshJWT refreshes an access token using a refresh token
func RefreshJWT(cfg *config.Config, refreshToken string) (string, string, error) {
	if cfg == nil {
		return "", "", errors.New("config required")
	}
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return "", "", errors.New("JWT_SECRET not set")
	}

	// Validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if !token.Valid {
		return "", "", errors.New("expired refresh token")
	}

	// Extract user info from refresh token (you might want to store this differently)
	// For now, we'll generate new tokens with default values
	return GenerateJWT(cfg, "user_id", "email@example.com", "user")
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value("user_role").(string)
	if !ok {
		return "", errors.New("user role not found in context")
	}
	return role, nil
}

// AddUserToContext adds user info to context
func AddUserToContext(ctx context.Context, userID, email, role string) context.Context {
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "user_email", email)
	ctx = context.WithValue(ctx, "user_role", role)
	return ctx
}