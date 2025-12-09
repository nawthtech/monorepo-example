package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cloudflare/cloudflare-go" // إذا لزم الأمر، يعتمد على Cloudflare SDK
)

// D1Manager يدير الاتصال بـ D1
type D1Manager struct {
	DB       *D1Database // كائن D1
	initOnce sync.Once
}

var d1Manager *D1Manager

// InitD1 تهيئة الاتصال بـ D1 (يجب استدعاؤه مرة واحدة عند بدء السيرفر)
func InitD1() error {
	if d1Manager != nil {
		return nil
	}

	var err error
	d1Manager = &D1Manager{}

	d1Manager.initOnce.Do(func() {
		dbURL := os.Getenv("D1_DATABASE_URL")
		if dbURL == "" {
			err = fmt.Errorf("D1_DATABASE_URL is required in environment variables")
			return
		}

		// إنشاء اتصال D1
		d1, dbErr := NewD1Database(dbURL)
		if dbErr != nil {
			err = fmt.Errorf("failed to connect to D1: %v", dbErr)
			return
		}

		d1Manager.DB = d1
		log.Println("✅ Connected to D1 successfully!")
	})

	return err
}

// GetD1 إعادة كائن D1Manager
func GetD1() *D1Manager {
	if d1Manager == nil {
		log.Fatal("D1Manager not initialized. Call InitD1() first.")
	}
	return d1Manager
}

// ==== وظائف مساعدة للعمل مع D1 ====

type D1Database struct {
	URL string
	// أضف هنا أي حقول أو SDK من Cloudflare إذا لزم الأمر
}

// NewD1Database إنشاء اتصال جديد بـ D1
func NewD1Database(url string) (*D1Database, error) {
	// حالياً مجرد مثال على الاتصال
	if url == "" {
		return nil, fmt.Errorf("D1 URL is empty")
	}

	db := &D1Database{
		URL: url,
	}

	// يمكن هنا إضافة اختبار Ping إذا كانت مكتبة D1 SDK متاحة
	// مثلا: db.Ping()

	return db, nil
}

// Query تنفيذ استعلام SQL في D1
func (d *D1Database) Query(ctx context.Context, sql string, args ...any) ([][]any, error) {
	// TODO: استبدل هذا بتنفيذ حقيقي باستخدام D1 SDK أو REST API
	// حالياً مجرد مثال فارغ
	return nil, fmt.Errorf("D1 Query not implemented yet")
}

// Exec تنفيذ استعلام تعديل البيانات (INSERT, UPDATE, DELETE)
func (d *D1Database) Exec(ctx context.Context, sql string, args ...any) error {
	// TODO: استبدل هذا بتنفيذ حقيقي باستخدام D1 SDK أو REST API
	return fmt.Errorf("D1 Exec not implemented yet")
}