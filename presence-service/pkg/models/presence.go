package models

type Presence struct {
	UserID    int  `json:"user_id"`
	ChannelID int  `json:"channel_id"`
	Online    bool `json:"online"`
}
