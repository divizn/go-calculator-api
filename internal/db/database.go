package db

import (
	"sync"

	"github.com/divizn/echo-calculator/internal/models"
)

type DbInstance struct {
	seq  *int
	data *map[int]*models.Calculation
	lock *sync.Mutex
}
