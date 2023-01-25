package namesgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GetRandom() string {
	fn := firstname[randUint8(len(firstname))]
	ln := lastname[randUint8(len(lastname))]

	return fmt.Sprintf("%s-%s", fn, ln)
}

func randUint8(max int) (n int) {
	if i, err := rand.Int(rand.Reader, big.NewInt(int64(max))); err == nil {
		n = int(i.Uint64())
	}

	return
}
