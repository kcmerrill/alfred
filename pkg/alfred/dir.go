package alfred

import "path/filepath"

func (t *Task) dir(context *Context) (string, bool) {
	if t.Dir != "" {
		return mkdir(t.Dir, context)
	}

	return filepath.Clean(context.rootDir+"/") + "/", true
}
