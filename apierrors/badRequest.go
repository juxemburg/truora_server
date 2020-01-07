package apierrors

/*ErrBadRequest is an not found error*/
type ErrBadRequest struct {
	message string
}

/*NewErrBadRequest creates a new not found error*/
func NewErrBadRequest(message string) *ErrBadRequest {
	return &ErrBadRequest{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrBadRequest) Error() string {
	return e.message
}
