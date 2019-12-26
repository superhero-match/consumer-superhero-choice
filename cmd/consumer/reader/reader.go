package reader

import (
	"github.com/consumer-superhero-choice/internal/cache"
	"github.com/consumer-superhero-choice/internal/config"
	"github.com/consumer-superhero-choice/internal/consumer"
	"github.com/consumer-superhero-choice/internal/db"
)

// Reader holds all the data relevant.
type Reader struct {
	DB       *db.DB
	Cache    *cache.Cache
	Consumer *consumer.Consumer
}

// NewReader configures Reader.
func NewReader(cfg *config.Config) (r *Reader, err error) {
	dbs, err := db.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	ch, err := cache.NewCache(cfg)
	if err != nil {
		return nil, err
	}

	cs := consumer.NewConsumer(cfg)

	return &Reader{
		DB:       dbs,
		Cache:    ch,
		Consumer: cs,
	}, nil
}
