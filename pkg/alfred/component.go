package alfred

// Component contains a name and function
type Component struct {
	Name string
	F    func(task Task, context *Context, tasks map[string]Task)
}
