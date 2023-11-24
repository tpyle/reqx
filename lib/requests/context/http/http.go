package http

import "time"

type HTTPRequestContext struct {
	Timeout       time.Duration
	SkipTLSVerify bool
}
