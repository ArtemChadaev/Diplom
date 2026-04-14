package dto

import "github.com/ima/diplom-backend/internal/domain"

type CountryResponse struct {
	Code   string `json:"code" example:"RUS"`
	NameRu string `json:"name_ru" example:"РФ"`
}

func FromCountryDomain(c domain.Country) CountryResponse {
	return CountryResponse{
		Code:   c.Code,
		NameRu: c.NameRu,
	}
}

type ATCCodeResponse struct {
	Code   string `json:"code" example:"J01CA04"`
	NameRu string `json:"name_ru" example:"Амоксициллин"`
}

func FromATCCodeDomain(a domain.ATCCode) ATCCodeResponse {
	return ATCCodeResponse{
		Code:   a.Code,
		NameRu: a.NameRu,
	}
}
