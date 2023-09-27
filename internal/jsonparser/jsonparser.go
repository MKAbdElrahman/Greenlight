package jsonparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxBytes = 1_048_576

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a JSON decoder with disallowing unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := decodeJSON(dec, dst); err != nil {
		return handleJSONDecodeError(err)
	}

	if err := ensureSingleJSONValue(dec); err != nil {
		return err
	}

	return nil
}

func decodeJSON(dec *json.Decoder, dst interface{}) error {
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func handleJSONDecodeError(err error) error {

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown key %s", fieldName)
	case err.Error() == "http: request body too large":
		return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}
}

func ensureSingleJSONValue(dec *json.Decoder) error {
	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
