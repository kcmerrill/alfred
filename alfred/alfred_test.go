package alfred

import (
	"sync"
)

func _testAlfred() *Alfred {
	return &Alfred{
		Tasks: _sampleTasks(),
		Lock:  &sync.Mutex{},
	}
}
