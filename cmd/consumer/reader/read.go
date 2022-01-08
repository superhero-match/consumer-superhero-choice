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
package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	cache "github.com/superhero-match/consumer-superhero-choice/internal/cache/model"
	"github.com/superhero-match/consumer-superhero-choice/internal/consumer/model"
	dbm "github.com/superhero-match/consumer-superhero-choice/internal/db/model"
)

const like = int64(1)

// Read consumes the Kafka topic and stores the choice made by superhero
// to DB and if it is a like to Cache as well.
func (r *reader) Read() error {
	ctx := context.Background()

	for {
		fmt.Print("before FetchMessage")
		m, err := r.Consumer.FetchMessage(ctx)
		fmt.Print("after FetchMessage")
		if err != nil {
			r.Logger.Error(
				"failed to fetch message",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}

		fmt.Printf(
			"message at topic/partition/offset \n%v/\n%v/\n%v: \n%s = \n%s\n",
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var c model.Choice

		if err := json.Unmarshal(m.Value, &c); err != nil {
			r.Logger.Error(
				"failed to unmarshal JSON to Choice consumer model",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		err = r.DB.StoreChoice(dbm.Choice{
			ID:                c.ID,
			Choice:            c.Choice,
			SuperheroID:       c.SuperheroID,
			ChosenSuperheroID: c.ChosenSuperheroID,
			CreatedAt:         c.CreatedAt,
		})
		if err != nil {
			r.Logger.Error(
				"failed to store choice to database",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}

		// If it is a like(1), then it should be saved to Cache.
		// Dislikes only go to the database.
		if c.Choice == like {
			ch := cache.Choice{
				ID:                c.ID,
				Choice:            c.Choice,
				SuperheroID:       c.SuperheroID,
				ChosenSuperheroID: c.ChosenSuperheroID,
				CreatedAt:         c.CreatedAt,
			}

			err = r.Cache.SetChoice(fmt.Sprintf(r.ChoiceKeyFormat, ch.SuperheroID, ch.ChosenSuperheroID), ch)
			if err != nil {
				r.Logger.Error(
					"failed to store choice in cache",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				err = r.Consumer.Close()
				if err != nil {
					r.Logger.Error(
						"failed to close consumer",
						zap.String("err", err.Error()),
						zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
					)

					return err
				}

				return err
			}

			err = r.Cache.SetLike(fmt.Sprintf(r.LikesKeyFormat, ch.ChosenSuperheroID), ch)
			if err != nil {
				r.Logger.Error(
					"failed to store like in cache",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				err = r.Consumer.Close()
				if err != nil {
					r.Logger.Error(
						"failed to close consumer",
						zap.String("err", err.Error()),
						zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
					)

					return err
				}

				return err
			}
		}

		err = r.Consumer.CommitMessages(ctx, m)
		if err != nil {
			r.Logger.Error(
				"failed to commit message",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}
	}
}
