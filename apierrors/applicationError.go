package apierrors

/*ApplicationError is an not found error*/
type ApplicationError struct {
	message string
}

/*NewApplicationError creates a new not found error*/
func NewApplicationError(message string) *ApplicationError {
	return &ApplicationError{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ApplicationError) Error() string {
	return e.message
}
