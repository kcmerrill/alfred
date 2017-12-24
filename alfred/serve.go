package alfred

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// serve is alfred's built in web server, useful for sharing private repos
func serve(task Task, context *Context, tasks map[string]Task) {
	if task.Serve == "" {
		return
	}

	dir, _ := task.dir(context)

	// TODO taskdir task.dir()
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + task.Serve,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		outOK("serving "+dir, "0.0.0.0:"+task.Serve, context)
		if err := srv.ListenAndServe(); err != nil {
			outFail("serve", err.Error(), context)
			task.Exit(context, tasks)
		}
	}()

	// Pause, let the port open up before continuing on ...
	//<-time.After(time.Second)
}
