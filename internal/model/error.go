package model

type ErrorNotFound struct {
	Msg string
}

func (e ErrorNotFound) Error() string {
	return e.Msg
}

func NewErrorNotFound(msg string) ErrorNotFound {
	return ErrorNotFound{Msg: msg}
}

type ErrorUnAuthorized struct {
	Msg string
}

func (e ErrorUnAuthorized) Error() string {
	return e.Msg
}

func NewErrorUnAuthorized(msg string) ErrorUnAuthorized {
	return ErrorUnAuthorized{Msg: msg}
}

