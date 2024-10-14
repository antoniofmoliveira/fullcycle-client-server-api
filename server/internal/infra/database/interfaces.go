package database

import (
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/entity"
)

type ExchangeRateInterface interface {
	Create(exchangeRate *entity.ExchangeRate) error
}
