package types

import (
	"errors"
)

var (
	ERR_NO_SESSION         = errors.New("ERR_NO_SESSION")
	ERR_NO_SERVICE         = errors.New("ERR_NO_SERVICE")
	ERR_NO_SERVICE_CLIENT  = errors.New("ERR_NO_SERVICE_CLIENT")
	ERR_INVALID_SERVICE_ID = errors.New("ERR_INVALID_SERVICE_ID")
)
