package onchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
)

const gasLimit = uint64(21000000) // TODO make sure this value more dynamic

func GetTx(client *ethclient.Client, config GetTxConfig) (*types.Transaction, error) {
	chainIdInt := int(config.ChainId.Int64())
	if chainIdInt == chains.Ethereum || chainIdInt == chains.Polygon {
		return GetDynamicFeeTx(client, config.ChainId, config.FromAddress, config.To, config.Value, config.Data)
	} else {
		return GetLegacyTx(client, config.FromAddress, config.To, config.Value, config.Data)
	}
}

func GetDynamicFeeTx(client *ethclient.Client, chainID *big.Int, fromAddress common.Address, to string, value *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas tip cap: %v", err)
	}

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas fee cap: %v", err)
	}

	toAddress := common.HexToAddress(to)

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	}), nil
}

func GetLegacyTx(client *ethclient.Client, fromAddress common.Address, to string, value *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	toAddress := common.HexToAddress(to)

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}), nil
}

// ReadContractName reads the 'name' public variable from a contract.
func ReadContractName(client *ethclient.Client, contractAddress common.Address) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return "", err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["name"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// Unpack the result
	var contractName string
	err = parsedABI.UnpackIntoInterface(&contractName, "name", result)
	if err != nil {
		return "", err
	}

	return contractName, nil
}

// ReadContractSymbol reads the 'symbol' public variable from a contract.
func ReadContractSymbol(client *ethclient.Client, contractAddress common.Address) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return "", err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["symbol"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// Unpack the result
	var contractName string
	err = parsedABI.UnpackIntoInterface(&contractName, "name", result)
	if err != nil {
		return "", err
	}

	return contractName, nil
}

// ReadContractDecimals reads the 'decimals' public variable from a contract.
func ReadContractDecimals(client *ethclient.Client, contractAddress common.Address) (uint8, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return 0, err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["decimals"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	// Unpack the result
	var decimals uint8
	err = parsedABI.UnpackIntoInterface(&decimals, "decimals", result)
	if err != nil {
		return 0, err
	}

	return decimals, nil
}

// ReadContractNonce reads the 'nonces' public variable from a contract.
func ReadContractNonce(client *ethclient.Client, publicAddress common.Address, contractAddress common.Address) (int64, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return -1, err
	}

	data, err := parsedABI.Pack("nonces", publicAddress)
	if err != nil {
		return -1, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return -1, err
	}

	// Unpack the result
	var nonce *big.Int
	err = parsedABI.UnpackIntoInterface(&nonce, "nonces", result)
	if err != nil {
		return -1, err
	}

	return nonce.Int64(), nil
}

// ReadContractAllowance reads the allowance a given contract has for a wallet.
func ReadContractAllowance(client *ethclient.Client, erc20Address common.Address, publicAddress common.Address, spenderAddress common.Address) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("allowance", publicAddress, spenderAddress)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &erc20Address,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	// Unpack the result
	var allowance *big.Int
	err = parsedABI.UnpackIntoInterface(&allowance, "allowance", result)
	if err != nil {
		return nil, err
	}

	return allowance, nil
}

// TODO function params can be clearer

func ApproveTokenForRouter(client *ethclient.Client, chainId int, key string, erc20Address common.Address, publicAddress common.Address, spenderAddress common.Address) error {
	// Parse the USDC contract ABI to get the 'Approve' function signature
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi))
	if err != nil {
		return fmt.Errorf("failed to parse USDC ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("approve", spenderAddress, amounts.BigMaxUint256)
	if err != nil {
		return fmt.Errorf("failed to pack data for approve: %v", err)
	}

	chainIdBig := big.NewInt(int64(chainId))

	getTxConfig := GetTxConfig{
		ChainId:     chainIdBig,
		FromAddress: publicAddress,
		Value:       big.NewInt(0),
		To:          erc20Address.Hex(),
		Data:        data,
	}
	approvalTx, err := GetTx(client, getTxConfig) // TODO improve common.Address <-> string conversions
	if err != nil {
		return fmt.Errorf("failed to get dynamic fee tx: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return fmt.Errorf("failed to convert private key: %v", err)
	}

	// Sign the transaction
	approvalTxSigned, err := types.SignTx(approvalTx, types.LatestSignerForChainID(chainIdBig), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), approvalTxSigned)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	fmt.Println("Approval transaction sent!")

	helpers.PrintBlockExplorerTxLink(chainId, approvalTxSigned.Hash().String())
	_, err = WaitForTransaction(client, approvalTxSigned.Hash())
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	return nil
}

func WaitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	periodCount := 0
	waitingForTxText := "Waiting for transaction to be mined"
	clearLine := strings.Repeat(" ", len(waitingForTxText)+3)
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			fmt.Println() // End the animated waiting text
			fmt.Println("Transaction complete!")
			return receipt, nil
		}
		if err != nil {
			fmt.Printf("\r%s", clearLine) // Clear the current line
			fmt.Printf("\r%s%s", waitingForTxText, strings.Repeat(".", periodCount))
			periodCount = (periodCount + 1) % 4
		}
		select {
		case <-time.After(1000 * time.Millisecond): // check again after a delay
		case <-context.Background().Done():
			fmt.Println() // End the animated waiting text
			fmt.Println("Context cancelled")
			return nil, context.Background().Err()
		}
	}
}
