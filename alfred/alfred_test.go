package alfred

import (
	"sync"
)

func _testAlfred() *Alfred {
	return &Alfred{
		Tasks: _testSampleTasks(),
		Lock:  &sync.Mutex{},
	}
}
