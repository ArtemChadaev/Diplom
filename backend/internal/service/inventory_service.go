package service

import (
	"context"
	"sort"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type inventoryService struct {
	repo        domain.InventoryRepository
	profileRepo domain.EmployeeProfileRepository
	productRepo domain.ProductRepository
}

func NewInventoryService(
	repo domain.InventoryRepository,
	profileRepo domain.EmployeeProfileRepository,
	productRepo domain.ProductRepository,
) domain.InventoryService {
	return &inventoryService{
		repo:        repo,
		profileRepo: profileRepo,
		productRepo: productRepo,
	}
}

func (s *inventoryService) ListSessions(ctx context.Context, limit, offset int) ([]domain.InventorySession, int, error) {
	sessions, total, err := s.repo.ListSessions(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	for i := range sessions {
		if sessions[i].Status == domain.InventoryStatusActive {
			for j := range sessions[i].Items {
				sessions[i].Items[j].SystemQuantity = 0
			}
		}
	}
	return sessions, total, nil
}

func (s *inventoryService) GetSession(ctx context.Context, id string) (*domain.InventorySession, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrInventorySessionNotFound
	}

	if session.Status == domain.InventoryStatusActive {
		for i := range session.Items {
			session.Items[i].SystemQuantity = 0
		}
	}

	return session, nil
}

func (s *inventoryService) StartSession(ctx context.Context, userID int, zoneID string) (*domain.InventorySession, error) {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, domain.ErrGDPTrainingRequired
	}
	if err := CheckGDPValid(profile); err != nil {
		return nil, err
	}

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

func (s *inventoryService) CalculateNetting(ctx context.Context, sessionID string) ([]domain.NettingLine, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.ErrInventorySessionNotFound
	}

	if session.Status != domain.InventoryStatusCompleted {
		return nil, domain.ErrInventoryNotCompleted
	}

	type prodQty struct {
		system   int
		physical int
	}
	agg := make(map[string]*prodQty)
	for _, item := range session.Items {
		q, exists := agg[item.ProductID]
		if !exists {
			q = &prodQty{}
			agg[item.ProductID] = q
		}
		q.system += item.SystemQuantity
		q.physical += item.PhysicalQuantity
	}

	var lines []domain.NettingLine
	for prodID, q := range agg {
		delta := q.physical - q.system
		if delta == 0 {
			continue
		}

		product, err := s.productRepo.GetByID(ctx, prodID)
		if err != nil {
			continue
		}
		if product == nil {
			continue
		}

		atcGroup := ""
		if len(product.ATCCode) >= 3 {
			atcGroup = product.ATCCode[:3]
		} else {
			atcGroup = product.ATCCode
		}

		lines = append(lines, domain.NettingLine{
			ProductID:   product.ID,
			ProductName: product.Name,
			ATCGroup:    atcGroup,
			Delta:       delta,
		})
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].ATCGroup < lines[j].ATCGroup
	})

	return lines, nil
}

