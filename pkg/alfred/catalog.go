package alfred

import "strings"

func isCatalog(task string) bool {
	return strings.HasPrefix(task, "@")
}

func updateCatalog(dir string, context *Context) {
	// catalogs at this point are just git repositories.
	// assume a git repository, and update
	_, ok := execute("git pull", dir)
	if ok {
		outOK("@catalog", "updated!", context)
		return
	}
	outWarn("@catalog", "Unable to update the catalog. It could be out of date.", context)
}
