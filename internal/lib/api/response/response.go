package response

import "net/http"

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"errs,omitempty"`
}

func OK() Response {
	return Response{
		Status: http.StatusOK,
	}
}

func NewError(errStatus int, msg string) Response {
	return Response{
		Status: errStatus,
		Error:  msg,
	}
}
