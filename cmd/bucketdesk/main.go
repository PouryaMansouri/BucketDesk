package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/PouryaMansouri/BucketDesk/internal/profiles"
	"github.com/PouryaMansouri/BucketDesk/internal/server"
)

func main() {
	port := flag.Int("port", defaultPort(), "HTTP port for the local UI")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	store, err := profiles.NewStore("")
	if err != nil {
		logger.Error("failed to initialize profile store", "error", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	app := server.New(store, logger)
	httpServer := &http.Server{
		Handler:           app.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	url := "http://" + listener.Addr().String()
	logger.Info("BucketDesk is running", "url", url, "os", runtime.GOOS)
	fmt.Printf("BucketDesk: %s\n", url)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
}

func defaultPort() int {
	if raw := os.Getenv("BUCKETDESK_PORT"); raw != "" {
		var port int
		if _, err := fmt.Sscanf(raw, "%d", &port); err == nil && port > 0 {
			return port
		}
	}
	return 5217
}
