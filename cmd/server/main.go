package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/egreco12/transit-backend-go/internal/config"
	"github.com/egreco12/transit-backend-go/internal/httpapi"
	"github.com/egreco12/transit-backend-go/internal/oba"
	"github.com/egreco12/transit-backend-go/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	obaClient := oba.NewClient(cfg.OBAURL, cfg.OBAApiKey)
	arrivalService := service.NewArrivalService(obaClient)
	arrivalHandler := httpapi.NewArrivalsHandler(arrivalService)
	router := httpapi.NewRouter(arrivalHandler)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
