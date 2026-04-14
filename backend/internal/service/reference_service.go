package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type referenceService struct {
	repo domain.ReferenceRepository
}

func NewReferenceService(repo domain.ReferenceRepository) domain.ReferenceService {
	return &referenceService{repo: repo}
}

func (s *referenceService) GetCountries(ctx context.Context) ([]domain.Country, error) {
	return s.repo.ListCountries(ctx)
}

func (s *referenceService) SearchATC(ctx context.Context, query string) ([]domain.ATCCode, error) {
	// Simple search with limit 100
	return s.repo.ListATCCodes(ctx, query, 100)
}
