// Package merr provides error handling for the application.
package merr

import "github.com/go-kratos/kratos/v2/errors"

// ErrorParamsNotSupportFileConfig returns an error indicating that the operation is not supported in file config mode.
func ErrorParamsNotSupportFileConfig() *errors.Error {
	return ErrorParams("operation is not supported in file config mode")
}
