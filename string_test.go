package xutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstr(t *testing.T) {
	s := "你好，世界"
	assert.Equal(t, "你好", Substr(s, 0, 2, ""))
	assert.Equal(t, "你好...", Substr(s, 0, 2, "..."))
	assert.Equal(t, "你好，世界", Substr(s, 0, 10, "..."))
}
