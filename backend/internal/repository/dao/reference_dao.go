package dao

import "github.com/ima/diplom-backend/internal/domain"

// CountryDAO maps to the 'countries' table.
type CountryDAO struct {
	Code   string `gorm:"column:code;primaryKey"`
	NameRu string `gorm:"column:name_ru"`
}

func (CountryDAO) TableName() string {
	return "countries"
}

func (c CountryDAO) ToDomain() domain.Country {
	return domain.Country{
		Code:   c.Code,
		NameRu: c.NameRu,
	}
}

// ATCCodeDAO maps to the 'atc_codes' table.
type ATCCodeDAO struct {
	Code   string `gorm:"column:code;primaryKey"`
	NameRu string `gorm:"column:name_ru"`
}

func (ATCCodeDAO) TableName() string {
	return "atc_codes"
}

func (a ATCCodeDAO) ToDomain() domain.ATCCode {
	return domain.ATCCode{
		Code:   a.Code,
		NameRu: a.NameRu,
	}
}
