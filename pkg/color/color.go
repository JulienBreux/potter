package color

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const (
	max = 255
)

var errInvalidFormat = errors.New("invalid format, hexa color must be #FFFFFF with a # as prefix")

// RGBA represents a traditional 32-bit alpha-premultiplied color, having 8 bits for each of red, green, blue and alpha.
type RGBA struct {
	R, G, B, A uint8
}

// Rand returns a random RGB color
func Rand() RGBA {
	return RGBA{
		R: randUint8(max),
		G: randUint8(max),
		B: randUint8(max),
		A: max,
	}
}

// ParseHex parses an hex color to RGB color
func ParseHex(color string) (RGBA, error) {
	const (
		four      = 4
		seven     = 7
		ten       = 10
		seventeen = 17
	)
	var (
		err error
		c   RGBA
	)

	c.A = 0xff

	if color[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + ten
		case b >= 'A' && b <= 'F':
			return b - 'A' + ten
		}
		err = errInvalidFormat
		return 0
	}

	switch len(color) {
	case seven:
		c.R = hexToByte(color[1])<<four + hexToByte(color[2])
		c.G = hexToByte(color[3])<<four + hexToByte(color[4])
		c.B = hexToByte(color[5])<<four + hexToByte(color[6])
	case four:
		c.R = hexToByte(color[1]) * seventeen
		c.G = hexToByte(color[2]) * seventeen
		c.B = hexToByte(color[3]) * seventeen
	default:
		err = errInvalidFormat
	}

	return c, err
}

// ToHex exports the color to hexadecimal
func (c RGBA) ToHex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// Invert returns the inverted RGB color
func (c RGBA) Invert() RGBA {
	return RGBA{
		R: max - c.R,
		G: max - c.G,
		B: max - c.B,
	}
}

func randUint8(max int64) (n uint8) {
	if i, err := rand.Int(rand.Reader, big.NewInt(max)); err == nil {
		n = uint8(i.Uint64())
	}

	return
}
