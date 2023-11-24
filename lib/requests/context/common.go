package context

import "github.com/tpyle/reqx/lib/requests/context/http"

type RequestContext struct {
	FileLocation string
	HTTPContext  http.HTTPRequestContext
}
