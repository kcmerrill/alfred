package alfred

import (
	"sync"
)

func _testAlfred() *Alfred {
	tasks := make(map[string]Task)
	return &Alfred{
		Tasks: tasks,
		Lock:  &sync.Mutex{},
	}
}
