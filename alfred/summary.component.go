package alfred

import "fmt"

// SummaryComponent prints out the task summary
func (t *Task) SummaryComponent(context Context, tasks Tasks) Context {
	fmt.Println("Summary:", t.Summary)
	context.Bleh += "woot,"
	return context
}
