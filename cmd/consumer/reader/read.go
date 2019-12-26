package reader

import (
	"context"
	"encoding/json"
	"fmt"

	cache "github.com/consumer-superhero-choice/internal/cache/model"
	"github.com/consumer-superhero-choice/internal/consumer/model"
	dbm "github.com/consumer-superhero-choice/internal/db/model"
)

// Read consumes the Kafka topic and stores the choice made by superhero
// to DB and if it is a like to Cache as well.
func (r *Reader) Read() error {
	ctx := context.Background()

	for {
		fmt.Print("before FetchMessage")
		m, err := r.Consumer.Consumer.FetchMessage(ctx)
		fmt.Print("after FetchMessage")
		if err != nil {
			err = r.Consumer.Consumer.Close()
			if err != nil {
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
			_ = r.Consumer.Consumer.Close()
			if err != nil {
				fmt.Println("Unmarshal")
				fmt.Println(err)
				err = r.Consumer.Consumer.Close()
				if err != nil {
					return err
				}

				return err
			}
		}

		err = r.DB.StoreChoice(dbm.Choice{
			ID:                c.ID,
			Choice:            c.Choice,
			SuperheroID:       c.SuperheroID,
			ChosenSuperheroID: c.ChosenSuperheroID,
			CreatedAt:         c.CreatedAt,
		}, )
		if err != nil {
			fmt.Println("DB")
			fmt.Println(err)
			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		// If it si a like(1), then it should be saved to Cache.
		if c.Choice == int64(1) {
			err = r.Cache.SetChoice(cache.Choice{
				ID:                c.ID,
				Choice:            c.Choice,
				SuperheroID:       c.SuperheroID,
				ChosenSuperheroID: c.ChosenSuperheroID,
				CreatedAt:         c.CreatedAt,
			}, )
			if err != nil {
				fmt.Println("Cache")
				fmt.Println(err)
				err = r.Consumer.Consumer.Close()
				if err != nil {
					return err
				}

				return err
			}
		}

		err = r.Consumer.Consumer.CommitMessages(ctx, m)
		if err != nil {
			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}
	}
}
