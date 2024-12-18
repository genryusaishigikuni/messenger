package memory

import (
	"sync"

	"github.com/genryusaishigikuni/messenger/presence-service/pkg/models"
)

type PresenceStore struct {
	mu       sync.RWMutex
	presence map[int]*models.Presence // user_id -> presence
}

func NewPresenceStore() *PresenceStore {
	return &PresenceStore{
		presence: make(map[int]*models.Presence),
	}
}

func (s *PresenceStore) SetOnline(userID int, channelID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.presence[userID] = &models.Presence{
		UserID:    userID,
		ChannelID: channelID,
		Online:    true,
	}
}

func (s *PresenceStore) SetOffline(userID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.presence, userID)
}

func (s *PresenceStore) GetAll() []models.Presence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []models.Presence
	for _, p := range s.presence {
		list = append(list, *p)
	}
	return list
}

func (s *PresenceStore) GetPresence(userID int) *models.Presence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.presence[userID]
	if !ok {
		return nil
	}
	return &models.Presence{
		UserID:    p.UserID,
		ChannelID: p.ChannelID,
		Online:    p.Online,
	}
}
