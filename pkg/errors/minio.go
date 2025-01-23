package errors

import "errors"

var MCCreateError = errors.New("minio create error")
var MCUpdateError = errors.New("minio update error")
var MCQueryError = errors.New("minio query error")
var MCDeleteError = errors.New("minio delete error")
