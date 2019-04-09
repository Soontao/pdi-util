package pdiutil

import (
	"strconv"
	"testing"

	"gotest.tools/assert"
)

func TestParseInt(t *testing.T) {

	v, _ := strconv.ParseInt("0012", 10, 64)
	assert.Equal(t, v, int64(12))

}
