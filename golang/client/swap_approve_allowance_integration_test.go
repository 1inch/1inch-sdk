package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproveAllowanceIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveControllerGetAllowanceParams
		expectedOutput swap.AllowanceResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress:  tokens.EthereumUsdc,
				WalletAddress: addresses.Vitalik,
			},
			expectedOutput: swap.AllowanceResponse{
				Allowance: "0",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			allowance, resp, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Allowance, allowance.Allowance)

			helpers.Sleep()
		})
	}
}
