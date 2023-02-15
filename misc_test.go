package xutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIF(t *testing.T) {
	i1, i2 := 1, 2
	assert.Equal(t, i1, IF(true, i1, i2))
	assert.Equal(t, i2, IF(false, i1, i2))

	s1, s2 := "abc", "xyz"
	assert.Equal(t, s1, IF(true, s1, s2))
	assert.Equal(t, s2, IF(false, s1, s2))

	type user struct {
		Id int
	}
	user1, user2 := &user{Id: 1}, &user{Id: 2}
	assert.Equal(t, s1, IF(user1.Id < user2.Id, s1, s2))
	assert.Equal(t, s2, IF(user1.Id > user2.Id, s1, s2))
}
