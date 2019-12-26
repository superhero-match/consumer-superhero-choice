package model

// Choice holds the information about the choice that a user made, e.g. like or dislike.
type Choice struct {
	ID                string `db:"id"`
	Choice            int64  `db:"choice"`
	SuperheroID       string `db:"superhero_id"`
	ChosenSuperheroID string `db:"chosen_superhero_id"`
	CreatedAt         string `db:"created_at"`
}
