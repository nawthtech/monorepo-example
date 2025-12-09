package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/d1"
)

type D1Manager struct {
	Db d1.Database
}

var manager *D1Manager

func GetD1Manager() (*D1Manager, error) {
	if manager != nil {
		return manager, nil
	}

	dbURL := os.Getenv("D1_DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("D1_DATABASE_URL is not set")
	}

	db, err := d1.NewDatabase(dbURL)
	if err != nil {
		return nil, err
	}

	manager = &D1Manager{Db: db}
	return manager, nil
}

func (m *D1Manager) HealthCheck(ctx context.Context) (string, error) {
	_, err := m.Db.Exec(ctx, "SELECT 1")
	if err != nil {
		return "unhealthy", err
	}
	return "healthy", nil
}