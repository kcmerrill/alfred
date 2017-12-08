package alfred

// Tasks will hold all of our raw tasks
type Tasks struct {
	Raw map[string]Task
}

func (t *Tasks) Get(name string) (Task, bool) {}
