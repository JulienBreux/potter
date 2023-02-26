package namesgen_test

import (
	"testing"

	"github.com/JulienBreux/potter/pkg/namesgen"
	"github.com/stretchr/testify/assert"
)

func TestRandomName(t *testing.T) {
	n0 := namesgen.GetRandom()
	n1 := namesgen.GetRandom()

	assert.NotEqual(t, n0, n1)
}
