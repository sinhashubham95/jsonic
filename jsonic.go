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
