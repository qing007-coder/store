package errors

import "errors"

var JsonMarshalError = errors.New("json marshal error")
var JsonUnmarshalError = errors.New("json unmarshal error")
var ShouldBindJsonError = errors.New("bind json error")
