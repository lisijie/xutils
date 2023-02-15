package xutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayColumn(t *testing.T) {
	type User struct {
		Id   int
		Age  int
		Name string
	}
	var users = []User{
		{1, 20, "Jack"},
		{2, 20, "Tom"},
		{3, 22, "Bob"},
	}
	assert.Equal(t, []int{1, 2, 3}, ArrayColumn(users, func(v User) int {
		return v.Id
	}, true))
	assert.Equal(t, []int{20, 20, 22}, ArrayColumn(users, func(v User) int {
		return v.Age
	}, false))
	assert.Equal(t, []string{"Jack", "Tom", "Bob"}, ArrayColumn(users, func(v User) string {
		return v.Name
	}, false))
}

func TestCartesian(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	result := [][]int{{1, 4, 6}, {1, 4, 7}, {1, 5, 6}, {1, 5, 7}, {2, 4, 6}, {2, 4, 7}, {2, 5, 6}, {2, 5, 7}, {3, 4, 6}, {3, 4, 7}, {3, 5, 6}, {3, 5, 7}}
	assert.Equal(t, result, Cartesian([][]int{a, b, c}))
}

func TestArrayUnique(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, ArrayUnique([]int{1, 1, 1, 2, 3}))
	assert.Equal(t, []string{"a", "b", "c"}, ArrayUnique([]string{"a", "a", "b", "c"}))
}

func TestInArray(t *testing.T) {
	assert.True(t, InArray(1, []int{1, 2, 3}))
	assert.True(t, InArray(1, []int64{1, 2, 3}))
	assert.True(t, InArray("a", []string{"a", "b"}))
	assert.False(t, InArray("nil", []string{"a", "b"}))
}

func TestArrayDiff(t *testing.T) {
	assert.Equal(t, []string{"blue"}, ArrayDiff([]string{"green", "red", "blue", "red"}, []string{"green", "yellow", "red"}))
}

func TestArrayReverse(t *testing.T) {
	assert.Equal(t, []int{3, 2, 1}, ArrayReverse([]int{1, 2, 3}))
}

func TestMap(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	var users = []User{
		{1, "bob"},
		{2, "jack"},
	}
	assert.Equal(t, map[int]User{1: {1, "bob"}, 2: {2, "jack"}}, Map(users, func(v User) int {
		return v.Id
	}))
}
