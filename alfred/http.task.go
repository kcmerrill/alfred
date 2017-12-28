package alfred

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// HTTPTask contains all of our task information
type HTTPTask struct {
	context *Context
	tasks   map[string]Task
}

func (h *HTTPTask) runner(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	c := h.context
	c.Args = strings.Split(vars["args"], "/")
	outOK(vars["task"]+" started ["+strings.Join(c.Args, ", ")+"]", "", c)
	stdin, _ := ioutil.ReadAll(req.Body)
	c.Stdin = strings.TrimSpace(string(stdin))
	c.Out = resp
	c.Text = TextConfig{}
	NewTask(vars["task"], c, h.tasks)
}

func httptasks(task Task, context *Context, tasks map[string]Task) {
	if task.HTTPTasks.Port == "" {
		return
	}

	dir, _ := task.dir(context)

	password := translate(evaluate(task.HTTPTasks.Password, dir), context)

	addr := []string{"0.0.0.0", task.HTTPTasks.Port}
	if strings.Contains(task.HTTPTasks.Port, ":") {
		addr = strings.Split(task.HTTPTasks.Port, ":")
	}

	runner := &HTTPTask{
		context: context,
		tasks:   tasks,
	}

	r := mux.NewRouter()
	r.HandleFunc(`/{task}/{args:[a-zA-Z0-9=\-\/\.]+}`, httptasksAuth(password, runner.runner)).Methods("GET", "POST")
	r.HandleFunc(`/{task}`, httptasksAuth(password, runner.runner)).Methods("GET", "POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         strings.Join(addr, ":"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	outOK("serving "+dir, strings.Join(addr, ":"), context)
	if err := srv.ListenAndServe(); err != nil {
		outFail("serve", err.Error(), context)
		task.Exit(context, tasks)
	}
}

func httptasksAuth(password string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := r.BasicAuth()
		if password != "" && password != token {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
