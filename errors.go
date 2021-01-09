package jsonic

import "errors"

// errors
var (
	ErrUnexpectedJSONData = errors.New("unexpected json data provided, neither array nor object")
	ErrIndexNotFound      = errors.New("expected index for json array but found something else")
	ErrIndexOutOfBound    = errors.New("index out of bounds of the json array")
	ErrNoDataFound        = errors.New("no tree satisfies the path elements provided")
	ErrInvalidType        = errors.New("data at the specified path does not match the expected type")
)
