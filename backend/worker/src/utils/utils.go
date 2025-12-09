package utils

import (
	"errors"
	"os"
	"time"

	"github.com/cloudflare/d1-go/d1"
)

// ========== نوع حالة الصحة ==========
type DBHealth struct {
	Status string
	Type   string
}

// ========== دوال D1 ==========
var d1DB *d1.DB

func GetD1DB() *d1.DB {
	if d1DB == nil {
		dsn := os.Getenv("D1_DATABASE_URL")
		d1DB = d1.MustConnect(dsn)
	}
	return d1DB
}

// التحقق من صحة قاعدة البيانات
func CheckDBHealth() DBHealth {
	db := GetD1DB()
	err := db.Ping()
	if err != nil {
		return DBHealth{Status: "unhealthy", Type: "D1"}
	}
	return DBHealth{Status: "healthy", Type: "D1"}
}

// ========== دوال المستخدم ==========
func GetUserByID(db *d1.DB, id string) (map[string]interface{}, error) {
	row := db.QueryRow("SELECT id, email, name FROM users WHERE id = ?", id)
	var uid, email, name string
	err := row.Scan(&uid, &email, &name)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return map[string]interface{}{
		"id":    uid,
		"email": email,
		"name":  name,
	}, nil
}

// ========== دوال الخدمات ==========
func GetAllServices(db *d1.DB) ([]map[string]interface{}, error) {
	rows, _ := db.Query("SELECT id, title, description, price FROM services LIMIT 50")
	defer rows.Close()
	var services []map[string]interface{}
	for rows.Next() {
		var id, title, description string
		var price float64
		rows.Scan(&id, &title, &description, &price)
		services = append(services, map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"price":       price,
		})
	}
	return services, nil
}

func GetServiceByID(db *d1.DB, id string) (map[string]interface{}, error) {
	row := db.QueryRow("SELECT id, title, description, price FROM services WHERE id = ?", id)
	var sid, title, description string
	var price float64
	err := row.Scan(&sid, &title, &description, &price)
	if err != nil {
		return nil, errors.New("service not found")
	}
	return map[string]interface{}{
		"id":          sid,
		"title":       title,
		"description": description,
		"price":       price,
	}, nil
}