package apierrors

/*ErrNotFound is an not found error*/
type ErrNotFound struct {
	message string
}

/*NewErrNotFound creates a new not found error*/
func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrNotFound) Error() string {
	return e.message
}
