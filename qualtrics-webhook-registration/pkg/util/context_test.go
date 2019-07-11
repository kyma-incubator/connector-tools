package util

import (
	"net/http"
	"testing"
)

func TestRequestContext_GetLoggerFields(t *testing.T) {
	//So far nothing to test
	context := RequestContext{
		TraceHeaders: http.Header{
			"Test": []string{"Temp"},
			"Hello": []string{"World"},
		},
	}


	fields := context.GetLoggerFields()

	if len(fields) != 0 {
		t.Errorf("no trace fields should have been returned, but there were %d", len(fields))
	}
}


func TestRequestContext_IncludeTraceHeaders(t *testing.T) {

	context := RequestContext{
		TraceHeaders: http.Header{
			"Test": []string{"Temp"},
			"Hello": []string{"World"},
		},
	}

	headers := http.Header{
		"Test": []string{"Temp2"},
		"Request": []string{"2"},
	}

	context.IncludeTraceHeaders(headers)

	if len(headers) != 3 {
		t.Fatalf("3 header fields should have been returned, but there were %d", len(headers))
	}

	if headers["Hello"][0] != "World" {
		t.Errorf("Expected header \"Hello\" has value \"World\", but has %q", headers["Hello"][0])
	}

	if headers["Test"][0] != "Temp" {
		t.Errorf("Expected header \"Test\" has value \"Temp\", but has %q", headers["Test"][0])
	}

	if headers["Request"][0] != "2" {
		t.Errorf("Expected header \"Request\" has value \"2\", but has %q", headers["Request"][0])
	}

}