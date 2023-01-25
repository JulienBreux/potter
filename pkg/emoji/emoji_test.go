package emoji_test

import (
	"testing"

	"github.com/julienbreux/potter/pkg/emoji"
	"github.com/stretchr/testify/assert"
)

func TestGetRandom(t *testing.T) {
	e0 := emoji.GetRandom()
	e1 := emoji.GetRandom()

	assert.NotEqual(t, e0, e1)
}
