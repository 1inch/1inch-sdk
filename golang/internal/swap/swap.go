package swap

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/1inch/1inch-sdk/golang/client/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/internal/onchain"
)

func ConfirmExecuteSwapWithUser(config *models.ExecuteSwapConfig) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmExecuteSwapWithUser(config, os.Stdin, stdOut)
}

func confirmExecuteSwapWithUser(config *models.ExecuteSwapConfig, reader io.Reader, writer helpers.Printer) (bool, error) {
	var permitType string
	if config.IsPermitSwap {
		permitType = "Permit1"
	} else {
		permitType = "Contract approval"
	}

	writer.Printf("Swap summary:\n")
	writer.Printf("    %-30s %s %s\n", "Selling: ", helpers.SimplifyValue(config.Amount, int(config.FromToken.Decimals)), config.FromToken.Symbol)
	writer.Printf("    %-30s %s %s\n", "Buying (estimation):", helpers.SimplifyValue(config.EstimatedAmountOut, int(config.ToToken.Decimals)), config.ToToken.Symbol)
	writer.Printf("    %-30s %v%s\n", "Slippage:", config.Slippage, "%")
	writer.Printf("    %-30s %s\n", "Permision type:", permitType)
	writer.Printf("\n")
	writer.Printf("WARNING: This swap will be executed onchain next. The results are irreversible. Make sure the proposed trade looks correct before continuing!\n")
	writer.Printf("Would you like to execute this swap onchain now? [y/N]: ")

	inputReader := bufio.NewReader(reader)
	input, _ := inputReader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "y":
		return true, nil
	default:
		return false, nil
	}
}

func ConfirmSwapDataWithUser(swapResponse *models.SwapResponse, fromAmount string, slippage float32) error {
	stdOut := helpers.StdOutPrinter{}
	return confirmSwapDataWithUser(swapResponse, fromAmount, slippage, stdOut)
}

func confirmSwapDataWithUser(swapResponse *models.SwapResponse, fromAmount string, slippage float32, writer helpers.Printer) error {
	writer.Printf("Swap summary:\n")
	writer.Printf("    %-30s %s %s\n", "Selling: ", helpers.SimplifyValue(fromAmount, int(swapResponse.FromToken.Decimals)), swapResponse.FromToken.Symbol)
	writer.Printf("    %-30s %s %s\n", "Buying (estimation):", helpers.SimplifyValue(swapResponse.ToAmount, int(swapResponse.ToToken.Decimals)), swapResponse.ToToken.Symbol)
	writer.Printf("    %-30s %v%s\n", "Slippage:", slippage, "%")
	writer.Printf("\n")
	writer.Printf("WARNING: Executing the transaction data generated by this function is irreversible. Make sure the proposed trade looks correct!\n")

	return nil
}

func ConfirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmApprovalWithUser(ethClient, publicAddress, tokenAddress, os.Stdin, stdOut)
}

func confirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string, reader io.Reader, writer helpers.Printer) (bool, error) {
	tokenName, err := onchain.ReadContractSymbol(ethClient, common.HexToAddress(tokenAddress))
	if err != nil {
		return false, fmt.Errorf("failed to read name: %v", err)
	}

	writer.Printf("The aggregator contract does not have enough allowance to execute this swap! The SDK can give an " +
		"unlimited approval on your behalf. If you would like to use custom approval amount instead, do that manually " +
		"onchain, then run the SDK again\n")
	writer.Printf("Approval summary:\n")
	writer.Printf("    %-30s %s\n", "Wallet:", publicAddress)
	writer.Printf("    %-30s %s\n", "Swapping: ", tokenName)
	writer.Printf("    %-30s %s\n", "Approval amount: ", "unlimited")
	writer.Printf("\n")
	writer.Printf("Would you like post an onchain unlimited approval now? [y/N]: ")

	inputReader := bufio.NewReader(reader)
	input, _ := inputReader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "y":
		return true, nil
	default:
		return false, nil
	}
}