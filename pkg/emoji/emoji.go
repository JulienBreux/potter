package emoji

import (
	"crypto/rand"
	"math/big"
)

func GetRandom() string {
	return emojis[randUint8(len(emojis))]
}

func randUint8(max int) (n int) {
	if i, err := rand.Int(rand.Reader, big.NewInt(int64(max))); err == nil {
		n = int(i.Uint64())
	}

	return
}
