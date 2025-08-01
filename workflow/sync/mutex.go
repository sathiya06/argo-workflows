package sync

import (
	sema "golang.org/x/sync/semaphore"
)

// newInternalMutex creates a size 1 semaphore
func newInternalMutex(name string, nextWorkflow NextWorkflow) *prioritySemaphore {
	logger := syncLogger{
		name:     name,
		lockType: lockTypeMutex,
	}
	return &prioritySemaphore{
		name:         name,
		limitGetter:  &mutexLimit{},
		pending:      &priorityQueue{itemByKey: make(map[string]*item)},
		semaphore:    sema.NewWeighted(int64(1)),
		lockHolder:   make(map[string]bool),
		nextWorkflow: nextWorkflow,
		logger:       logger.get,
	}
}
