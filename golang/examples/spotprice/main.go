package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/tokenprices"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the token prices request
	priceParameters := tokenprices.ChainControllerByAddressesParams{
		Currency: tokenprices.GetCurrencyParameter(tokenprices.USD),
	}

	// Execute token prices request
	tokenPrices, _, err := c.TokenPrices.GetPrices(context.Background(), priceParameters)
	if err != nil {
		log.Fatalf("Failed to get token prices: %v", err)
	}

	prettyPrintMap(*tokenPrices)
}

// Helper function to pretty print a map of token prices
func prettyPrintMap(m map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("{")
	for _, k := range keys {
		fmt.Printf("  %v: %v,\n", k, m[k])
	}
	fmt.Println("}")
}