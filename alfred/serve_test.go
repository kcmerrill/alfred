package alfred

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServeComponent(t *testing.T) {

	/* I think travis-ci.org isn't allowing ports anymore?
	Needs investgating ... and this test needs fixing if that's the case.
	For now, lets limit it to just my username /cringe
	*/

	if evaluate("whoami", "./") != "kcmerrill" {
		return
	}

	tasks := make(map[string]Task)
	tasks["http.serve"] = Task{
		Serve: "8088",
	}
	context := InitialContext([]string{})
	go serve(tasks["http.serve"], context, tasks)
	response, err := http.Get("http://localhost:8088/serve_test.go")
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("Expected a proper 200 response from localhost")
	}
	if response.StatusCode != 200 {
		t.Fatalf("Status code 200 expected")
	}
}
