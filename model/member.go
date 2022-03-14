package model

import "time"

// Member 群成员
type Member struct {
	GuildID  string    `json:"guild_id"`
	JoinedAt time.Time `json:"joined_at"`
	Nick     string    `json:"nick"`
	User     *User     `json:"user"`
	Roles    []string  `json:"roles"`
}
