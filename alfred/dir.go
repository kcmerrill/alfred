package alfred

import (
	"os"
)

func (t *Task) dir(context *Context) (string, bool) {
	if t.Dir != "" {
		dir := evaluate(t.Template(t.Dir, context), ".")
		if _, err := os.Stat(dir); err == nil {
			// woot!
			return dir, true
		}

		// ok, we have some work to do
		if err := os.MkdirAll(dir, 0755); err != nil {
			// problem making directory
			return "", false
		}
		return dir, true
	}
	return ".", true
}
