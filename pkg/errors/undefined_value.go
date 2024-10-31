package errors

import "errors"

func UndefinedValue(value string) error {
	return errors.New("undefined " + value + " error")
}
