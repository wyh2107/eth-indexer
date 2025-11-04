package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"
)

func FetchBlockData(client *ethclient.Client, blockNum *big.Int) (*types.Block, []*types.Receipt, error) {
	block, _ := client.BlockByNumber(context.Background(), blockNum)
	var receipts []*types.Receipt
	for _, tx := range block.Transactions() {
		receipt, _ := client.TransactionReceipt(context.Background(), tx.Hash())
		receipts = append(receipts, receipt)
	}
	return block, receipts, nil
}
