/*
  Copyright (C) 2019 - 2021 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package cache

import (
	"fmt"
	"github.com/superhero-match/consumer-superhero-choice/internal/cache/model"
)

// SetChoice stores choice(like only, dislikes only go to DB) into Redis cache.
func (c *Cache) SetLike(choice model.Choice) error {
	members := make([]interface{}, 0)
	members = append(members, choice.SuperheroID)

	err := c.Redis.SAdd(fmt.Sprintf(c.LikesKeyFormat, choice.ChosenSuperheroID), members..., ).Err()
	if err != nil {
		return err
	}

	return nil
}