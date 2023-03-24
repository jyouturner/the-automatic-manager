package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("started webhookserver")

	defer func() {
		if r := recover(); r != nil {
			log.Error("error", r)
			panic(r)
		}
	}()

	h := NewHandler(&Config{})

	r := NewRouter(h)

	url := fmt.Sprintf("%s:%s", "localhost", "8080")
	log.Infof("starting-server %s", url)
	runServer(url, r)
}

func runServer(url string, r *mux.Router) {
	srv := &http.Server{
		Addr:         url,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("error", r)
				panic(r)
			}
		}()
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("error", r)
			}
			log.Info("srv.ListenAndServe-shutdown")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	shutdownTimeout := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Error("srv.Shutdown", err)
	}
	log.Info("server-is-shutdown")
	cancel()
	os.Exit(0)
}

func NewRouter(handler *Handler) *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/health_check", healthCheck)
	route.HandleFunc("/echo", echo)
	return route
}

// HealthCheck is the ELB health check
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func echo(w http.ResponseWriter, r *http.Request) {
	//print out stuff
	for k, v := range r.Header {
		fmt.Printf("%v: %v\n", k, v)
	}
	w.WriteHeader(http.StatusOK)
}

type Handler struct {
}

type Config struct {
}

// New returns a new Handler
func NewHandler(config *Config) *Handler {

	return &Handler{}
}
