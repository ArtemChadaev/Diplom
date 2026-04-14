package service

import (
	"context"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type inventoryService struct {
	repo domain.InventoryRepository
}

func NewInventoryService(repo domain.InventoryRepository) domain.InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) ListSessions(ctx context.Context, limit, offset int) ([]domain.InventorySession, int, error) {
	return s.repo.ListSessions(ctx, limit, offset)
}

func (s *inventoryService) GetSession(ctx context.Context, id string) (*domain.InventorySession, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrInventorySessionNotFound
	}
	return session, nil
}

func (s *inventoryService) StartSession(ctx context.Context, userID int, zoneID string) (*domain.InventorySession, error) {
	session := &domain.InventorySession{
		ZoneID:    zoneID,
		Status:    domain.InventoryStatusActive,
		StartedBy: userID,
		StartedAt: time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *inventoryService) FinishSession(ctx context.Context, id string) error {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return err
	}
	if session == nil {
		return domain.ErrInventorySessionNotFound
	}

	now := time.Now()
	session.Status = domain.InventoryStatusCompleted
	session.CompletedAt = &now

	return s.repo.UpdateSession(ctx, session)
}

func (s *inventoryService) SubmitCount(ctx context.Context, sessionID string, items []domain.InventoryItem) error {
	for i := range items {
		items[i].SessionID = sessionID
		if err := s.repo.AddCount(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}
