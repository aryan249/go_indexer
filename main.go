package main

import (
	"context"
	"go-indexer/abis/Vault"
	"go-indexer/eth"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rpcURL := os.Getenv("RPC_URL")
	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	vaultAddr := common.HexToAddress(contractAddr)
	vault, err := Vault.NewVault(vaultAddr, client)
	if err != nil {
		log.Fatalf("Failed to create Vault binding: %v", err)
	}

	eventCh := make(chan []eth.EventData)

	// Start event listener goroutine
	go func() {
		listener := eth.NewEventListener(ctx, client, vault, eventCh)
		listener.Listen()
		close(eventCh)
	}()

	// Start event processor goroutine
	processor := eth.NewEventProcessor(eventCh, vault, auth)
	go processor.Start()

	// Run until interrupted
	select {
	case <-ctx.Done():
	case <-time.After(24 * time.Hour): // or use a signal handler
	}
}
