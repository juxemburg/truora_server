package apierrors

/*ErrUnsupportedMediaType is an Unsupported Media Type*/
type ErrUnsupportedMediaType struct {
	message string
}

/*NewErrUnsupportedMediaType creates a new Unsupported Media Type*/
func NewErrUnsupportedMediaType(message string) *ErrUnsupportedMediaType {
	return &ErrUnsupportedMediaType{
		message: message,
	}
}

/*Error returns the error message*/
func (e *ErrUnsupportedMediaType) Error() string {
	return e.message
}
