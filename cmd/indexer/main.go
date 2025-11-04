package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"eth-indexer/internal/db"
	"eth-indexer/internal/eth"
	"eth-indexer/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	// 读取环境变量
	_ = godotenv.Load()
	rpc := os.Getenv("ETH_RPC_URL")
	if rpc == "" {
		log.Fatal("ETH_RPC_URL not set ........")
	}

	// 初始化 MySQL
	if err := db.Init(
		getenv("MYSQL_USER", "root"),
		getenv("MYSQL_PASSWORD", "123456"),
		getenv("MYSQL_HOST", "localhost"),
		getenv("MYSQL_PORT", "3306"),
		getenv("MYSQL_DB", "ethdata"),
	); err != nil {
		log.Fatal(err)
	}

	// 连接以太坊
	client, err := eth.NewClient(rpc)
	if err != nil {
		log.Fatal(err)
	}

	networkChainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("RPC ChainId failed:%v", err)
	}
	log.Printf("network chain id: %s", networkChainID)

	fmt.Println("🚀 indexer starting...")
	// 可换成从 DB/Redis 读取上次断点，这里演示从固定高度起
	start := uint64(18_900_000)
	service.StartSync(client, start, networkChainID)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
