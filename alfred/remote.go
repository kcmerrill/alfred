package alfred

import (
	"strings"
)

// Remote holds our remote mappings
type Remote struct {
	Repos map[string]string
}

// NewRemote create a new remote with our common defaults
func NewRemote(repos map[string]string) *Remote {
	r := &Remote{}
	r.Repos = repos

	// these repos are stored for official buisness ;)
	r.Repos["common"] = "https://raw.githubusercontent.com/kcmerrill/alfred/master/modules/"
	return r
}

// Exists returns true/false depending if a remote exists or not
func (r *Remote) Exists(remote string) bool {
	// does it exist?
	if _, exists := r.Repos[remote]; exists {
		return true
	}

	// boo! doesn't exist!
	return false
}

//URL returns the url given a specific remote
func (r *Remote) URL(remote, module string) string {
	// assuming it exists
	if r.Exists(remote) {
		return r.Repos[remote] + module + "/alfred.yml"
	}

	// hmmmmm ... can't find anything
	return ""
}

//ModulePath returns the path of a given remote
func (r *Remote) ModulePath(remote string) string {
	// is it a proper remote? If not, lets assume it's common
	if !strings.Contains(remote, "/") {
		remote = "common/" + remote
	}

	return remote
}

//Parse parses a string to see if it's a remote or not
func (r *Remote) Parse(remoteModule string) (string, string) {
	if strings.Contains(remoteModule, "/") {
		s := strings.SplitN(remoteModule, "/", 2)
		// Set some defaults
		if s[0] == "" {
			s[0] = "common"
		}
		if s[1] == "" {
			s[1] = "self"
		}
		return s[0], s[1]
	}

	return "", ""
}
