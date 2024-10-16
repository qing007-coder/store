package rules

import "errors"

var InsufficientPermissions = errors.New("insufficient permissions")
var OtherError = errors.New("other error")
var NoMatchingPolicy = errors.New("no matching policy")
