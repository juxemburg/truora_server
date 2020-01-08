package httpresult

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/juxemburg/truora_server/apierrors"
)

const (
	maxBodySize = 1048576
)

/*HandleRequestBody tries to process the request's body, returning an error otherwise*/
func HandleRequestBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	if r.Header.Get("Content-type") != "application/json" {
		return apierrors.NewErrUnsupportedMediaType("The Content-Type must be application/json")
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)

	if err == nil {
		return nil
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var msg string
	switch {
	case errors.As(err, &syntaxError):
		msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg = fmt.Sprintf("Request body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)
	case errors.Is(err, io.EOF):
		msg = "Request body must not be empty"
	case err.Error() == "http: request body too large":
		msg = "Request body must not be larger than 1MB"
	default:
		return err
	}
	return apierrors.NewErrBadRequest(msg)
}

/*HandleRequestResponse process every request made, listening for errors and returning
an adequate HTTP response*/
func HandleRequestResponse(w http.ResponseWriter, r *http.Request, fn func() (v interface{}, err error)) {
	w.Header().Set("Content-Type", "application/json")

	if result, err := fn(); err != nil {
		log.Println("error while handling the request: ", err.Error())
		switch err.(type) {
		case *apierrors.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 - %v", err.Error())))
		case *apierrors.ErrBadRequest:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 - %v", err.Error())))
		case *apierrors.ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("401 - %v", err.Error())))
		case *apierrors.ApplicationError:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 - %v", err.Error())))
		case *apierrors.ErrBadGateway:
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("502 - %v", err.Error())))
		case *apierrors.ErrSQL:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`500 - There was some issue with the server while handling 
			the request, pelase contact the administrator. 😔`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`500 - There was some issue with the server while handling 
			the request, pelase contact the administrator. 😔`))
		}
	} else {
		b, _ := json.Marshal(result)
		w.Write(b)
	}
}
