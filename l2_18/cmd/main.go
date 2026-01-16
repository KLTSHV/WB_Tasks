package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"cmd/internal/calendar"
	"cmd/internal/httpapi"
)

func main() {
	port := flag.String("port", "", "port to listen on")
	flag.Parse()

	p := *port
	if p == "" {
		p = os.Getenv("PORT")
	}
	if p == "" {
		p = "8080"
	}

	svc := calendar.NewService()
	h := &httpapi.Handler{Svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.EventsForDay)
	mux.HandleFunc("/events_for_week", h.EventsForWeek)
	mux.HandleFunc("/events_for_month", h.EventsForMonth)
	handler := httpapi.Logging(mux)
	srv := &http.Server{
		Addr:              ":" + p,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("listening on :%s", p)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
