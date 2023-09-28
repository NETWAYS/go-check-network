package checkhttpmock

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Read all data from a io.ReadCloser, return the data as string and return a new io.ReadCloser to pass on
//
// This can be quite tricky and is only used for mocking and testing here.
func dumpAndBuffer(r io.ReadCloser) (string, io.ReadCloser) {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	_ = r.Close()

	return string(data), io.NopCloser(bytes.NewReader(data))
}

// Data structure to store information about a http.Request and http.Response in a simplified way
type Record struct {
	URL    string
	Method string
	Query  string
	Status string
	Body   string
}

// Build a new Record from an http.Request
func NewRecord(request *http.Request) (r *Record) {
	r = &Record{
		URL:    request.URL.String(),
		Method: request.Method,
	}

	// read the query from the request
	r.Query = extractFormQuery(request)

	slog.Info("recording request", "url", r.URL, "method", r.Method)

	return
}

// Update the Record with a http.Response to get Body and Status
func (r *Record) Complete(response *http.Response) {
	body, newReader := dumpAndBuffer(response.Body)
	response.Body = newReader

	r.Status = response.Status
	r.Body = body

	slog.Info("recording response", "status", response.Status)
}

// Write a YAML representation of the Record to an io.Writer
func (r Record) EmitYAML(w io.Writer) (err error) {
	out := yaml.NewEncoder(w)
	out.SetIndent(2)

	_, _ = fmt.Fprintln(w, "---")

	err = out.Encode(r)

	return
}
