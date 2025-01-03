package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/divizn/echo-calculator/internal/models"
)

var (
	REDIS_CACHE_TIMEOUT = time.Minute * 60
	REDIS_CACHE_RETRIES = 3
)

func createUserCachePrefix(id int, role string) string {
	return fmt.Sprintf("%v:%v:calculations", role, id)
}

// caches calculations for given user id
func (db *Database) CacheUserCalculations(id int, role string) error {
	key := createUserCachePrefix(id, role)
	calc, err := db.GetAllCalculations()
	if err != nil {
		return err
	}

	// marshal into json so can use as value otherwise error
	parsedCalcs, err := json.Marshal(calc)
	if err != nil {
		return err
	}

	return db.Cache.Set(*db.Ctx, key, string(parsedCalcs), REDIS_CACHE_TIMEOUT).Err()
}

// gets calculations from the cache for a given user
func (db *Database) GetCalculationsFromCache(id int, role string) (*[]models.Calculation, error) {
	key := createUserCachePrefix(id, role)

	retry := func() (*[]models.Calculation, error) {
		var lastErr error
		for range REDIS_CACHE_RETRIES {
			cached, err := db.Cache.Get(*db.Ctx, key).Result()
			if err != nil {
				if err.Error() == "redis: nil" {
					if cacheErr := db.CacheUserCalculations(id, role); cacheErr != nil {
						lastErr = cacheErr
						continue
					}
					continue
				}
				lastErr = err
				continue
			}
			parsedCalcs, err := parseCalculations(cached)
			if err != nil {
				return nil, err
			}
			return parsedCalcs, nil
		}
		return nil, lastErr
	}

	calcs, err := retry()
	if err != nil {
		return nil, err
	}

	return calcs, nil
}
func (db *Database) GetCalculationsFromCacheUser(id int, role string) (*[]models.Calculation, error) {
	for range REDIS_CACHE_RETRIES {
		calcs, err := db.GetCalculationsFromCache(id, role)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if calcs != nil {
			return calcs, nil
		}
	}
	return nil, fmt.Errorf("could not retrieve from cache")
}

func parseCalculations(calcs string) (*[]models.Calculation, error) {
	var parsedCalcs *[]models.Calculation
	err := json.Unmarshal([]byte(calcs), &parsedCalcs)
	if err != nil {
		return nil, err
	}

	return parsedCalcs, nil
}
