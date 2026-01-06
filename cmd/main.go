package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rancher/rancher-manager-mcp/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		transport          = flag.String("transport", "stdio", "Transport type: stdio or http")
		httpAddr           = flag.String("http-addr", ":8080", "HTTP server address (for http transport)")
		rancherURL         = flag.String("rancher-url", "", "Rancher Manager API URL")
		rancherToken       = flag.String("rancher-token", "", "Rancher API token")
		insecureSkipVerify = flag.Bool("insecure-skip-verify", false, "Skip SSL certificate verification (not recommended)")
		logLevel           = flag.String("log-level", "info", "Log level: debug, info, warn, error")
	)
	flag.Parse()

	// Read from environment variables if not provided via flags
	if *rancherURL == "" {
		*rancherURL = os.Getenv("RANCHER_URL")
	}
	if *rancherToken == "" {
		*rancherToken = os.Getenv("RANCHER_TOKEN")
	}
	// Check environment variable for SSL verification (flag takes precedence)
	if !*insecureSkipVerify {
		if os.Getenv("RANCHER_INSECURE_SKIP_VERIFY") == "true" || os.Getenv("RANCHER_INSECURE_SKIP_VERIFY") == "1" {
			*insecureSkipVerify = true
		}
	}

	// Set log level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)

	// Create server
	srv := server.NewServer(*rancherURL, *rancherToken, *insecureSkipVerify)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	switch *transport {
	case "stdio":
		logrus.Info("Starting MCP server with stdio transport")
		if err := srv.ServeStdio(ctx); err != nil {
			logrus.Fatalf("Failed to serve stdio: %v", err)
		}
	case "http":
		logrus.Infof("Starting MCP server with HTTP transport on %s", *httpAddr)
		httpServer := &http.Server{
			Addr:    *httpAddr,
			Handler: srv.HTTPHandler(),
		}

		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("HTTP server error: %v", err)
			}
		}()

		<-sigChan
		logrus.Info("Shutting down HTTP server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Errorf("Error shutting down server: %v", err)
		}
	default:
		log.Fatalf("Invalid transport: %s (must be 'stdio' or 'http')", *transport)
	}
}
