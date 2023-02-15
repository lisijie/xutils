package xutils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestZipDir(t *testing.T) {
	fn := "testdata/test.zip"
	defer os.Remove(fn)
	assert.Nil(t, ZipDir("testdata/files", fn, true))
}
