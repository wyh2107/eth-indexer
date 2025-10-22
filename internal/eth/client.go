package eth
import "github.com/ethereum/go-ethereum/ethclient"

func NewClient(rpc string) (*ethclient.Client, error) {
    return ethclient.Dial(rpc)
}

