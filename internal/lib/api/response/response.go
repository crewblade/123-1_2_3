package response

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"errs,omitempty"`
}

func NewSuccess(successStatus int) Response {
	return Response{
		Status: successStatus,
	}
}

func NewError(errStatus int, msg string) Response {
	return Response{
		Status: errStatus,
		Error:  msg,
	}
}
