package color_test

import (
	"testing"

	"github.com/julienbreux/potter/pkg/color"
	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	c := color.Rand().ToHex()

	t.Logf("Color: %s", c)

	assert.Len(t, c, 7)
	assert.Equal(t, "#", c[0:1])
}

func TestParseHex(t *testing.T) {
	white := "#ffffff"
	black := "#000000"

	c, err := color.ParseHex(white)

	assert.NoError(t, err)
	assert.Equal(t, white, c.ToHex())
	assert.Equal(t, black, c.Invert().ToHex())

	blue := "#2974bf"
	orange := "#d68b40"

	c, err = color.ParseHex(blue)

	assert.NoError(t, err)
	assert.Equal(t, blue, c.ToHex())
	assert.Equal(t, orange, c.Invert().ToHex())

	shortRed := "#f00"
	red := "#ff0000"
	cyan := "#00ffff"

	c, err = color.ParseHex(shortRed)

	assert.NoError(t, err)
	assert.Equal(t, red, c.ToHex())
	assert.Equal(t, cyan, c.Invert().ToHex())
}

func TestBadParseHex(t *testing.T) {
	var err error

	white := "ffffff"

	_, err = color.ParseHex(white)
	assert.Error(t, err)

	unknown := "#ffff"
	_, err = color.ParseHex(unknown)
	assert.Error(t, err)

	unknown = "#FGH"
	_, err = color.ParseHex(unknown)
	assert.Error(t, err)
}
