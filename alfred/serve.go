package alfred

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func Serve(dir, port string) {
	r := mux.NewRouter()
	r.PathPrefix("/shares/").Handler(http.StripPrefix("/shares/", http.FileServer(http.Dir(dir))))

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
