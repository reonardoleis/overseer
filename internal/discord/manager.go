package discord

import "sync"

type Manager struct {
	guildID        string
	audioQueue     *playableItemQueue
	audioQueueLock sync.Mutex
}

var (
	managersLock = sync.Mutex{}
	managers     = map[string]*Manager{} // guildID -> Manager
)

func managerExists(guildID string, useLock bool) bool {
	if useLock {
		managersLock.Lock()
		defer managersLock.Unlock()
	}

	_, exists := managers[guildID]
	return exists
}

func createManager(guildID string, useLock bool) *Manager {
	if useLock {
		managersLock.Lock()
		defer managersLock.Unlock()
	}

	manager := &Manager{
		guildID: guildID,
		audioQueue: &playableItemQueue{
			guildID: guildID,
			queue:   make(chan playableItem, 1024),
		},
	}

	managers[guildID] = manager

	return manager
}

func getManager(guildID string) *Manager {
	managersLock.Lock()
	defer managersLock.Unlock()

	if !managerExists(guildID, false) {
		createManager(guildID, false)
	}

	return managers[guildID]
}
