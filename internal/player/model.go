package player

import "time"

type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Stats     *Stats    `json:"stats,omitempty"` // stays nil for now
}

type Stats struct {
	MatchesPlayed int     `json:"matchesPlayed"`
	MatchesWon    int     `json:"matchesWon"`
	AverageScore  float64 `json:"averageScore"`
	BestCheckout  *int    `json:"bestCheckout,omitempty"`
}

// JSON body for POST /players
type CreatePlayerRequest struct {
	Name string `json:"name"`
}
