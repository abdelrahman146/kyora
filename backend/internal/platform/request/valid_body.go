package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ValidBody(c *gin.Context, obj any) error {
	if c.Request == nil || c.Request.Body == nil {
		err := errors.New("request body is required")
		response.Error(c, problem.BadRequest("invalid request body").WithError(err).WithCode("request.invalid_body"))
		return err
	}

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(obj); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			response.Error(c, problem.PayloadTooLarge("request body too large").WithError(err).WithCode("request.body_too_large"))
			return err
		}

		if errors.Is(err, io.EOF) {
			err = errors.New("request body is required")
		}
		response.Error(c, problem.BadRequest("invalid request body").WithError(err).WithCode("request.invalid_body"))
		return err
	}

	// Reject trailing JSON tokens (e.g., `{...}{...}` or `{...} garbage`).
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("invalid request body")
		}
		response.Error(c, problem.BadRequest("invalid request body").WithError(err).WithCode("request.invalid_body"))
		return err
	}

	if binding.Validator != nil {
		if err := binding.Validator.ValidateStruct(obj); err != nil {
			response.Error(c, problem.BadRequest("invalid request body").WithError(err).WithCode("request.invalid_body"))
			return err
		}
	}

	return nil
}
