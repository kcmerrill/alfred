package remote

import (
	"strings"
)

type Remote struct {
	Repos []string
}

func New() *Remote {
	r := &Remote{}

	/* common repo is stored for official buisness ;) */
	r.Repos["common"] = "https://raw.githubusercontent.com/kcmerrill/alfred/master/modules/"
	return r
}

func (r *Remote) Exists(remote string) bool {
	/* is it a proper remote? If not, lets assume it's common */
	if !strings.Contains(remote, "/") {
		remote = "common/" + remote
	}

	/* does it exist? */
	if _, exists := r.Repos[remote]; exists {
		return true
	}

	/* boo! doesn't exist! */
	return false
}

func (r *Remote) URL(remote, module string) string {
	/* assuming it exists */
	if r.Exists(remote) {
		return r.Repos[remote] + module + "/alfred.yml"
	}

	/* hmmmmm ... can't find anything */
	return ""
}
