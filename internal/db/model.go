package db

import "time"

type Block struct {
    Number      uint64 `gorm:"primaryKey"`
    Hash        string `gorm:"uniqueIndex"`
    Timestamp   time.Time
    Miner       string
    GasUsed     uint64
    GasLimit    uint64
    TxCount     int
}

type Transaction struct {
    Hash        string `gorm:"primaryKey"`
    BlockNumber uint64
    FromAddress string
    ToAddress   string
    Value       string
    GasPrice    string
}

type Receipt struct {
    TxHash    string `gorm:"primaryKey"`
    Status    bool
    GasUsed   uint64
    Logs      string `gorm:"type:json"`
}

