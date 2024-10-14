package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/antoniofmoliveira/fullcycle-client-server-api/server/docs"
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/entity"
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/infra/database"
	"github.com/antoniofmoliveira/fullcycle-client-server-api/server/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// @title           Api Cotação do Dólar
// @version         1.0
// @description     Api para buscar e guardar cotação do dólar
// @termsOfService  http://swagger.io/terms/

// @contact.name   Antonio Francisco Macedo de Oliveira
// @contact.url    http://github.com/antoniofmoliveira/fullcycle-client-server-api
// @contact.email  antoniofmoliveira@gmail.com

// @license.name   MIT
// @license.url    http://github.com/antoniofmoliveira/fullcycle-client-server-api/license.txt

// @host      localhost:8080
// @BasePath  /

func main() {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database", err)
		return
	}
	db.AutoMigrate(&entity.ExchangeRate{})

	exchangeRateDB := database.NewExchangeRate(db)
	exchangeRateHandler := handlers.NewExchangeRateHandler(exchangeRateDB)

	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/cotacao", exchangeRateHandler.GetExchangeRate)

	http.HandleFunc("/docs/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/swagger.json")))

	http.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	go func() {
		fmt.Println("Server is running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	log.Println("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown the server: %v\n", err)
	}
	fmt.Println("Server stopped")
	os.Exit(0)
}
