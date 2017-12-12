package alfred

import (
	"net/http"
	"testing"
)

func TestServeComponent(t *testing.T) {
	tasks := _sampleTasks()

	context := &Context{}

	go serve(tasks["http.serve"], context, _sampleTasks())

	response, err := http.Get("http://localhost:8080/serve_test.go")
	if err != nil {
		t.Fatalf("Unable to start HTTP webserver")
	}
	if response.StatusCode != 200 {
		t.Fatalf("Status code 200 expected")
	}
}
