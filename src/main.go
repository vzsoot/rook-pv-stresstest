package main

import (
	"github.com/elliotxx/healthcheck"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

const (
	VolumePath1 = "VOLUME_PATH1"
	VolumePath2 = "VOLUME_PATH2"
	Role        = "ROLE"
)

func server() {
	r := gin.Default()

	r.GET("livez", healthcheck.NewHandler(healthcheck.NewDefaultHandlerConfig()))
	r.GET("readyz", healthcheck.NewHandler(healthcheck.NewDefaultHandlerConfigFor()))

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		slog.Error("Failed to start server", "err", err)
	}
	slog.Info("Server started")
}

func producer() {
	volumePath1 := os.Getenv(VolumePath1)
	volumePath2 := os.Getenv(VolumePath2)
	slog.Info("Volume path", VolumePath1, volumePath1)
	slog.Info("Volume path", VolumePath2, volumePath2)

	entries, err := os.ReadDir(volumePath1)
	if err != nil {
		slog.Error("Failed to read volume", "err", err)
		os.Exit(2)
	}

	slog.Info("Volume1 content:")
	for _, entry := range entries {
		slog.Info(entry.Name())
	}

	entries, err = os.ReadDir(volumePath2)
	if err != nil {
		slog.Error("Failed to read volume", "err", err)
		os.Exit(3)
	}

	slog.Info("Volume2 content:")
	for _, entry := range entries {
		slog.Info(entry.Name())
	}

	server()
}

func consumer() {
	server()
}

func main() {
	applicationRole := os.Getenv(Role)
	slog.Info("Application role", Role, applicationRole)

	if applicationRole == "producer" {
		producer()
	} else if applicationRole == "consumer" {
		consumer()
	} else {
		slog.Error("Application role not recognized", "applicationRole", applicationRole)
		os.Exit(1)
	}
}
