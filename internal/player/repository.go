package player

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrPlayerHasGames = errors.New("player has games and cannot be deleted")
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
	name = strings.ToUpper(name)

	var p Player
	err := r.db.QueryRow(ctx, `
INSERT INTO players (name, avatar_data)
VALUES ($1, $2)
RETURNING id::text, name, created_at, avatar_data;
`, name, avatarData).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.AvatarData)

	return p, err
}

func (r *Repository) UpdatePlayer(ctx context.Context, id, name string, avatarData *string) (Player, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Player{}, errors.New("name cannot be empty")
	}
	name = strings.ToUpper(name)

	var p Player
	err := r.db.QueryRow(ctx, `
UPDATE players
SET name = $2,
    avatar_data = COALESCE($3, avatar_data)
WHERE id = $1
RETURNING id::text, name, created_at, avatar_data;
`, id, name, avatarData).Scan(&p.ID, &p.Name, &p.CreatedAt, &p.AvatarData)

	if err != nil {
		return Player{}, err
	}

	return p, nil
}

func (r *Repository) DeletePlayer(ctx context.Context, id string) error {
	cmdTag, err := r.db.Exec(ctx, `
DELETE FROM players
WHERE id = $1;
`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			// foreign key violation: player is referenced by game_players
			return ErrPlayerHasGames
		}
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrPlayerNotFound
	}

	return nil
}
