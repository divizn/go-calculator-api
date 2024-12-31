package services

import "github.com/divizn/echo-calculator/internal/utils"

type Service struct {
	Config *utils.IConfig
}

func NewService() *Service {
	cfg, _ := utils.NewConfig()
	return &Service{Config: cfg}
}
