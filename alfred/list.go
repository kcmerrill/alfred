package alfred

import (
	"fmt"
	"sort"
	"strings"
)

func list(context *Context, tasks map[string]Task) {
	maxLabel := 0
	maxSummary := 0
	labels := make([]string, 0)
	for label, task := range tasks {
		// lets add the label to the list(so we an alphabatize the list)
		labels = append(labels, label)

		// figure out max label size
		if len(label) >= maxLabel {
			maxLabel = len(label)
		}

		// figure out max summary size
		if len(task.Summary) >= maxSummary {
			maxSummary = len(task.Summary)
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

	// signifigant tasks
	for _, label := range labels {
		task := tasks[label]
		if task.Summary != "" {
			fmt.Print(translate(" {{ .Text.Task }}"+padLeft(label, maxLabel, " "), context))
			fmt.Println(translate("{{ .Text.Grey }} | {{ .Text.Reset }}"+padRight(task.Summary, maxSummary, " "), context))
			if task.Usage != "" {
				usage := strings.Join(strings.Split(task.Usage, " "), "{{ .Text.Grey }}, {{ .Text.Args }}")
				usage = strings.Join(strings.Split(task.Usage, " "), "{{ .Text.Grey }},{{ .Text.Args }}")
				fmt.Println(translate(padLeft("", maxLabel, " ")+"{{ .Text.Args }}    "+usage+"", context))
				//fmt.Print(translate("{{ .Text.Grey }} ⇨ {{ .Text.Args }}"+usage+"", context))
				//fmt.Print(translate("{{ .Text.Grey }} | {{ .Text.Args }}("+usage+") ➡ {{ .Text.Reset }}"+task.Summary, context))
				//fmt.Print(translate("{{ .Text.Grey }} | {{ .Text.Args }}("+usage+") ⇨ {{ .Text.Reset }}"+task.Summary, context))
				//fmt.Print(translate("{{ .Text.Grey }} | {{ .Text.Args }}"+usage+" {{ .Text.Reset }}", context))
				//fmt.Print(translate("{{ .Text.Grey }} | {{ .Text.Reset }}"+task.Summary+"{{ .Text.Args }} ⇨ ("+usage+"){{ .Text.Reset }}", context))
				//fmt.Print(translate("{{ .Text.Grey }} | {{ .Text.Reset }}"+task.Summary+"{{ .Text.Args }}("+usage+"){{ .Text.Reset }}", context))
			}
		}

	}
}
