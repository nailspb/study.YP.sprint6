package main

import (
	"context"
	"errors"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/http/handlers"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/storage/mem"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/helpers/slogHelper"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/http/middleware"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/http/routing"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	serverAddress = ":8088"
)

func main() {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				a.Value = slog.AnyValue(time.Now().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	}))
	log.Info("start TODO-LIST service")
	storage := mem.New()

	r := routing.New().
		Get("/{$}", handlers.GetAllTask(log, storage)).
		Post("/{$}", handlers.AddTask(log, storage)).
		Get("/tasks/{id}", handlers.GetTask(log, storage)).
		Delete("/tasks/{id}", handlers.DeleteTask(log, storage)).
		UseMiddleware(middleware.Logging(log)).
		UseMiddleware(middleware.Recovery(log))

	//configure http server
	srv := &http.Server{
		Addr:           serverAddress,
		Handler:        r,
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 * 1024 * 1024, //1Mb
	}

	//start http server
	log.Info("starting http server", slog.String("address", serverAddress))
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error("Failed start server", slogHelper.GetErrAttr(err))
			os.Exit(1)
		}
	}()

	//wait os signals
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt)
	signal.Notify(osSig, os.Kill)
	sig := <-osSig
	log.Info("Stop signal received", slog.String("signal", sig.String()))
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", slogHelper.GetErrAttr(err))
	}
	log.Info("Server shutdown successfully")
}
