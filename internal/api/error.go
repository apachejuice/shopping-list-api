package api

// used to differentiate between error values that warrant a 500 and ones that indicate a malformed request
type ApiError struct {
	error
	code          int
	isServerError bool
}

func NewApiError(err error, isServerError bool) *ApiError {
	return &ApiError{error: err, isServerError: isServerError, code: -1}
}

func NewApiErrorWithCode(err error, code int) *ApiError {
	return &ApiError{error: err, code: code, isServerError: false}
}

func (a ApiError) IsServerError() bool {
	return a.isServerError
}

func (a ApiError) GetCode() int {
	return a.code
}
