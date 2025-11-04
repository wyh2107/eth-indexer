package db

import "time"

type Block struct {
	Number    uint64 `gorm:"primaryKey"`
	Hash      string `gorm:"type:char(66);uniqueIndex"`
	Timestamp time.Time
	Miner     string `gorm:"type:char(66);index"`
	GasUsed   uint64
	GasLimit  uint64
	TxCount   int
}

type Transaction struct {
	Hash        string `gorm:"type:char(66);primaryKey"`
	BlockNumber uint64
	FromAddress string `gorm:"type:char(66);index"`
	ToAddress   string `gorm:"type:char(66);index"`
	Value       string
	GasPrice    string
}

type Receipt struct {
	TxHash  string `gorm:"type:char(66);primaryKey"`
	Status  bool
	GasUsed uint64
	Logs    string `gorm:"type:json"`
}
