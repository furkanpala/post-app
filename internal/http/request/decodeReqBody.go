package request

import (
	"encoding/json"
	"net/http"

	httperror "github.com/furkanpala/post-app/internal/http/error"
)

// DecodeRequestBody function decodes json request r and stores to the value pointed by v
func DecodeRequestBody(r *http.Request, v interface{}) *httperror.HTTPError {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Invalid JSON",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

	return nil
}
