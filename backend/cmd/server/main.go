package main

import (
	"log/slog"
	"net/http"
	"os"

	httpadapter "github.com/gallardo/temp-converter/internal/adapter/http"
	"github.com/gallardo/temp-converter/internal/adapter/memory"
	"github.com/gallardo/temp-converter/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo := memory.NewHistoryRepository()
	convertOnly := usecase.NewConvertOnly()
	convertTemp := usecase.NewConvertTemperature(repo)
	getHistory := usecase.NewGetHistory(repo)
	handler := httpadapter.NewHandler(convertTemp, getHistory, convertOnly)

	mux := http.NewServeMux()
	handler.Register(mux)

	port := "8080"
	slog.Info("servidor iniciado", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("error al iniciar servidor", "err", err)
		os.Exit(1)
	}
}
