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

	i, err = j.GetInt("e")
	assert.Equal(t, 0, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	i64, err := j.GetInt64("z")
	assert.Equal(t, int64(0), i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i64, err = j.GetInt64("a")
	assert.Equal(t, int64(1), i64)
	assert.NoError(t, err)

	i64, err = j.GetInt64("e")
	assert.Equal(t, int64(0), i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f, err := j.GetFloat("z")
	assert.Equal(t, float32(0.0), f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f, err = j.GetFloat("b")
	assert.Equal(t, float32(2.2), f)
	assert.NoError(t, err)

	f, err = j.GetFloat("e")
	assert.Equal(t, float32(0.0), f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f64, err := j.GetFloat64("z")
	assert.Equal(t, 0.0, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f64, err = j.GetFloat64("b")
	assert.Equal(t, 2.2, f64)
	assert.NoError(t, err)

	f64, err = j.GetFloat64("e")
	assert.Equal(t, 0.0, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	b, err := j.GetBool("z")
	assert.Equal(t, false, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	b, err = j.GetBool("c")
	assert.Equal(t, true, b)
	assert.NoError(t, err)

	b, err = j.GetBool("e")
	assert.Equal(t, false, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	s, err := j.GetString("z")
	assert.Equal(t, "", s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	s, err = j.GetString("d")
	assert.Equal(t, "naruto", s)
	assert.NoError(t, err)

	s, err = j.GetString("e")
	assert.Equal(t, "", s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)
}

func TestTypedArray(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test2.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)

	i, err := j.GetIntArray("z")
	assert.Nil(t, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i, err = j.GetIntArray("e")
	assert.Equal(t, []int{1, 2}, i)
	assert.NoError(t, err)

	i, err = j.GetIntArray("a")
	assert.Nil(t, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	i64, err := j.GetInt64Array("z")
	assert.Nil(t, i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i64, err = j.GetInt64Array("e")
	assert.Equal(t, []int64{1, 2}, i64)
	assert.NoError(t, err)

	i64, err = j.GetInt64Array("a")
	assert.Nil(t, i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f, err := j.GetFloatArray("z")
	assert.Nil(t, f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f, err = j.GetFloatArray("f")
	assert.Equal(t, []float32{1.1, 2.2}, f)
	assert.NoError(t, err)

	f, err = j.GetFloatArray("a")
	assert.Nil(t, f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f64, err := j.GetFloat64Array("z")
	assert.Nil(t, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f64, err = j.GetFloat64Array("f")
	assert.Equal(t, []float64{1.1, 2.2}, f64)
	assert.NoError(t, err)

	f64, err = j.GetFloat64Array("a")
	assert.Nil(t, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	b, err := j.GetBoolArray("z")
	assert.Nil(t, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	b, err = j.GetBoolArray("g")
	assert.Equal(t, []bool{true, false}, b)
	assert.NoError(t, err)

	b, err = j.GetBoolArray("a")
	assert.Nil(t, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	s, err := j.GetStringArray("z")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	s, err = j.GetStringArray("h")
	assert.Equal(t, []string{"naruto", "boruto"}, s)
	assert.NoError(t, err)

	s, err = j.GetStringArray("a")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)
}

func TestTypedMap(t *testing.T) {
	j, err := jsonic.New(readFromFile("test_data/test2.json", t))
	assert.NotNil(t, j)
	assert.NoError(t, err)

	i, err := j.GetIntMap("z")
	assert.Nil(t, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i, err = j.GetIntMap("i")
	assert.Equal(t, map[string]int{"naruto": 1}, i)
	assert.NoError(t, err)

	i, err = j.GetIntMap("a")
	assert.Nil(t, i)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	i64, err := j.GetInt64Map("z")
	assert.Nil(t, i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	i64, err = j.GetInt64Map("i")
	assert.Equal(t, map[string]int64{"naruto": 1}, i64)
	assert.NoError(t, err)

	i64, err = j.GetInt64Map("a")
	assert.Nil(t, i64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f, err := j.GetFloatMap("z")
	assert.Nil(t, f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f, err = j.GetFloatMap("j")
	assert.Equal(t, map[string]float32{"naruto": 1.1}, f)
	assert.NoError(t, err)

	f, err = j.GetFloatMap("a")
	assert.Nil(t, f)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	f64, err := j.GetFloat64Map("z")
	assert.Nil(t, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	f64, err = j.GetFloat64Map("j")
	assert.Equal(t, map[string]float64{"naruto": 1.1}, f64)
	assert.NoError(t, err)

	f64, err = j.GetFloat64Map("a")
	assert.Nil(t, f64)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	b, err := j.GetBoolMap("z")
	assert.Nil(t, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	b, err = j.GetBoolMap("k")
	assert.Equal(t, map[string]bool{"naruto": true}, b)
	assert.NoError(t, err)

	b, err = j.GetBoolMap("a")
	assert.Nil(t, b)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)

	s, err := j.GetStringMap("z")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrNoDataFound, err)

	s, err = j.GetStringMap("l")
	assert.Equal(t, map[string]string{"naruto": "rocks"}, s)
	assert.NoError(t, err)

	s, err = j.GetStringMap("a")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, jsonic.ErrInvalidType, err)
}
