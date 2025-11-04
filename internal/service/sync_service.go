package service

import (
	"eth-indexer/internal/db"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"
)

// internal/service/sync_service.go
func StartSync(client *ethclient.Client, start uint64, networkChainID *big.Int) {
	current := big.NewInt(int64(start))

	for {
		block, err := client.BlockByNumber(context.Background(), current)
		if err != nil {
			log.Println("Fetch block error:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		// ... 先存 block 省略

		for _, tx := range block.Transactions() {
			var signer types.Signer
			if tx.Protected() {
				// 受保护交易：使用网络的 chainID（1/…），避免 tx.ChainId() 为 nil
				signer = types.LatestSignerForChainID(networkChainID)
			} else {
				// 老交易：使用 HomesteadSigner
				signer = (types.HomesteadSigner{})
			}

			from, err := types.Sender(signer, tx)
			if err != nil {
				log.Printf("sender decode error: %v", err)
				continue
			}

			to := ""
			if tx.To() != nil {
				to = tx.To().Hex()
			}

			dbTx := db.Transaction{
				Hash:        tx.Hash().Hex(),
				BlockNumber: block.NumberU64(),
				FromAddress: from.Hex(),
				ToAddress:   to,
				Value:       tx.Value().String(),
				GasPrice:    tx.GasPrice().String(),
			}
			db.DB.Create(&dbTx)
		}

		log.Printf("✅ Block %d synced (%d txs)", block.NumberU64(), len(block.Transactions()))
		current.Add(current, big.NewInt(1))
	}
}

//func StartSync(client *ethclient.Client, start uint64, networkChainID *big.Int) {
//	current := big.NewInt(int64(start))
//
//	for {
//		block, err := client.BlockByNumber(context.Background(), current)
//		if err != nil {
//			log.Println("Fetch block error:", err)
//			time.Sleep(time.Second * 3)
//			continue
//		}
//
//		dbBlock := db.Block{
//			Number:    block.NumberU64(),
//			Hash:      block.Hash().Hex(),
//			Timestamp: time.Unix(int64(block.Time()), 0),
//			GasUsed:   block.GasUsed(),
//			TxCount:   len(block.Transactions()),
//		}
//		db.DB.Create(&dbBlock)
//
//		for _, tx := range block.Transactions() {
//			chainId := tx.ChainId()
//			if chainId == nil {
//				log.Println("Transaction has no chain ID")
//				continue
//			}
//			signer := types.LatestSignerForChainID(networkChainID)
//			from, err := types.Sender(signer, tx)
//			if err != nil {
//				log.Printf("Failed to get sender from transaction: %v", err)
//				continue
//			}
//
//			rec := db.Transaction{
//				Hash:        tx.Hash().Hex(),
//				BlockNumber: block.NumberU64(),
//				FromAddress: from.Hex(),
//				ToAddress: func() string {
//					if tx.To() == nil {
//						return ""
//					}
//					return tx.To().Hex()
//				}(),
//				Value:    tx.Value().String(),
//				GasPrice: tx.GasPrice().String(),
//			}
//			db.DB.Create(&rec)
//		}
//
//		log.Printf("✅ Block %d synced (%d txs)", block.NumberU64(), len(block.Transactions()))
//		current.Add(current, big.NewInt(1))
//	}
//
//
//}
