package domain

import (
	"context"
)

// Country — справочник стран (ISO alpha-3).
type Country struct {
	Code   string `json:"code"`    // e.g. "RUS"
	NameRu string `json:"name_ru"` // e.g. "РФ"
}

// ATCCode — справочник кодов АТХ (Anatomical Therapeutic Chemical Classification).
type ATCCode struct {
	Code   string `json:"code"`    // e.g. "J01CA04"
	NameRu string `json:"name_ru"` // e.g. "Амоксициллин"
}

// ReferenceRepository — интерфейс для работы со справочниками.
type ReferenceRepository interface {
	ListCountries(ctx context.Context) ([]Country, error)
	ListATCCodes(ctx context.Context, query string, limit int) ([]ATCCode, error)
	GetCountryByCode(ctx context.Context, code string) (*Country, error)
	GetATCByCode(ctx context.Context, code string) (*ATCCode, error)
}

// ReferenceService — бизнес-логика справочников.
type ReferenceService interface {
	GetCountries(ctx context.Context) ([]Country, error)
	SearchATC(ctx context.Context, query string) ([]ATCCode, error)
}
