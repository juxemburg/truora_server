package apierrors

/*ErrBadGateway is an not found error*/
type ErrBadGateway struct {
	message string
}

/*NewErrBadGateway creates a new bad gateway error*/
func NewErrBadGateway(message string) *ErrBadGateway {
	return &ErrBadGateway{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrBadGateway) Error() string {
	return e.message
}
