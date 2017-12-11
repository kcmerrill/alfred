package alfred

import (
	"net/http"
	"testing"
)

func TestServeComponent(t *testing.T) {
	tasks := _sampleTasks()
	tasks["http.server"] = Task{
		Serve: "8088",
	}
	context := &Context{}

	NewTask("http.server", context, tasks)

	response, err := http.Get("http://localhost:8088/main.go")
	if err != nil {
		t.Fatalf("Unable to start HTTP webserver")
	}
	if response.StatusCode != 200 {
		t.Fatalf("Status code 200 expected")
	}
}
