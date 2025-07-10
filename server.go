package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	apiv1 "github.com/rainbowmga/timetravel/api/v1"
	apiv2 "github.com/rainbowmga/timetravel/api/v2"
	"github.com/rainbowmga/timetravel/service"
)

// logError logs all non-nil errors
func logError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func main() {
	router := mux.NewRouter()
	//service := service.NewInMemoryRecordService()
	service, err := service.NewSqLiteRecordService("records.db")
	logError(err)

	//api := apiv1.NewAPI(&service)

	v1API := apiv1.NewAPI(&service)
	v2API := apiv2.NewAPI(&service)

	v1Router := router.PathPrefix("/api/v1").Subrouter()
	v1Router.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	v1API.CreateRoutes(v1Router)

	v2Router := router.PathPrefix("/api/v2").Subrouter()
	v2Router.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	v2API.CreateRoutes(v2Router)

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s", address)
	log.Fatal(srv.ListenAndServe())
}
