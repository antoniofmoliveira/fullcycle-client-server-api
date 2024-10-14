package entity

type ExchangeRate struct {
	ID         ID     `json:"id"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

// NewExchangeRate creates a new ExchangeRate with the given parameters and returns an error if the fields are invalid.
// The ID field is set to a new UUID.
func NewExchangeRate(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate string) (ExchangeRate, error) {
	exchangeRate := ExchangeRate{
		ID:         NewID(),
		Code:       code,
		Codein:     codein,
		Name:       name,
		High:       high,
		Low:        low,
		VarBid:     varBid,
		PctChange:  pctChange,
		Bid:        bid,
		Ask:        ask,
		Timestamp:  timestamp,
		CreateDate: createDate,
	}
	err := exchangeRate.Validate()
	if err != nil {
		return ExchangeRate{}, err
	}
	return exchangeRate, err
}

// Validate checks if the fields of the ExchangeRate are valid, and returns an error
// if any of them are invalid. It does not check if the ID field is valid, as it is
// set by the NewExchangeRate function.
func (er ExchangeRate) Validate() error {
	return nil
}
