package alfred

import (
	"sync"
)

// Alfred coordinates the tasks
type Alfred struct {
	Lock  *sync.Mutex
	Tasks map[string]Task
}

// Initialize will create and setup a new version of alfred
func Initialize(alfred *Alfred) *Alfred {
	// make sure we don't stomp on anything ....
	if alfred.Tasks == nil {
		alfred.Tasks = make(map[string]Task)
	}

	// no biggie to stomp on
	alfred.Lock = &sync.Mutex{}

	// return alfred
	return alfred
}
