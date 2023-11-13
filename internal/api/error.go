package api

// used to differentiate between error values that warrant a 500 and ones that indicate a malformed request
type ApiError struct {
	error
	isServerError bool
}

func NewApiError(err error, isServerError bool) *ApiError {
	return &ApiError{error: err, isServerError: isServerError}
}

func (a ApiError) IsServerError() bool {
	return a.isServerError
}
