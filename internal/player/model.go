package player

import "time"

type Player struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	Stats      *Stats    `json:"stats,omitempty"`      // stays nil for now
	AvatarData *string   `json:"avatarData,omitempty"` // NEW: base64 data URL
}

type Stats struct {
	MatchesPlayed int     `json:"matchesPlayed"`
	MatchesWon    int     `json:"matchesWon"`
	AverageScore  float64 `json:"averageScore"`
	BestCheckout  *int    `json:"bestCheckout,omitempty"`
}

// JSON body for POST /players
type CreatePlayerRequest struct {
	Name       string  `json:"name"`
	AvatarData *string `json:"avatarData,omitempty"` // NEW: optional image data URL
}

// JSON body for PUT /players/{id}
type UpdatePlayerRequest struct {
	Name       string  `json:"name"`
	AvatarData *string `json:"avatarData,omitempty"` // if nil -> keep existing avatar
}
