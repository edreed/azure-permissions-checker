package controller

import "fmt"

type invalidObjectStateErr struct {
}

var InvalidObjectStateErr error = &invalidObjectStateErr{}

func (e *invalidObjectStateErr) Error() string {
	return "object is not in a valid state"
}

type invalidResourcePathErr struct {
	resourcePath string
}

func NewInvalidResourcePathErr(resourcePath string) error {
	return &invalidResourcePathErr{resourcePath: resourcePath}
}

func (e *invalidResourcePathErr) Error() string {
	return fmt.Sprintf("the resource path %q is invalid", e.resourcePath)
}
