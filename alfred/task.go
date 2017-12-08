package alfred

// Task holds all of our task components
type Task struct {
	Summary     string
	Description string
	Dir         string
	Command     string
	Script      string
}

// NewTask will create a new task
func NewTask(name string, tasks Tasks) {
	// does the task exist?
	if task, exists := tasks.Get(); exists {

	}
	// TODO: Task does not exist
}
