package apierrors

/*ErrSQL is an not found error*/
type ErrSQL struct {
	message string
}

/*NewErrSQL creates a new not found error*/
func NewErrSQL(message string) *ErrSQL {
	return &ErrSQL{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrSQL) Error() string {
	return e.message
}
