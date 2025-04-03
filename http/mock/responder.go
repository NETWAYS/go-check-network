package checkhttpmock

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jarcoal/httpmock"
)

// Where response data is stored, relative to the package being tested
const TestData = "./testdata"
const contentTypeHeader = "Content-Type"
const contentTypeUrlencoded = "application/x-www-form-urlencoded"

// Extract a URL query from the request body, when the Content-Type is set to be urlencoded
func extractFormQuery(request *http.Request) string {
	if strings.Contains(request.Header.Get(contentTypeHeader), contentTypeUrlencoded) {
		query, newReader := dumpAndBuffer(request.Body)
		request.Body = newReader

		return query
	}

	return ""
}

// Mapping a partial form request to a response
//
// The response will be expected and read from the local testdata directory of your package.
//
//		QueryMap{
//	 	"test=1": "response.json",
//		}
type QueryMap map[string]string

// Register a NewQueryMapResponder with httpmock
//
// See QueryMap and NewQueryMapResponder
func RegisterQueryMapResponder(method, url string, queryMap QueryMap) {
	httpmock.RegisterResponder(method, url, NewQueryMapResponder(queryMap)) //nolint:bodyclose
}

// Return a responder function for httpmock, to return different results based on a QueryMap
//
// Queries from the QueryMap are matched partially and the response is read from `./testdata`
func NewQueryMapResponder(queryMap QueryMap) func(request *http.Request) (*http.Response, error) {
	return func(request *http.Request) (*http.Response, error) {
		query := extractFormQuery(request)

		for part, file := range queryMap {
			if strings.Contains(query, part) {
				body, err := os.ReadFile(filepath.Clean(path.Join(TestData, file)))
				return httpmock.NewStringResponse(200, string(body)), err
			}
		}

		// When a recorder is enabled use it - we don't have a way to access the NoResponder from here
		if CurrentRecorder != nil {
			return CurrentRecorder.Respond(request)
		}

		return nil, fmt.Errorf("no matching query found for: %s", query)
	}
}
