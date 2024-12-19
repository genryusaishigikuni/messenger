package memory

import (
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/models"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
	"strconv"
	"sync"
)

type PresenceStore struct {
	mu       sync.RWMutex
	presence map[int]*models.Presence // user_id -> presence
}

func NewPresenceStore() *PresenceStore {
	utils.Info("Initializing new PresenceStore")
	return &PresenceStore{
		presence: make(map[int]*models.Presence),
	}
}

func (s *PresenceStore) SetOnline(userID int, channelID int) {
	utils.Info("Setting user online: userID=" + strconv.Itoa(userID) + ", channelID=" + strconv.Itoa(channelID))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.presence[userID] = &models.Presence{
		UserID:    userID,
		ChannelID: channelID,
		Online:    true,
	}
	utils.Info("User set online: userID=" + strconv.Itoa(userID))
}

func (s *PresenceStore) SetOffline(userID int) {
	utils.Info("Setting user offline: userID=" + strconv.Itoa(userID))
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.presence, userID)
	utils.Info("User set offline: userID=" + strconv.Itoa(userID))
}

func (s *PresenceStore) GetAll() []models.Presence {
	utils.Info("Retrieving all online users")
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []models.Presence
	for _, p := range s.presence {
		list = append(list, *p)
	}
	utils.Info("Retrieved online users: count=" + strconv.Itoa(len(list)))
	return list
}

func (s *PresenceStore) GetPresence(userID int) *models.Presence {
	utils.Info("Fetching presence for userID=" + strconv.Itoa(userID))
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.presence[userID]
	if !ok {
		utils.Info("No presence found for userID=" + strconv.Itoa(userID))
		return nil
	}
	utils.Info("Presence found for userID=" + strconv.Itoa(userID))
	return &models.Presence{
		UserID:    p.UserID,
		ChannelID: p.ChannelID,
		Online:    p.Online,
	}
}
