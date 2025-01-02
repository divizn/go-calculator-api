package db

import (
	"fmt"

	"github.com/divizn/echo-calculator/internal/models"
)

func createUserCachePrefix(user_id int) string {
	return fmt.Sprintf("user:%v:calculations", user_id)
}

// caches calculations for given user
func (db *Database) CacheUserCalc(user_id int, calc *[]models.Calculation) error {
	key := createUserCachePrefix(user_id)

	err := db.Cache.Set(*db.Ctx, key, calc, REDIS_CACHE_TIMEOUT).Err()
	if err != nil {
		return err
	}
	return nil
}
