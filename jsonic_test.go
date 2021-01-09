package jsonic_test

import (
	"io/ioutil"
	"testing"

	"github.com/sinhashubham95/jsonic"
	"github.com/stretchr/testify/assert"
)

func readFromFile(path string, t *testing.T) []byte {
	data, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	return data
}

func TestNewSuccess(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test1.json", t))
	assert.NotNil(t, j)
	assert.Nil(t, err)
}

func TestNewError(t *testing.T) {
	j, err := jsonic.New(nil)
	assert.Nil(t, j)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestChild(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test1.json", t))
	assert.NotNil(t, j)
	assert.Nil(t, err)

	// simple
	c1, err := j.Child("c")
	assert.Nil(t, err)
	assert.NotNil(t, c1)
	s, err := c1.GetString(".")
	assert.Nil(t, err)
	assert.Equal(t, s, "d")

	// same
	c2, err := j.Child(".")
	assert.Nil(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, c2, j)
}
