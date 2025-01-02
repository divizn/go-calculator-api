package services

import (
	"fmt"

	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/utils"
)

type Service struct {
	Config *utils.IConfig
	Db     *db.Database
}

func NewService() *Service {
	cfg, _ := utils.NewConfig()
	db, err := db.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	return &Service{Config: cfg, Db: db}
}

// returns error if id not valid
func (s *Service) validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("id cannot be 0 or negative")
	}
	return nil
}
