package cache

import (
	"fmt"
	"time"

	"github.com/consumer-superhero-choice/internal/cache/model"
)

// SetChoice stores choice(like only, dislikes only go to DB) into Redis cache.
func (c *Cache) SetChoice(choice model.Choice) error {
	err := c.Redis.Set(
		fmt.Sprintf("choice.%s.%s", choice.SuperheroID, choice.ChosenSuperheroID),
		choice,
		time.Hour*time.Duration(24)*time.Duration(7),
	).Err()
	if err != nil {
		return err
	}

	return nil
}
