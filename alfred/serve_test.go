package alfred

import (
	"net/http"
	"testing"
)

func TestServeComponent(t *testing.T) {
	tasks := _testSampleTasks()

	context := _testSilentContext()

	serve(tasks["http.serve"], context, _testSampleTasks())

	response, _ := http.Get("http://localhost:8080/serve_test.go")
	if response.StatusCode != 200 {
		t.Fatalf("Status code 200 expected")
	}
}
