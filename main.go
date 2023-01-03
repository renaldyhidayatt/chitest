package main

import (
	"chigitaction/config"
	"chigitaction/router"
	"chigitaction/utils"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

func main() {
	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	}

	if _, ok := os.LookupEnv("GO_ENV"); !ok {
		err := utils.Viper("./")
		if err != nil {
			log.Fatalf(".env file not load: %v", err)
		}
	}

	r := chi.NewRouter()

	db, err := config.Database()

	if err != nil {
		log.Fatal(err.Error())
	}

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
	}))
	r.Use(middleware.Compress(gzip.BestCompression))
	r.Use(middleware.NoCache)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	router.NewTodoRoutes("/api/todo", db, r).TodoRoute()
	router.NewActivityRoutes("/api/activity-groups", db, r).ActivityRoute()

	serve := &http.Server{
		Addr:           fmt.Sprintf(":%s", viper.GetString("PORT")),
		ReadTimeout:    time.Duration(time.Second) * 60,
		WriteTimeout:   time.Duration(time.Second) * 30,
		IdleTimeout:    time.Duration(time.Second) * 120,
		MaxHeaderBytes: 3145728,
		Handler:        r,
	}

	go func() {
		err := serve.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Connected to port:", viper.GetString("PORT"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serve.Shutdown(ctx)
	os.Exit(0)

}
