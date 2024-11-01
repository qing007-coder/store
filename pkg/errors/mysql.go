package errors

import "errors"

var DBInitFailed = errors.New("db init failed")
var AutoMigrateFailed = errors.New("auto migrate failed")
var RecordNotFound = errors.New("record not found")

var DBCreateError = errors.New("db create error")
var DBUpdateError = errors.New("db update error")
var DBQueryError = errors.New("db query error")
var DBDeleteError = errors.New("db delete error")
