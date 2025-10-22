package service

import (
    "context"
    "log"
    "math/big"
    "os"
    "time"

    "eth-indexer/internal/db"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

func StartSync(client *ethclient.Client, start uint64) {
    cur := big.NewInt(int64(start))
    for {
        blk, err := client.BlockByNumber(context.Background(), cur)
        if err != nil {
            log.Printf("fetch block %s err: %v", cur.String(), err)
            time.Sleep(2 * time.Second)
            continue
        }

        b := db.Block{
            Number:    blk.NumberU64(),
            Hash:      blk.Hash().Hex(),
            Timestamp: time.Unix(int64(blk.Time()), 0),
            Miner:     blk.Coinbase().Hex(),
            GasUsed:   blk.GasUsed(),
            GasLimit:  blk.GasLimit(),
            TxCount:   len(blk.Transactions()),
        }
        db.DB.Clauses().Create(&b)

        for _, tx := range blk.Transactions() {
            msg, _ := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
            rec := db.Transaction{
                Hash:        tx.Hash().Hex(),
                BlockNumber: blk.NumberU64(),
                FromAddress: msg.From().Hex(),
                ToAddress: func() string {
                    if tx.To() == nil { return "" }
                    return tx.To().Hex()
                }(),
                Value:    tx.Value().String(),
                GasPrice: tx.GasPrice().String(),
            }
            db.DB.Clauses().Create(&rec)
        }
        log.Printf("✅ synced block %d (%d txs)", b.Number, b.TxCount)
        cur.Add(cur, big.NewInt(1))

        _ = os.Setenv("LAST_SYNCED", cur.String()) // 简易断点（可换 Redis/DB 记录）
        time.Sleep(1 * time.Second)
    }
}

