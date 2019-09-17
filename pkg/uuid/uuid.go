package uuid

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

// GenUUID ...
func GenUUID() (string, error) {
	rand.Seed(time.Now().Unix())
	node, err := snowflake.NewNode(int64(rand.Intn(1023)))
	if err != nil {
		return "", err
	}

	return node.Generate().Base58(), nil
}
