package discord

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	guildID             string
	audioQueue          *playableItemQueue
	lastInteractionTime time.Time
	vc                  *discordgo.VoiceConnection
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

	m, exists := managers[guildID]
	if exists {
		m.lastInteractionTime = time.Now()
	}

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
			queue:   make(chan *playableItem, 1024),
		},
		lastInteractionTime: time.Now(),
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

func removeManager(guildID string) {
	managersLock.Lock()
	defer managersLock.Unlock()

	delete(managers, guildID)
}

func (m *Manager) setVC(vc *discordgo.VoiceConnection) {
	m.vc = vc
	go handleVoice(m.vc.OpusRecv, m)
}

func (m *Manager) audioPlayerWorker() {
	m.audioQueue.audioPlayerWorker(m.vc)
}

func managerCleanupJob() {
	for {
		managersLock.Lock()
		for guildID, manager := range managers {
			if time.Since(manager.lastInteractionTime) > time.Minute*5 {
				log.Println("cleaning up manager for guild", guildID)
				delete(managers, guildID)
			}
		}

		managersLock.Unlock()

		time.Sleep(time.Second * 10)
	}
}
