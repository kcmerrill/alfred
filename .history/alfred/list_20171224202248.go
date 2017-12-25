package alfred

import (
	"fmt"
	"sort"
)

func list(tasks map[string]Task) {
	max := 0
	labels := make([]string, 0)
	for label := range tasks {
		// lets add the label to the list(so we an alphabatize the list)
		labels = append(labels, label)

		// figure out max label size
		if len(label) > max {
			max = len(label)
		}
	}
	// alphabatize the list
	sort.Strings(labels)

	noLabels := 0
	// insignifigant tasks
	// still chewing on this one. Not sure if we should include them or not
	/*for _, label := range labels {
		task := tasks[label]
		if task.Summary == "" {
			noLabels++
			fmt.Print(translate("{{ .Text.Grey }}"+label+"{{ .Text.Reset }}", emptyContext()), "\t")
		}
	}*/

	if noLabels != 0 {
		// TODO: we need to determine if we should show this or not
		fmt.Println()
	}

	ec := emptyContext()
	// signifigant tasks
	for _, label := range labels {
		task := tasks[label]
		if task.Summary != "" {
			fmt.Println(translate("{{ .Text.Task }}"+padRight(label, max, " ")+"{{ .Text.Grey }} | {{ .Text.Reset }}"+task.Summary, ec))
		}
	}
}
