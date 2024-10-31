package errors

import "fmt"

func HandleError(err error) {
	fmt.Println("err: ", err)
}
