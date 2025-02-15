package errors

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Handler creates a middleware that handles panics and errors encountered during HTTP request processing.
func Handler(logger *zap.Logger) routing.Handler {
	return func(c *routing.Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				var ok bool
				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}

				logger.Error("recovered from panic.", zap.String("error", err.Error()),
					zap.String("stack", string(debug.Stack())))
			}

			if err != nil {
				res := buildErrorResponse(err)
				if res.StatusCode() == http.StatusInternalServerError {
					logger.Error("encountered internal server.", zap.String("error", err.Error()))
				}
				c.Response.WriteHeader(res.StatusCode())
				if err = c.Write(res); err != nil {
					logger.Error("failed writing error response.", zap.String("error", err.Error()))
				}
				c.Abort() // skip any pending handlers since an error has occurred
				err = nil // return nil because the error is already handled
			}
		}()
		return c.Next()
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err error) ErrorResponse {
	switch err.(type) {
	case ErrorResponse:
		return err.(ErrorResponse)
	case validation.Errors:
		return InvalidInput(err.(validation.Errors))
	case routing.HTTPError:
		switch err.(routing.HTTPError).StatusCode() {
		case http.StatusNotFound:
			return NotFound("")
		default:
			return ErrorResponse{
				Status:  err.(routing.HTTPError).StatusCode(),
				Message: err.Error(),
			}
		}
	default:
		if err == sql.ErrNoRows {
			return NotFound("")
		}
		return InternalServerError("")
	}
}
