package checkhttpmock

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jarcoal/httpmock"
)

func ExampleActivateRecorder() {
	// Activate the normal httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Activate recorder
	_ = os.Remove(RecordFile) // Remove any prior recording
	ActivateRecorder()

	// We don't set any mock examples here
	//httpmock.RegisterResponder("GET", "http://localhost:8080/test",
	//	func(request *http.Request) (*http.Response, error) {
	//		return httpmock.NewStringResponse(200, "Hello World"), nil
	//	})

	// Start a simple HTTP server
	runHTTP()

	// Do any HTTP request
	resp, err := http.Get("http://localhost:64888/test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s\n", data)

	// Print recording
	data, _ = os.ReadFile(RecordFile)
	fmt.Printf("%s\n", data)

	_ = resp.Body.Close()

	// Output:
	// Hello World
	// ---
	// url: http://localhost:64888/test
	// method: GET
	// query: ""
	// status: 200 OK
	// body: Hello World
}

func runHTTP() {
	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, `Hello World`)
	})

	go http.ListenAndServe(":64888", nil)
}
