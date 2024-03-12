package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/models"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/web3providers"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []client.Web3ProviderConfig{
			{
				ChainId: chains.Polygon,
				Url:     web3providers.Polygon,
			},
		},
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	quoteParams := models.GetQuoteParams{
		ChainId: chains.Polygon,
		AggregationControllerGetQuoteParams: models.AggregationControllerGetQuoteParams{
			Src:    tokens.PolygonFrax,
			Dst:    tokens.PolygonWeth,
			Amount: amounts.Ten16,
		},
	}

	swapData, _, err := c.SwapApi.GetQuote(context.Background(), quoteParams)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}

	swapDataRawIndented, err := json.MarshalIndent(swapData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal swap data: %v", err)
	}

	fmt.Printf("%s\n", string(swapDataRawIndented))
}
