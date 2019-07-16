package httphandler

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	requestID    = "X-Request-Id"
	traceId      = "X-B3-Traceid"
	spanId       = "X-B3-Spanid"
	parentSpanId = "X-b3-Parentspanid"
	sampled      = "X-B3-Sampled"
	flags        = "X-B3-Flags"
	spanContext  = "X-Ot-Span-Context"
)

type Handler interface {
	HandleRequest(r *http.Request, ctx *RequestContext) *Response
}

type RequestContext struct {
	TraceHeaders http.Header
}

type Response struct {
	IsSuccess    bool
	ResponseCode int
	Response     interface{}
}

type HandlerContext struct {
	Metrics     *Metrics
	NextHandler Handler
}

type Metrics struct {
	HTTPCalls2xx        prometheus.Counter
	HTTPCalls3xx        prometheus.Counter
	HTTPCalls4xx        prometheus.Counter
	HTTPCalls5xx        prometheus.Counter
	ServerResponseTimes prometheus.Summary
	InFlightRequests    prometheus.Gauge
}

type JsonError struct {
	Message string `json:"message"`
}

//ToJSON converts Response to json
func (r *Response) ToJSON() []byte {
	bytes, err := json.Marshal(r.Response)

	if err != nil {
		return []byte("{\"message:\": \"Error converting error message to json format, hence this is meaningless :-(\"}")
	}

	return bytes
}

// extracts istio trace headers for request context
// (https://istio.io/docs/tasks/telemetry/distributed-tracing/overview/)
func extractTraceHeaders(src http.Header, dst http.Header) {
	if val, ok := src[requestID]; ok {
		dst[requestID] = val
	}
	if val, ok := src[traceId]; ok {
		dst[traceId] = val
	}
	if val, ok := src[spanId]; ok {
		dst[spanId] = val
	}
	if val, ok := src[parentSpanId]; ok {
		dst[parentSpanId] = val
	}
	if val, ok := src[sampled]; ok {
		dst[sampled] = val
	}
	if val, ok := src[flags]; ok {
		dst[flags] = val
	}
	if val, ok := src[spanContext]; ok {
		dst[spanContext] = val
	}
}

//IncludeTraceHeaders enriches headers with Trace headers
func (ctx *RequestContext) IncludeTraceHeaders(dst http.Header) {
	for h, v := range ctx.TraceHeaders {
		dst[h] = v
	}
}

func (ctx *RequestContext) GetLoggerFields() log.Fields {
	return log.Fields{
		requestID: ctx.TraceHeaders[requestID],
		traceId:   ctx.TraceHeaders[traceId],
		spanId:    ctx.TraceHeaders[spanId],
	}
}

func (h *HandlerContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//Collect Duration Metric
	startTime := time.Now()

	// Set Content Type to "application/json in any case
	w.Header().Add("Content-Type", "application/json")

	ctx := RequestContext{
		TraceHeaders: http.Header{},
	}

	if log.GetLevel() == log.TraceLevel {
		r.ParseForm()
		log.WithFields(
			ctx.GetLoggerFields(),
		).Tracef("Request for %s processed received: %+v", r.URL.Opaque, r.Form)
	}

	extractTraceHeaders(r.Header, ctx.TraceHeaders)

	//Make request and manage in-flight
	h.Metrics.InFlightRequests.Inc()
	resp := h.NextHandler.HandleRequest(r, &ctx)
	h.Metrics.InFlightRequests.Dec()

	if resp == nil {
		_, _ = w.Write([]byte("{\"message:\": \"Internal Server Error, please contact an administrator\"}"))
		w.WriteHeader(500)

		//if not level trace
		if log.GetLevel() != log.TraceLevel {
			log.WithFields(
				ctx.GetLoggerFields(),
			).Errorf("Request for %s processed with fatal error (nil response from Handler)", r.URL.Opaque)
		} else {

			//we are already doomed, so no need for additional error handling
			requestBody, _ := ioutil.ReadAll(r.Body)
			_ = r.Body.Close()

			log.WithFields(
				ctx.GetLoggerFields(),
			).WithFields(log.Fields{
				"requestBody": string(requestBody),
			}).Tracef("Request for %s processed with fatal error (nil response from Handler)", r.URL.Opaque)
		}
		//not to forget about metric
		h.Metrics.HTTPCalls5xx.Inc()
		h.Metrics.ServerResponseTimes.Observe(float64(time.Since(startTime)))
		return
	}

	//deal with metrics
	if resp.ResponseCode > 199 && resp.ResponseCode < 300 {
		h.Metrics.HTTPCalls2xx.Inc()
	} else if resp.ResponseCode > 299 && resp.ResponseCode < 400 {
		h.Metrics.HTTPCalls3xx.Inc()
	} else if resp.ResponseCode > 399 && resp.ResponseCode < 500 {
		h.Metrics.HTTPCalls4xx.Inc()
	} else if resp.ResponseCode > 499 && resp.ResponseCode < 600 {
		h.Metrics.HTTPCalls5xx.Inc()
	}

	w.WriteHeader(resp.ResponseCode)
	_, _ = w.Write(resp.ToJSON())
	h.Metrics.ServerResponseTimes.Observe(float64(time.Since(startTime)))

	if log.GetLevel() != log.TraceLevel {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Debugf("Request for %s processed with response code %d", r.URL.Opaque, resp.ResponseCode)
	} else {

		log.WithFields(
			ctx.GetLoggerFields(),
		).Tracef("Request processed with response code %d: %1v", resp.ResponseCode, resp)
	}

}
