package errors

import "errors"

var ConfigFileNotFound = errors.New("config file not found")
var OtherError = errors.New("other error")
var UnmarshalError = errors.New("unmarshal error")
