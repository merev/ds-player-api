package player

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPlayers(ctx context.Context) ([]Player, error) {
	rows, err := r.db.Query(ctx, `
SELECT id::text, name, created_at, avatar_data
FROM players
ORDER BY created_at ASC;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Player, 0)

	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.AvatarData); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

func (r *Repository) CreatePlayer(ctx context.Context, name string, avatarData *string) (Player, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Player{}, errors.New("name cannot be empty")
	}

	var p Player
	err := r.db.QueryRow(ctx, `
INSERT INTO players (name, avatar_data)
VALUES ($1, $2)
RETURNING id::text, name, created_at, avatar_data;
`, name, avatarData).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.AvatarData)

	return p, err
}
