package alfred

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	event "github.com/kcmerrill/hook"
)

// serve is alfred's built in web server, useful for sharing private repos
func serve(task Task, context *Context) {
	if task.Serve == "" {
		return
	}
	dir := "."
	if task.Dir != "" {
		dir = task.Dir
	}
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + task.Serve,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	event.Trigger("speak", "Serving "+dir+" 0.0.0.0:"+task.Serve, task, context)
	if err := srv.ListenAndServe(); err != nil {
		event.Trigger("speak", "{{ .Text.Failure }}"+err.Error()+"{{ .Text.Reset }}", task, context)
	}
}
