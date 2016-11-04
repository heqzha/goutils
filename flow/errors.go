package flow

// Most of errors.go codes are copied from gin-gonic/gin/errors.go
// Thanks for their great works!

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// ErrorType describes error types
type ErrorType uint64

const (
	ErrorTypePrivate ErrorType = 1 << 0
	ErrorTypePublic  ErrorType = 1 << 1

	ErrorTypeAny ErrorType = 1<<64 - 1
	ErrorTypeNu            = 2
)

type (
	//Error is error struct with ErrorType and Meta data
	Error struct {
		Err  error
		Type ErrorType
		Meta interface{}
	}

	errors []*Error
)

var _ error = &Error{} //inherit??

func (msg *Error) SetType(flags ErrorType) *Error {
	msg.Type = flags
	return msg
}

func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

func (msg *Error) JSON() interface{} {
	json := H{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				json[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			json["meta"] = msg.Meta
		}
	}
	if _, ok := json["error"]; !ok {
		json["error"] = msg.Error()
	}
	return json
}

// Implements the json.Marshaller interface
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

// Implements the error interface
func (msg *Error) Error() string {
	return msg.Err.Error()
}

func (msg *Error) IsType(flags ErrorType) bool {
	return (msg.Type & flags) > 0
}

// Returns a readonly copy filterd the byte.
// ie ByType(gin.ErrorTypePublic) returns a slice of errors with type=ErrorTypePublic
func (a errors) ByType(typ ErrorType) errors {
	if len(a) == 0 {
		return nil
	}
	if typ == ErrorTypeAny {
		return a
	}
	var result errors = nil
	for _, msg := range a {
		if msg.IsType(typ) {
			result = append(result, msg)
		}
	}
	return result
}

// Returns the last error in the slice. It returns nil if the array is empty.
// Shortcut for errors[len(errors)-1]
func (a errors) Last() *Error {
	length := len(a)
	if length > 0 {
		return a[length-1]
	}
	return nil
}

// Returns an array will all the error messages.
// Example:
//		c.Error(errors.New("first"))
//		c.Error(errors.New("second"))
//		c.Error(errors.New("third"))
//		c.Errors.Errors() // == []string{"first", "second", "third"}
func (a errors) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

func (a errors) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

func (a errors) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

func (a errors) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", (i + 1), msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
