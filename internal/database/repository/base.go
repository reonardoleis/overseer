package repository

import (
	"database/sql"

	"github.com/reonardoleis/overseer/internal/database"
)

type Repository struct {
	db      *sql.DB
	guildID string
	userID  string
}

func Prepare(guildID, userID string) *Repository {
	return &Repository{database.Conn, guildID, userID}
}
