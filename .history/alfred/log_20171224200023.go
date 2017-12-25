package alfred

func log(task Task, context *Context, tasks map[string]Task) {
	if task.Log != "" {
		f, err := os.OpenFile(task.Log, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			context.Log[task.Log] = 
		}
	}
}

func logger(text string, context *Context) {

}
