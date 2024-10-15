package entity

import "errors"

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

var ErrIDInvalid = errors.New("id is invalid")
var ErrCodeEmpty = errors.New("code is required")
var ErrCodeinEmpty = errors.New("codein is required")
var ErrNameEmpty = errors.New("name is required")
var ErrHighEmpty = errors.New("high is required")
var ErrLowEmpty = errors.New("low is required")
var ErrVarBidEmpty = errors.New("varBid is required")
var ErrPctChangeEmpty = errors.New("pctChange is required")
var ErrBidEmpty = errors.New("bid is required")
var ErrAskEmpty = errors.New("ask is required")
var ErrTimestampEmpty = errors.New("timestamp is required")
var ErrCreateDateEmpty = errors.New("create_date is required")

// Validate checks if the fields of the ExchangeRate are valid, and returns an error
// if any of them are invalid. It does not check if the ID field is valid, as it is
// set by the NewExchangeRate function.
func (er ExchangeRate) Validate() error {
	_, err := ParseID(er.ID.String())
	if err != nil {
		return ErrIDInvalid
	}
	if er.Code == "" {
		return ErrCodeEmpty
	}
	if er.Codein == "" {
		return ErrCodeinEmpty
	}
	if er.Name == "" {
		return ErrNameEmpty
	}
	if er.High == "" {
		return ErrHighEmpty
	}
	if er.Low == "" {
		return ErrLowEmpty
	}
	if er.VarBid == "" {
		return ErrVarBidEmpty
	}
	if er.PctChange == "" {
		return ErrPctChangeEmpty
	}
	if er.Bid == "" {
		return ErrBidEmpty
	}
	if er.Ask == "" {
		return ErrAskEmpty
	}
	if er.Timestamp == "" {
		return ErrTimestampEmpty
	}
	if er.CreateDate == "" {
		return ErrCreateDateEmpty
	}
	return nil
}
