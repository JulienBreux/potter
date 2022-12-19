package color

import (
	"crypto/rand"
	"math/big"

	"github.com/go-playground/colors"
)

// RandomColor returns a random HTML color
func RandomColor() string {
	rgb, _ := colors.RGB(
		mustGenerateRandomNumber(),
		mustGenerateRandomNumber(),
		mustGenerateRandomNumber(),
	)

	return rgb.ToHEX().String()
}

// mustGenerateRandomNumber returns a random number between 0 and 255
func mustGenerateRandomNumber() uint8 {
	const full = 255
	if i, err := rand.Int(rand.Reader, big.NewInt(full)); err != nil {
		return uint8(i.Uint64())
	}

	return 0
}
