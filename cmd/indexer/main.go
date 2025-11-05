package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	/*"os"
	"time"*/

	"eth-indexer/internal/config"
	"eth-indexer/internal/db"
	"eth-indexer/internal/eth"
	"eth-indexer/internal/service"
	//"github.com/joho/godotenv"
)

func main() {
	// 读取环境变量
	/*_ = godotenv.Load()
	rpc := os.Getenv("ETH_RPC_URL")
	if rpc == "" {
		log.Fatal("ETH_RPC_URL not set ........")
	}*/

	// 初始化 MySQL
	/*if err := db.Init(
		getenv("MYSQL_USER", os.Getenv("MYSQL_USER")),
		getenv("MYSQL_PASSWORD", os.Getenv("MYSQL_PASSWORD")),
		getenv("MYSQL_HOST", os.Getenv("MYSQL_HOST")),
		getenv("MYSQL_PORT", os.Getenv("MYSQL_PORT")),
		getenv("MYSQL_DB", os.Getenv("MYSQL_DB")),
	); err != nil {
		log.Fatal(err)
	}*/

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Init(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}

	// 连接以太坊
	client, err := eth.NewClient(cfg.RPCURL)
	if err != nil {
		log.Fatalf("new eth client error: %v", err)
	}

	if closer, ok := any(client).(interface {
		Close()
	}); ok {
		defer closer.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ChainIdTimeout)
	defer cancel()

	networkChainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("RPC ChainId failed:%v", err)
	}
	log.Printf("network chain id: %s", bigIntToString(networkChainID))

	fmt.Println("indexer starting...")
	// 换成从 DB/Redis 读取上次断点
	//这里演示从固定高度起
	start := uint64(18_900_000)
	service.StartSync(client, start, networkChainID)
}

/*func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}*/

func bigIntToString(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}
