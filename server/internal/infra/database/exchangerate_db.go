package database

import (
	"context"
	"log"
	"time"

	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/entity"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	DB *gorm.DB
}

// NewExchangeRate returns a new instance of ExchangeRate with the given database.
func NewExchangeRate(db *gorm.DB) *ExchangeRate {
	return &ExchangeRate{
		DB: db,
	}
}

// Create creates a new exchange rate in the database.
//
// It uses a context with a 10 milliseconds timeout to create the exchange rate.
// If the context is canceled, it logs the error and returns it.
// Otherwise, it returns the error returned by the database, or nil if it succeeds.
func (er *ExchangeRate) Create(exchangeRate *entity.ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	
	select {
	case <-ctx.Done():
		log.Println("db: ", ctx.Err())
		return ctx.Err()
	default:
		return er.DB.WithContext(ctx).Create(exchangeRate).Error
	}
}
