package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

var (
	D1DB *cloudflare.D1Database
)

type DBHealth struct {
	Status string `json:"status"`
	Type   string `json:"type"`
}

// ==========================
// تهيئة البيئة
// ==========================
func LoadEnv() {
	if os.Getenv("ENVIRONMENT") == "" {
		os.Setenv("ENVIRONMENT", "development")
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "")
	}
	if os.Getenv("DATABASE_NAME") == "" {
		os.Setenv("DATABASE_NAME", "nawthtech")
	}
}

// ==========================
// تهيئة D1 Cloudflare
// ==========================
func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("DATABASE_NAME")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("CLOUDFLARE_API_TOKEN environment variable is required")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to init Cloudflare API: %v", err)
	}

	db := api.D1(databaseName)
	D1DB = db

	log.Println("✅ D1 Cloudflare initialized successfully")
}

// ==========================
// Health Check
// ==========================
func HealthCheck() *DBHealth {
	if D1DB == nil {
		return &DBHealth{
			Status: "disconnected",
			Type:   "d1",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// تجربة استعلام بسيط
	_, err := D1DB.Query(ctx, "SELECT 1")
	if err != nil {
		return &DBHealth{
			Status: "unhealthy",
			Type:   "d1",
		}
	}

	return &DBHealth{
		Status: "healthy",
		Type:   "d1",
	}
}

// ==========================
// Users
// ==========================
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserByID(ctx context.Context, userID string) (*User, error) {
	if D1DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := "SELECT id, name, email FROM users WHERE id = ?"
	row := D1DB.QueryRow(ctx, query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}

func GetUsers(ctx context.Context, limit int) ([]User, error) {
	if D1DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := "SELECT id, name, email FROM users LIMIT ?"
	rows, err := D1DB.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// ==========================
// Services
// ==========================
type Service struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func GetServiceByID(ctx context.Context, serviceID string) (*Service, error) {
	if D1DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := "SELECT id, title, description, price FROM services WHERE id = ?"
	row := D1DB.QueryRow(ctx, query, serviceID)

	var service Service
	if err := row.Scan(&service.ID, &service.Title, &service.Description, &service.Price); err != nil {
		return nil, fmt.Errorf("service not found: %v", err)
	}

	return &service, nil
}

func GetServices(ctx context.Context, limit int) ([]Service, error) {
	if D1DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := "SELECT id, title, description, price FROM services LIMIT ?"
	rows, err := D1DB.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Price); err != nil {
			continue
		}
		services = append(services, s)
	}

	return services, nil
}