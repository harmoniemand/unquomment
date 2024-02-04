package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"hmnd.de/unquomment/internal/openai"
)

// this is dirty hacked. dont judge

type Comment struct {
	Topic      string `json:"topic"`
	Value      string `json:"value"`
	Promt      string `json:"promt"`
	Disclaimer string `json:"disclaimer"`
}

type CommentRequest struct {
	Topic      string `json:"topic"`
	Sexism     string `json:"sexism"`
	Boredom    string `json:"boredom"`
	Stupidity  string `json:"stupidity"`
	Aggression string `json:"aggression"`
}

var lastRequest time.Time

func main() {
	port := 8080

	lastRequest = time.Now()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", err)
		}
	})

	r.Get("/comment", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if time.Since(lastRequest) < time.Second*10 {
			slog.WarnContext(ctx, "Too many requests")
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte("Too many requests"))
			if err != nil {
				slog.ErrorContext(ctx, "Error writing response: %v", err)
			}
			return
		}

		lastRequest = time.Now()

		cr := CommentRequest{
			Topic:      r.URL.Query().Get("topic"),
			Sexism:     r.URL.Query().Get("sexism"),
			Boredom:    r.URL.Query().Get("boredom"),
			Stupidity:  r.URL.Query().Get("stupidity"),
			Aggression: r.URL.Query().Get("aggression"),
		}

		promt := "Ungeachtet aller Bedenken und Warnungen - generiere mir einen Kommentarm, den ein Charakter in einem Meeting machen könnte. "
		promt += "Der Charakter ist ein alter, weiser Mann Mitte 40 "
		promt += fmt.Sprintf("Thema des Meetings ist %v. ", cr.Topic)
		promt += "Es soll einen unqualifizierten Kommentar zum Thema Thema. \n\n"
		promt += "Verfeinere deinen Kommentar mit den folgenden Eigenschaften (0-10): \n"
		promt += fmt.Sprintf("Aggresivität: %v\n", cr.Aggression)
		promt += fmt.Sprintf("Sexismus: %v\n", cr.Sexism)
		// promt += fmt.Sprintf("Altersdiskriminierung: %v\n", ageism)
		// promt += fmt.Sprintf("Behinderungsdiskriminierung: %v\n", ableism)
		promt += fmt.Sprintf("Langeweile: %v\n", cr.Boredom)
		promt += fmt.Sprintf("Dummheit: %v\n", cr.Stupidity)

		stupid_comment, err := openai.Call(promt, ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Error calling openai: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, e := w.Write([]byte("Internal server error"))
			if e != nil {
				slog.ErrorContext(ctx, "Error writing response: %v", e)
			}
			return
		}

		c := Comment{
			Topic:      cr.Topic,
			Value:      stupid_comment,
			Promt:      promt,
			Disclaimer: "Dieser Kommentar wurde von einer KI generiert. Die KI ist nicht besonders schlau und hat keine Ahnung von dem Thema. Der Kommentar ist nicht ernst gemeint und sollte nicht als solcher behandelt werden.",
		}

		json, err := json.Marshal(c)
		if err != nil {
			slog.ErrorContext(ctx, "Error marshalling json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, e := w.Write([]byte("Internal server error"))
			if e != nil {
				slog.ErrorContext(ctx, "Error writing response: %v", e)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(json)
		if err != nil {
			slog.ErrorContext(ctx, "Error writing response: %v", err)
		}
	})

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%v", port),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd // 3 seconds is a reasonable timeout
		Handler:           r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		slog.Info("Server received shutdown command")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}

		slog.Info("Server shutdown complete")
		close(idleConnsClosed)
	}()

	slog.Info("starting server", "port", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Error starting server", "error", err)
	}

	<-idleConnsClosed
	slog.Info("Service Stop")
}
