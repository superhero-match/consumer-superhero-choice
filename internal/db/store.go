package db

import (
	"github.com/consumer-superhero-choice/internal/db/model"
)

// StoreChoice saves newly registered Superhero.
func(db *DB) StoreChoice (c model.Choice) error {
	_, err := db.stmtInsertNewChoice.Exec(
		c.ID,
		c.Choice,
		c.SuperheroID,
		c.ChosenSuperheroID,
		c.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
