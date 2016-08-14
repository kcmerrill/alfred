package remote

import (
	"strings"
)

type Remote struct {
	Repos map[string]string
}

func New(repos map[string]string) *Remote {
	r := &Remote{}
	r.Repos = repos

	/* these repos are stored for official buisness ;) */
	r.Repos["common"] = "https://raw.githubusercontent.com/kcmerrill/alfred/master/modules/"
	r.Repos["official"] = "https://raw.githubusercontent.com/kcmerrill/alfred/master/modules/"
	r.Repos["alfred"] = "https://raw.githubusercontent.com/kcmerrill/alfred/master/modules/"
	return r
}

func (r *Remote) Exists(remote string) bool {
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

func (r *Remote) ModulePath(remote string) string {
	/* is it a proper remote? If not, lets assume it's common */
	if !strings.Contains(remote, "/") {
		remote = "common/" + remote
	}

	return remote
}

func (r *Remote) Parse(remoteModule string) (string, string) {
	if strings.Contains(remoteModule, "/") {
		s := strings.Split(remoteModule, "/")
		return s[0], s[1]
	}

	return "", ""
}
