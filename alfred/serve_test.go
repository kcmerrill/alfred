package alfred

import (
	"net/http"
	"testing"
)

func TestServeComponent(t *testing.T) {
	tasks := make(map[string]Task)
	tasks["http.serve"] = Task{
		Serve: "8088",
	}
	context := _testSilentContext()
	go serve(tasks["http.serve"], context, tasks)
	response, _ := http.Get("http://localhost:8080/serve_test.go")
	if response.StatusCode != 200 {
		t.Fatalf("Status code 200 expected")
	}
}
