package model

// Error Not Found
type ErrorNotFound struct {
	Msg string
}

func (e ErrorNotFound) Error() string {
	return e.Msg
}

func NewErrorNotFound(msg string) ErrorNotFound {
	return ErrorNotFound{Msg: msg}
}

// Error Unauthorized
type ErrorUnAuthorized struct {
	Msg string
}

func (e ErrorUnAuthorized) Error() string {
	return e.Msg
}

func NewErrorUnAuthorized(msg string) ErrorUnAuthorized {
	return ErrorUnAuthorized{Msg: msg}
}

// Error Bad request
type ErrorBadRequest struct {
	Msg string
}

func (e ErrorBadRequest) Error() string {
	return e.Msg
}

func NewErrorBadRequest(msg string) ErrorBadRequest {
	return ErrorBadRequest{Msg: msg}
}
