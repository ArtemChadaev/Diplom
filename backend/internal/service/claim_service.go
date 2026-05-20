package service

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
)

type claimService struct {
	repo      domain.ClaimRepository
	batchRepo domain.BatchRepository
}

func NewClaimService(repo domain.ClaimRepository, batchRepo domain.BatchRepository) domain.ClaimService {
	return &claimService{
		repo:      repo,
		batchRepo: batchRepo,
	}
}

func (s *claimService) ListClaims(ctx context.Context, limit, offset int) ([]domain.Claim, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *claimService) GetClaim(ctx context.Context, id string) (*domain.Claim, error) {
	claim, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, domain.ErrClaimNotFound
	}
	return claim, nil
}

func (s *claimService) CreateClaim(ctx context.Context, c *domain.Claim) (*domain.Claim, error) {
	c.Status = domain.ClaimStatusNew
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	if c.Type == "recall" && c.ProductID != nil && *c.ProductID != "" {
		if err := s.batchRepo.BlockAllByProductID(ctx, *c.ProductID); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (s *claimService) UpdateStatus(ctx context.Context, callerRole domain.UserRole, id string, status domain.ClaimStatus) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrClaimNotFound
	}

	// Status flow logic
	// Only Admin or WarehouseManager can move to resolved/rejected
	if status == domain.ClaimStatusResolved || status == domain.ClaimStatusRejected {
		if callerRole != domain.RoleAdmin && callerRole != domain.RoleWarehouseManager {
			return domain.ErrInsufficientPerms
		}
	}

	existing.Status = status
	return s.repo.Update(ctx, existing)
}
