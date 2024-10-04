package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Usdbrl struct {
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
type ExchangeRate struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}
type Response struct {
	Bid string `json:"bid"`
}

const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

// main starts an HTTP server that listens on port 8080 and
// responds to GET /cotacao with the current dollar quotation.
func main() {
	http.HandleFunc("/cotacao", getExchangeRate)
	http.ListenAndServe(":8080", nil)
}

// getExchangeRate responds to GET /cotacao with the current dollar exchange rate.
//
// It will timeout if the request takes more than 300ms, if the query to the exchange rate API takes more than 200ms or if the database insertion takes more than 10ms.
//
// It will log the received exchange rate and the sent response.
func getExchangeRate(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database", err)
		return
	}
	db.AutoMigrate(&Usdbrl{})

	ctxClient := r.Context()

	ctxQueryExchangeRate, queryExchangeRateCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer queryExchangeRateCancel()

	ctxDb, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	select {
	case <-ctxQueryExchangeRate.Done():
		log.Println("query exchange rate: context timeout exceeded")
		return
	case <-ctxClient.Done():
		log.Println("client: context timeout exceeded")
		return
	case <-ctxDb.Done():
		log.Println("db: context timeout exceeded")
		return
	default:
		cotacao, err := execQuery(ctxQueryExchangeRate)
		if err != nil {
			log.Println("error while querying exchange rate")
			return
		}

		log.Print("received: ")
		json.NewEncoder(log.Writer()).Encode(cotacao.Usdbrl)

		resposta := Response{Bid: cotacao.Usdbrl.Bid}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resposta)

		log.Print("sent: ")
		json.NewEncoder(log.Writer()).Encode(resposta)

		err = saveExchangeRate(ctxDb, db, cotacao)
		if err != nil {
			log.Println("db: " + err.Error())
			return
		}
	}
}


// saveExchangeRate saves the given exchange rate to the database.
//
// If the context times out before the operation is finished, it returns the
// context error. Otherwise, it returns nil.
func saveExchangeRate(ctx context.Context, db *gorm.DB, c *ExchangeRate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		db.WithContext(ctx).Create(&c.Usdbrl)
		return nil
	}
}

// execQuery queries the exchange rate API and returns the received exchange rate.
//
// If the context times out before the operation is finished, it returns the
// context error. Otherwise, it returns the received exchange rate and nil.
func execQuery(ctx context.Context) (*ExchangeRate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			panic(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, error := io.ReadAll(res.Body)
		if error != nil {
			return nil, error
		}
		var c ExchangeRate
		error = json.Unmarshal(body, &c)
		if error != nil {
			return nil, error
		}
		return &c, nil
	}

}
