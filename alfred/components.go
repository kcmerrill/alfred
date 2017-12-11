package alfred

import (
	event "github.com/kcmerrill/hook"
)

func init() {
	// Speak will simply print out lines of text(for now)
	event.Register("speak", speak)

	// Summary will display the task summary
	event.Register("task.summary.header", summaryHeader)

	// Command will execute a command
	event.Register("task.command", command)

	// Register both ok/fail task groups
	event.Register("task.group", taskGroup)

	// Serve will serve up a static HTTP webserver
	event.Register("task.serve", serve)

	// SummaryFooter displays the summary footer
	event.Register("task.summary.footer", summaryFooter)

	// Wait a given duration
	event.Register("task.wait", wait)
}
