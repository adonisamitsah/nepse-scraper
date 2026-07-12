package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/voidarchive/go-nepse"
)

func main() {
	// Initialize default scraping configuration parameters
	opts := nepse.DefaultOptions()
	opts.TLSVerification = false

	// Spin up the core client execution engine
	client, err := nepse.NewClient(opts)
	if err != nil {
		log.Fatal("Engine initialization failure: ", err)
	}

	log.Println("Acquiring live market data matrix for all active instruments...")
	prices, err := client.LiveMarket(context.Background())
	if err != nil {
		log.Fatal("Failed to fetch live market dataset: ", err)
	}

	// Format raw payload data structures into human-readable indented JSON
	payload, err := json.MarshalIndent(prices, "", "  ")
	if err != nil {
		log.Fatal("JSON serialization parsing error: ", err)
	}

	// Force the creation of the target historical storage directory path
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to compile target data storage directory: ", err)
	}

	// Lock execution clock explicitly to Nepal Time (UTC + 5 hours 45 minutes)
	loc := time.FixedZone("NPT", 5*60*60+45*60)
	nepalTime := time.Now().In(loc)
	fileName := fmt.Sprintf("%s.json", nepalTime.Format("2006-01-02"))
	filePath := filepath.Join(dataDir, fileName)

	// Save the structured daily market snapshot payload to disk file system
	err = os.WriteFile(filePath, payload, 0644)
	if err != nil {
		log.Fatal("Failed to execute local data file commit: ", err)
	}

	fmt.Printf("Daily snapshot successfully generated and archived to: %s\n", filePath)
}