package alfred

// Tasks will hold all of our raw tasks
type Tasks struct {
	Raw map[string]Task
}

// Get returns a task based on a name, false if no task is found
func (t *Tasks) Get(name string) (Task, bool) {
	if task, exists := t.Raw[name]; exists {
		return task, true
	}
	return Task{}, false
}
