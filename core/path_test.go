package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilStringConcat(t *testing.T) {
	path := buildPath(nil)

	assert.Equal(t, "", path)
}

func TestEmptyStringConcatArraySizeOne(t *testing.T) {
	path := buildPath([]string{""})

	assert.Equal(t, "", path)
}

func TestStringConcatArraySizeOne(t *testing.T) {
	path := buildPath([]string{"/a"})

	assert.Equal(t, "/a", path)
}

func TestStringConcatArraySizeTwo(t *testing.T) {
	path := buildPath([]string{"/a", "/b"})

	assert.Equal(t, "/a/b", path)
}

func TestStringConcatArraySizeThree(t *testing.T) {
	path := buildPath([]string{"/a/", "b", "/c"})

	assert.Equal(t, "/a/b/c", path)
}

func TestStringConcatArraySizeFour(t *testing.T) {
	path := buildPath([]string{"/a/", "b", "/c", "/d"})

	assert.Equal(t, "/a/b/c/d", path)
}
