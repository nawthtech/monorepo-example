package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ==================== JSON Response ====================

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ==================== JWT Handling ====================

// GenerateJWT يولد JWT جديد باستخدام JWT_SECRET
// expDuration تحدد مدة صلاحية التوكن
func GenerateJWT(userID string, expDuration time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT يتحقق من صحة JWT ويعيد userID إذا صالح
func ValidateJWT(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
	}

	return "", errors.New("invalid token claims")
}

// ==================== Session Encryption (اختياري) ====================

// EncryptSession يشفر نص الجلسة باستخدام SESSION_SECRET
func EncryptSession(data string) (string, error) {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		return "", errors.New("SESSION_SECRET environment variable is not set")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	signature := h.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	signed := encoded + "." + base64.StdEncoding.EncodeToString(signature)
	return signed, nil
}

// DecryptSession يتحقق ويفك تشفير نص الجلسة باستخدام SESSION_SECRET
func DecryptSession(signed string) (string, error) {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		return "", errors.New("SESSION_SECRET environment variable is not set")
	}

	parts := splitOnce(signed, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid session format")
	}

	encodedData := parts[0]
	signature := parts[1]

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(decodeBase64(encodedData)))
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return "", errors.New("invalid session signature")
	}

	return decodeBase64(encodedData), nil
}

// ==================== Helpers ====================

func splitOnce(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func decodeBase64(s string) string {
	data, _ := base64.StdEncoding.DecodeString(s)
	return string(data)
}