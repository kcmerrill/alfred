package alfred

func (t *Task) dir(context *Context) (string, bool) {
	if t.Dir != "" {
		return mkdir(t.Dir, context)
	}

	if context.rootDir == "" {
		return "./", true
	}

	return context.rootDir, true
}
