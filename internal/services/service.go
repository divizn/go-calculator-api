package services

import (
	"fmt"

	"github.com/divizn/echo-calculator/internal/utils"
)

type Service struct {
	Config *utils.IConfig
}

func NewService() *Service {
	cfg, _ := utils.NewConfig()
	return &Service{Config: cfg}
}

// returns error if id not valid
func (s *Service) validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("id cannot be 0 or negative")
	}
	return nil
}
