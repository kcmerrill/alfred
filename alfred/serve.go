package alfred

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Serve is alfred's built in web server, useful for sharing private repos
func Serve(dir, port string) {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		say("error", err.Error())
	}
}
