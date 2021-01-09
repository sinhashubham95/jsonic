package jsonic

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

// Jsonic is the type to hold the JSON data
type Jsonic struct {
	data  interface{}
	mu    *sync.RWMutex
	cache map[string]*Jsonic
}

type pathElement struct {
	nature int
	key    string
	index  int
}

// New is used to crete a new parser for the JSON data
func New(data []byte) (*Jsonic, error) {
	var unmarshalled interface{}
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		// not a valid json
		return nil, err
	}
	return new(unmarshalled), nil
}

// Child returns the json tree at the path specified.
//
// It returns an error in case there is nothing that can be resolved
// at the specified path.
//
// Path should be like this for example - a.[0].b, [0].a.b, etc.
// The path elements should be separated with dots.
// Now the path elements can either be the index in case of an array
// with the index enclosed within square brackets or it can be
// the key of the object.
func (j *Jsonic) Child(path string) (*Jsonic, error) {
	return j.child(strings.Split(path, dot))
}

// Get is used to get the data at the path specified.
func (j *Jsonic) Get(path string) (interface{}, error) {
	child, err := j.Child(path)
	if err != nil {
		return nil, err
	}
	return child.data, nil
}

// GetInt is used to get the integer at the path specified.
func (j *Jsonic) GetInt(path string) (int, error) {
	val, err := j.Get(path)
	if err != nil {
		return 0, err
	}
	if i, ok := val.(int); ok {
		return i, nil
	}
	return 0, ErrInvalidType
}

// GetInt64 is used to get the integer at the path specified.
func (j *Jsonic) GetInt64(path string) (int64, error) {
	val, err := j.Get(path)
	if err != nil {
		return 0, err
	}
	if i, ok := val.(int64); ok {
		return i, nil
	}
	return 0, ErrInvalidType
}

// GetFloat is used to get the floating point number at the path specified.
func (j *Jsonic) GetFloat(path string) (float32, error) {
	val, err := j.Get(path)
	if err != nil {
		return 0, err
	}
	if f, ok := val.(float32); ok {
		return f, nil
	}
	return 0, ErrInvalidType
}

// GetFloat64 is used to get the floating point number at the path specified.
func (j *Jsonic) GetFloat64(path string) (float64, error) {
	val, err := j.Get(path)
	if err != nil {
		return 0, err
	}
	if f, ok := val.(float64); ok {
		return f, nil
	}
	return 0, ErrInvalidType
}

// GetBool is used to get the integer at the path specified.
func (j *Jsonic) GetBool(path string) (bool, error) {
	val, err := j.Get(path)
	if err != nil {
		return false, err
	}
	if b, ok := val.(bool); ok {
		return b, nil
	}
	return false, ErrInvalidType
}

// GetString is used to get the string at the path specified.
func (j *Jsonic) GetString(path string) (string, error) {
	val, err := j.Get(path)
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", ErrInvalidType
}

// GetArray is used to get the data array at the path specified.
func (j *Jsonic) GetArray(path string) ([]interface{}, error) {
	val, err := j.Get(path)
	if err != nil {
		return nil, err
	}
	if a, ok := val.([]interface{}); ok {
		return a, nil
	}
	return nil, ErrInvalidType
}

// GetIntArray is used to get the integer array at the path specified.
func (j *Jsonic) GetIntArray(path string) ([]int, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	iArr := make([]int, len(val))
	for index, v := range val {
		if i, ok := v.(int); ok {
			iArr[index] = i
		}
	}
	return iArr, nil
}

// GetInt64Array is used to get the 64-bit integer array at the path specified.
func (j *Jsonic) GetInt64Array(path string) ([]int64, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	iArr := make([]int64, len(val))
	for index, v := range val {
		if i, ok := v.(int64); ok {
			iArr[index] = i
		}
	}
	return iArr, nil
}

// GetFloatArray is used to get the floating point number array at the path specified.
func (j *Jsonic) GetFloatArray(path string) ([]float32, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	fArr := make([]float32, len(val))
	for index, v := range val {
		if f, ok := v.(float32); ok {
			fArr[index] = f
		}
	}
	return fArr, nil
}

// GetFloat64Array is used to get the 64-bit floating point number array at the path specified.
func (j *Jsonic) GetFloat64Array(path string) ([]float64, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	fArr := make([]float64, len(val))
	for index, v := range val {
		if f, ok := v.(float64); ok {
			fArr[index] = f
		}
	}
	return fArr, nil
}

// GetBoolArray is used to get the boolean array at the path specified.
func (j *Jsonic) GetBoolArray(path string) ([]bool, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	bArr := make([]bool, len(val))
	for index, v := range val {
		if b, ok := v.(bool); ok {
			bArr[index] = b
		}
	}
	return bArr, nil
}

// GetStringArray is used to get the string array at the path specified.
func (j *Jsonic) GetStringArray(path string) ([]string, error) {
	val, err := j.GetArray(path)
	if err != nil {
		return nil, err
	}
	sArr := make([]string, len(val))
	for index, v := range val {
		if s, ok := v.(string); ok {
			sArr[index] = s
		}
	}
	return sArr, nil
}

// GetMap is used to get the data map at the path specified.
func (j *Jsonic) GetMap(path string) (map[string]interface{}, error) {
	val, err := j.Get(path)
	if err != nil {
		return nil, err
	}
	if m, ok := val.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, ErrInvalidType
}

// GetIntMap is used to get the integer map at the path specified.
func (j *Jsonic) GetIntMap(path string) (map[string]int, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	iMap := make(map[string]int)
	for k, v := range val {
		if i, ok := v.(int); ok {
			iMap[k] = i
		}
	}
	return iMap, nil
}

// GetInt64Map is used to get the 64-bit integer map at the path specified.
func (j *Jsonic) GetInt64Map(path string) (map[string]int64, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	iMap := make(map[string]int64)
	for k, v := range val {
		if i, ok := v.(int64); ok {
			iMap[k] = i
		}
	}
	return iMap, nil
}

// GetFloatMap is used to get the floating point number map at the path specified.
func (j *Jsonic) GetFloatMap(path string) (map[string]float32, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	fMap := make(map[string]float32)
	for k, v := range val {
		if f, ok := v.(float32); ok {
			fMap[k] = f
		}
	}
	return fMap, nil
}

// GetFloat64Map is used to get the 64-bit floating point number map at the path specified.
func (j *Jsonic) GetFloat64Map(path string) (map[string]float64, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	fMap := make(map[string]float64)
	for k, v := range val {
		if f, ok := v.(float64); ok {
			fMap[k] = f
		}
	}
	return fMap, nil
}

// GetBoolMap is used to get the boolean map at the path specified.
func (j *Jsonic) GetBoolMap(path string) (map[string]bool, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	bMap := make(map[string]bool)
	for k, v := range val {
		if b, ok := v.(bool); ok {
			bMap[k] = b
		}
	}
	return bMap, nil
}

// GetStringMap is used to get the string map at the path specified.
func (j *Jsonic) GetStringMap(path string) (map[string]string, error) {
	val, err := j.GetMap(path)
	if err != nil {
		return nil, err
	}
	sMap := make(map[string]string)
	for k, v := range val {
		if s, ok := v.(string); ok {
			sMap[k] = s
		}
	}
	return sMap, nil
}

func new(data interface{}) *Jsonic {
	return &Jsonic{
		data:  data,
		mu:    &sync.RWMutex{},
		cache: make(map[string]*Jsonic),
	}
}

func (j *Jsonic) childFromArray(array []interface{}, path []string) (*Jsonic, error) {
	// get the index, which should be there as the first path element
	index, err := getIndex(path[0])
	if err != nil {
		// some problem with the index
		return nil, ErrIndexNotFound
	}
	if index < 0 || index >= len(array) {
		// index out of bound
		return nil, ErrIndexOutOfBound
	}
	// check in cache for the child
	if cached := j.checkInCache(strconv.Itoa(index)); cached != nil {
		return cached.child(path[1:])
	}
	// create a child, and save it
	child := new(array[index])
	j.saveInCache(strconv.Itoa(index), child)
	return child.child(path[1:])
}

func (j *Jsonic) childFromObject(object map[string]interface{}, path []string) (*Jsonic, error) {
	current := path[0]
	// this loop is to handle the following scenario
	// say the path elements are as follows a, b and c
	// it might so happen that each of a, a.b and a.b.c are
	// present in the json data as keys, so we should give
	// each of them a fair chance. the only thing is we
	// are giving preference in the following order a > a.b > a.b.c
	for i, p := range path {
		if cached := j.checkInCache(current); cached != nil {
			result, err := cached.child(path[i+1:])
			if err == nil {
				// result found successfully
				return result, nil
			}
		} else if data, ok := object[current]; ok {
			child := new(data)
			j.saveInCache(current, child)
			result, err := child.child(path[i+1:])
			if err == nil {
				// result found successfully
				return result, nil
			}
		}
		// nothing here, check further
		current = joinPath(current, p)
	}
	return nil, ErrNoDataFound
}

func (j *Jsonic) child(path []string) (*Jsonic, error) {
	// first the base condition
	if len(path) == 0 {
		// we have reached the result
		return j, nil
	}
	// either data is array or object
	// we need to check that
	// and accordingly proceed
	if array, ok := j.data.([]interface{}); ok {
		return j.childFromArray(array, path)
	}
	if object, ok := j.data.(map[string]interface{}); ok {
		return j.childFromObject(object, path)
	}
	return nil, ErrUnexpectedJSONData

}

func getIndex(element string) (int, error) {
	// it should be enclosed within curly braces
	return strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(element, closeBracket), openBracket))
}

func (j *Jsonic) checkInCache(path string) *Jsonic {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.cache[path]
}

func (j *Jsonic) saveInCache(path string, child *Jsonic) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.cache[path] = child
}

func joinPath(first, second string) string {
	return first + dot + second
}
