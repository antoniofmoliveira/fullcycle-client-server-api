package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/dto"
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/entity"
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/infra/database"
)

const msgQERTimeOut = "query exchange rate: context timeout exceeded"
const msgClientTimeout = "client: context timeout exceeded"
const msgInternalError = "error while querying exchange rate"

type ExchangeRateHandler struct {
	exchangeRateDB database.ExchangeRateInterface
}

// NewExchangeRateHandler creates a new instance of ExchangeRateHandler with the given exchange rate db.
func NewExchangeRateHandler(exchangeRateDB database.ExchangeRateInterface) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		exchangeRateDB: exchangeRateDB,
	}
}

// GetExchangeRate godoc
// @Summary      Get a exchange rate for dollar
// @Description  Get a exchange rate for dollar
// @Tags         exchange rate
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ResponseOutput
// @Failure      404
// @Failure      500  {object}  dto.Message
// @Failure      504  {object}  dto.Message
// @Router       /cotacao/ [get]
func (er *ExchangeRateHandler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {

	ctxClient := r.Context()

	ctxQER, qerCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer qerCancel()

	select {
	case <-ctxQER.Done():
		go log.Println(msgQERTimeOut)
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(dto.Message{Message: msgQERTimeOut})
		return
	case <-ctxClient.Done():
		go log.Println(msgClientTimeout)
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(dto.Message{Message: msgClientTimeout})
		return
	default:
		exchangeRateInput, err := execQuery(ctxQER)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Message{Message: msgInternalError + ": " + err.Error()})
			return
		}

		go func() {
			exchangeRate, err := entity.NewExchangeRate(
				exchangeRateInput.Usdbrl.Code,
				exchangeRateInput.Usdbrl.Codein,
				exchangeRateInput.Usdbrl.Name,
				exchangeRateInput.Usdbrl.High,
				exchangeRateInput.Usdbrl.Low,
				exchangeRateInput.Usdbrl.VarBid,
				exchangeRateInput.Usdbrl.PctChange,
				exchangeRateInput.Usdbrl.Bid,
				exchangeRateInput.Usdbrl.Ask,
				exchangeRateInput.Usdbrl.Timestamp,
				exchangeRateInput.Usdbrl.CreateDate,
			)
			if err != nil {
				log.Println(err)
				return
			}
			err = er.exchangeRateDB.Create(&exchangeRate)
			if err != nil {
				log.Println(err)
			}
		}()

		response := dto.ResponseOutput{
			Bid: exchangeRateInput.Usdbrl.Bid,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		go func() {
			log.Print("received: ")
			json.NewEncoder(log.Writer()).Encode(exchangeRateInput.Usdbrl)
			log.Print("sent: ")
			json.NewEncoder(log.Writer()).Encode(response)
		}()
	}
}

// execQuery sends a GET request to the given url and returns the response body as a *dto.GetExchangeRateInput.
// The request is sent with the given context.
// If the context is canceled, the function returns immediately with the context error.
// Otherwise, it returns the response body as a *dto.GetExchangeRateInput, or an error if it fails to do so.
func execQuery(ctx context.Context) (*dto.GetExchangeRateInput, error) {
	const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println("could not read body: ", err)
			return nil, err
		}
		var c dto.GetExchangeRateInput
		err = json.Unmarshal(body, &c)
		if err != nil {
			log.Println("could not unmarshal body: ", err)
			return nil, err
		}
		return &c, nil
	}
}
