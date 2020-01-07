package apierrors

/*ErrUnauthorized is an not found error*/
type ErrUnauthorized struct {
	message string
}

/*NewErrUnauthorized creates a new not found error*/
func NewErrUnauthorized(message string) *ErrUnauthorized {
	return &ErrUnauthorized{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrUnauthorized) Error() string {
	return e.message
}
