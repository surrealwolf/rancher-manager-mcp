package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rancher/rancher-manager-mcp/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		rancherURL   = flag.String("rancher-url", "", "Rancher Manager API URL (required)")
		rancherToken = flag.String("rancher-token", "", "Rancher API token (required)")
	)
	flag.Parse()

	if *rancherURL == "" || *rancherToken == "" {
		fmt.Fprintf(os.Stderr, "Error: rancher-url and rancher-token are required\n")
		flag.Usage()
		os.Exit(1)
	}

	client := server.NewRancherClient(*rancherURL, *rancherToken)
	ctx := context.Background()

	logrus.Infof("Verifying token against Rancher at %s", *rancherURL)
	if err := client.VerifyToken(ctx); err != nil {
		logrus.Fatalf("Token verification failed: %v", err)
	}

	logrus.Info("Token verified successfully!")
}
