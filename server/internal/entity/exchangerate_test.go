package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewExchangeRate tests that NewExchangeRate returns a valid ExchangeRate with the given parameters,
// and that the ID field is set to a new UUID.
func TestNewExchangeRate(t *testing.T) {
	exchangeRate, err := NewExchangeRate("USD", "BRL", "Dollar", "10.00", "5.00", "0.50", "0.50", "10.00", "10.00", "2022-01-01 00:00:00", "2022-01-01 00:00:00")

	assert.Nil(t, err)
	assert.NotEmpty(t, exchangeRate.ID)
	assert.Equal(t, "USD", exchangeRate.Code)
	assert.Equal(t, "BRL", exchangeRate.Codein)
	assert.Equal(t, "Dollar", exchangeRate.Name)
	assert.Equal(t, "10.00", exchangeRate.High)
	assert.Equal(t, "5.00", exchangeRate.Low)
	assert.Equal(t, "0.50", exchangeRate.VarBid)
	assert.Equal(t, "0.50", exchangeRate.PctChange)
	assert.Equal(t, "10.00", exchangeRate.Bid)
	assert.Equal(t, "10.00", exchangeRate.Ask)
	assert.Equal(t, "2022-01-01 00:00:00", exchangeRate.Timestamp)
	assert.Equal(t, "2022-01-01 00:00:00", exchangeRate.CreateDate)
}
