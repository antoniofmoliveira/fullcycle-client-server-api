package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

type Message struct {
	Message string `json:"message"`
}

const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const msgQERTimeOut = "query exchange rate: context timeout exceeded"
const msgClientTimeout = "client: context timeout exceeded"
const msgInternalError = "error while querying exchange rate"

var DBS gorm.DB


// main initializes the database and starts the server. It listens for SIGINT,
// SIGTERM and SIGHUP signals and shuts down the server when it receives one.
// It logs "server: shutting down" before exiting.
func main() {
	initializeDb()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-termChan
		log.Println("server: shutting down")
		os.Exit(0)
	}()

	http.HandleFunc("/cotacao", getExchangeRate)
	http.ListenAndServe(":8080", nil)

}

// initializeDb opens a connection to a SQLite database named "cotacao.db" and
// uses it to initialize the global dbs variable. If the connection cannot be
// established, it logs the error and does not set the dbs variable. It also
// migrates the Usdbrl table if it does not exist yet.
func initializeDb() {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database", err)
		return
	}
	db.AutoMigrate(&Usdbrl{})
	DBS = *db
}

// getExchangeRate responds to GET requests at /cotacao.
//
// It will timeout if the query to obtain the exchange rate takes more than 200ms.
//
// It will return the received exchange rate in the format
// {"bid": "{value}"} or an error.
//
// If the context is canceled or the query takes too long, it will return a
// StatusGatewayTimeout response with the error message.
//
// If there is an internal error, it will return a StatusInternalServerError
// response with the error message.
func getExchangeRate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	ctxClient := r.Context()

	ctxQER, qerCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer qerCancel()

	select {
	case <-ctxQER.Done():

		go log.Println(msgQERTimeOut)
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(Message{Message: msgQERTimeOut})
		return
	case <-ctxClient.Done():
		go log.Println(msgClientTimeout)
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(Message{Message: msgClientTimeout})
		return
	default:
		cotacao, err := execQuery(ctxQER)
		if err != nil {
			go log.Println(msgInternalError, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Message{Message: msgInternalError})
			return
		}

		resposta := Response{Bid: cotacao.Usdbrl.Bid}

		go func() {
			log.Print("received: ")
			json.NewEncoder(log.Writer()).Encode(cotacao.Usdbrl)
			log.Print("sent: ")
			json.NewEncoder(log.Writer()).Encode(resposta)
		}()

		go func() {
			err = saveExchangeRate(cotacao)
			if err != nil {
				log.Println("db: ", err.Error())
			}
		}()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resposta)

	}
}

// saveExchangeRate saves the given exchange rate to the database.
//
// It will timeout if the save operation takes more than 10ms.
//
// It will return the error if the save operation fails or the error if the context is canceled.
func saveExchangeRate(c *ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("db: ", ctx.Err())
		return ctx.Err()
	default:
		result := DBS.WithContext(ctx).Create(&c.Usdbrl)
		if result.Error != nil {
			log.Println("db: ", result.Error)
			return result.Error
		}
		return nil
	}
}

// execQuery performs a GET request to http://economia.awesomeapi.com.br/json/last/USD-BRL and
// returns the received exchange rate.
//
// It will timeout if the request takes more than 200ms.
//
// It will return the received exchange rate as a ExchangeRate or an error.
func execQuery(ctx context.Context) (*ExchangeRate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Println("could not create request: ", err)
			return nil, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("could not send request: ", err)
			return nil, err
		}
		defer res.Body.Close()

		body, error := io.ReadAll(res.Body)
		if error != nil {
			log.Println("could not read body: ", error)
			return nil, error
		}
		var c ExchangeRate
		error = json.Unmarshal(body, &c)
		if error != nil {
			log.Println("could not unmarshal body: ", error)
			return nil, error
		}
		return &c, nil
	}

}
