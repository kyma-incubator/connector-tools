package util

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RequestContext struct {
	TraceHeaders http.Header
}


//IncludeTraceHeaders enriches headers with Trace headers
func (ctx *RequestContext) IncludeTraceHeaders(dst http.Header) {
	for h, v := range ctx.TraceHeaders {
		dst[h] = v
	}
}

func (ctx *RequestContext) GetLoggerFields() log.Fields {
	return log.Fields{}
}