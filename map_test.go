package xutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapFlip(t *testing.T) {
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, MapFlip(map[int]string{1: "a", 2: "b", 3: "c"}))
}
