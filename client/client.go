package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}


// main starts an HTTP client that sends a GET /cotacao to
// http://localhost:8080, waits for the response and saves the received
// exchange rate to a file named "cotacao.txt" in the current directory.
//
// The client timeouts if the request takes more than 300ms.
//
// It logs the received response.
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("client: context timeout exceeded")
		return
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("recebido: ")
		json.NewEncoder(log.Writer()).Encode(response)

		file, err := os.Create("cotacao.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(" Dólar: {" + response.Bid + "}\n")
		if err != nil {
			log.Fatal(err)
		}
	}

}
