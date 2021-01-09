package jsonic

import "errors"

// errors
var (
	ErrUnexpectedJSONData = errors.New("unexpected json data provided, neither array nor object")
	ErrIndexNotFound      = errors.New("expected index for json array but found something else")
	ErrIndexOutOfBound    = errors.New("index out of bounds of the json array")
)
