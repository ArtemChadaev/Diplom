package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type ProductResponse struct {
	ID                string     `json:"id"`
	SKU               string     `json:"sku"`
	Name              string     `json:"name"`
	GenericName       string     `json:"generic_name"`
	ATCCode           string     `json:"atc_code"`
	DosageForm        string     `json:"dosage_form"`
	Strength          string     `json:"strength"`
	PackageSize       int        `json:"package_size"`
	IsJNVLP           bool       `json:"is_jnvlp"`
	ManufacturerID    *string    `json:"manufacturer_id,omitempty"`
	StorageConditions string     `json:"storage_conditions"`
	PhotoURL          string     `json:"photo_url"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func ToProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:                p.ID,
		SKU:               p.SKU,
		Name:              p.Name,
		GenericName:       p.GenericName,
		ATCCode:           p.ATCCode,
		DosageForm:        p.DosageForm,
		Strength:          p.Strength,
		PackageSize:       p.PackageSize,
		IsJNVLP:           p.IsJNVLP,
		ManufacturerID:    p.ManufacturerID,
		StorageConditions: p.StorageConditions,
		PhotoURL:          p.PhotoURL,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		DeletedAt:         p.DeletedAt,
	}
}

type CreateProductRequest struct {
	SKU               string  `json:"sku" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	GenericName       string  `json:"generic_name" validate:"required"`
	ATCCode           string  `json:"atc_code" validate:"required"`
	DosageForm        string  `json:"dosage_form" validate:"required"`
	Strength          string  `json:"strength" validate:"required"`
	PackageSize       int     `json:"package_size" validate:"required,min=1"`
	IsJNVLP           bool    `json:"is_jnvlp"`
	ManufacturerID    *string `json:"manufacturer_id" validate:"omitempty,uuid"`
	StorageConditions string  `json:"storage_conditions"`
	PhotoURL          string  `json:"photo_url"`
}

func (r CreateProductRequest) ToDomain() domain.Product {
	return domain.Product{
		SKU:               r.SKU,
		Name:              r.Name,
		GenericName:       r.GenericName,
		ATCCode:           r.ATCCode,
		DosageForm:        r.DosageForm,
		Strength:          r.Strength,
		PackageSize:       r.PackageSize,
		IsJNVLP:           r.IsJNVLP,
		ManufacturerID:    r.ManufacturerID,
		StorageConditions: r.StorageConditions,
		PhotoURL:          r.PhotoURL,
	}
}

type UpdateProductRequest struct {
	Name              *string `json:"name"`
	GenericName       *string `json:"generic_name"`
	ATCCode           *string `json:"atc_code"`
	DosageForm        *string `json:"dosage_form"`
	Strength          *string `json:"strength"`
	PackageSize       *int    `json:"package_size"`
	IsJNVLP           *bool   `json:"is_jnvlp"`
	ManufacturerID    *string `json:"manufacturer_id"`
	StorageConditions *string `json:"storage_conditions"`
	PhotoURL          *string `json:"photo_url"`
}

func (r UpdateProductRequest) ApplyTo(p *domain.Product) {
	if r.Name != nil { p.Name = *r.Name }
	if r.GenericName != nil { p.GenericName = *r.GenericName }
	if r.ATCCode != nil { p.ATCCode = *r.ATCCode }
	if r.DosageForm != nil { p.DosageForm = *r.DosageForm }
	if r.Strength != nil { p.Strength = *r.Strength }
	if r.PackageSize != nil { p.PackageSize = *r.PackageSize }
	if r.IsJNVLP != nil { p.IsJNVLP = *r.IsJNVLP }
	if r.ManufacturerID != nil { p.ManufacturerID = r.ManufacturerID }
	if r.StorageConditions != nil { p.StorageConditions = *r.StorageConditions }
	if r.PhotoURL != nil { p.PhotoURL = *r.PhotoURL }
}

type ProductListResponse struct {
	Total    int               `json:"total"`
	Products []ProductResponse `json:"products"`
}
