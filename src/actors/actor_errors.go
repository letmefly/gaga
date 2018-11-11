package actors

import (
	"errors"
)

var (
	ERR_FUNCTION_INVALID   = errors.New("ERR_INVALID_FUNCTION")
	ERR_FUNCTION_TYPE      = errors.New("ERR_FUNCTION_TYPE")
	ERR_NO_SERVICE_CLIENT  = errors.New("ERR_NO_SERVICE_CLIENT")
	ERR_INVALID_SERVICE_ID = errors.New("ERR_INVALID_SERVICE_ID")
)
