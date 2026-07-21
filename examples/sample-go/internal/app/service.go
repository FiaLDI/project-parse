package app

import "example.com/sample-go/internal/domain"

// Service is a sample application service.
type Service struct{}

func (s *Service) Load(id string) domain.Entity {
	return domain.Entity{ID: id}
}
