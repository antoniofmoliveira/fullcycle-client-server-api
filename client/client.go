package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

var errClientTimeout = errors.New("client: context timeout exceeded")

// main gets the current dollar exchange rate from the server,
// logs it and then saves it to a file named "cotacao.txt".
func main() {

	response, err := getExchangeRate()
	if err != nil {
		log.Fatal("could not get exchange rate: ", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = saveToFile(*response)
		if err != nil {
			log.Fatal("could not save to file: ", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Print("recebido: ")
		json.NewEncoder(log.Writer()).Encode(response)
	}()

	wg.Wait()
}

// getExchangeRate performs a GET request to http://localhost:8080/cotacao and
// returns the received exchange rate.
//
// It will timeout if the request takes more than 300ms.
//
// It will return the received exchange rate as a response or an error.
func getExchangeRate() (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println(errClientTimeout)
		return nil, errClientTimeout
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
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
			log.Println("could not read response: ", err)
			return nil, err
		}
		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("could not unmarshal response: ", err)
			return nil, err
		}
		return &response, nil
	}
}

// saveToFile saves the given exchange rate to a file named "cotacao.txt".
//
// The saved line is in the format " Dólar: {value}\n".
//
// If the file cannot be created or written to, it returns the error.
func saveToFile(response Response) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("could not create file: ", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(" Dólar: {" + response.Bid + "}\n")
	if err != nil {
		log.Println("could not write to file: ", err)
		return err
	}
	return nil
}
