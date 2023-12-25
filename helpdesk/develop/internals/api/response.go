package api

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ResponseData any

var (
	ResponseSuccess = true
	ResponseFail    = false
)

type Response struct {
	Status bool         `json:"status"`
	Data   ResponseData `json:"data"`
	Error  error        `json:"error"`
}

func NewResponse(status bool, data ResponseData, err error) Response {
	return Response{
		Status: status,
		Data:   data,
		Error:  err,
	}
}

func Success(data any) Response {
	return NewResponse(ResponseSuccess, data, nil)
}

func Fail(err error) Response {
	return NewResponse(
		ResponseFail, nil, err)
}
