package implementations

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PaymentProcessor struct {
	Client       *ethclient.Client
	Ledger       string
	LedgerPublic string
}

func (p *PaymentProcessor) ProcessNative(to string) (string, error) {
	fromAddress := common.HexToAddress(p.LedgerPublic)
	toAddress := common.HexToAddress(to)

	privateKey, err := crypto.HexToECDSA(p.Ledger)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	nonce, err := p.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := p.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	gasLimit := uint64(21000)

	value := new(big.Int)
	value.SetString("100000000000000000000", 10) // 100 ETH in wei (1 ETH = 10^18 wei)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     nil,
	})

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = p.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}
