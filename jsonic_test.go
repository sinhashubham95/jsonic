package jsonic_test

import (
	"io/ioutil"
	"testing"

	"github.com/sinhashubham95/jsonic"
	"github.com/stretchr/testify/assert"
)

func readFromFile(path string, t *testing.T) []byte {
	data, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	return data
}

func TestNewSuccess(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test1.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)
}

func TestNewError(t *testing.T) {
	j, err := jsonic.New(nil)
	assert.Nil(t, j)
	assert.Error(t, err)
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}

func TestChild(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test1.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)

	// simple
	c1, err := j.Child("c")
	assert.NoError(t, err)
	assert.NotNil(t, c1)
	s, err := c1.GetString(".")
	assert.NoError(t, err)
	assert.Equal(t, "d", s)

	// same
	c2, err := j.Child(".")
	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, j, c2)

	// not found
	c3, err := j.Child("a.m")
	assert.Error(t, err)
	assert.Nil(t, c3)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	// array
	c4, err := j.Child("a.arr")
	assert.NoError(t, err)
	assert.NotNil(t, c4)
	m, err := c4.GetMap("[0]")
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, "b", m["a"])

	// array nested query
	c5, err := j.Child("a.arr.[0].c.d")
	assert.NoError(t, err)
	assert.NotNil(t, c5)
	s, err = c5.GetString("e")
	assert.NoError(t, err)
	assert.Equal(t, "f", s)

	// array index not found
	c6, err := j.Child("a.arr.xyz")
	assert.Error(t, err)
	assert.Nil(t, c6)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	// array index out of bound
	c7, err := j.Child("a.arr.[1]")
	assert.Error(t, err)
	assert.Nil(t, c7)
	assert.Equal(t, jsonic.ErrNoDataFound, err)
}

func TestGet(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test1.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)

	// child exists
	v, err := j.Get("a.x")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "p", v)

	// child not exists
	v, err = j.Get("a.y")
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestTyped(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test2.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)

	i, err := j.GetInt("z")
	assert.Equal(t, 0, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i, err = j.GetInt("a")
	assert.Equal(t, 1, i)
	assert.NoError(t, err)

	i64, err := j.GetInt64("z")
	assert.Equal(t, int64(0), i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i64, err = j.GetInt64("a")
	assert.Equal(t, int64(1), i64)
	assert.NoError(t, err)
}
