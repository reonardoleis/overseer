package repository

import (
	"database/sql"
	"log"

	"github.com/reonardoleis/overseer/internal/database/models"
)

func (r *Repository) GetFunction(name string) (*models.Function, bool, error) {
	query := `SELECT  id, name, code, guild_id, user_id FROM functions 
            WHERE name = $1 AND
            guild_id = $2   AND
            user_id  = $3`

	row := r.db.QueryRow(query, name, r.guildID, r.userID)
	if row.Err() != nil && row.Err() == sql.ErrNoRows {
		return nil, false, nil
	} else if row.Err() != nil {
		log.Println("error getting function from database", row.Err())
		return nil, false, row.Err()
	}

	f := &models.Function{}
	err := row.Scan(&f.ID, &f.Name, &f.Code, &f.GuildID, &f.UserID)
	if err != nil {
		log.Println("error scanning function", err)
		return nil, true, err
	}

	return f, true, nil
}

func (r *Repository) CreateFunction(function *models.Function) (bool, error) {
	query := `SELECT COUNT(id) FROM functions 
            WHERE name = $1 AND
            guild_id = $2 AND
            user_id = $3`

	row := r.db.QueryRow(query, function.Name, r.guildID, r.userID)
	if row.Err() != nil {
		log.Println("error counting functions before creation", row.Err())
		return false, row.Err()
	}

	var count int64
	row.Scan(&count)

	if count > 0 {
		return true, nil
	}

	query = `INSERT INTO functions(name, code, guild_id, user_id)
           VALUES($1, $2, $3, $4)
           RETURNING id`

	err := r.db.
		QueryRow(query, function.Name, function.Code, r.guildID, r.userID).
		Scan(&function.ID)
	if err != nil {
		log.Println("error scanning function id", err)
		return false, err
	}

	return false, nil
}
