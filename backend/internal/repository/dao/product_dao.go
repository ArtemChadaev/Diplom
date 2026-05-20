package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"gorm.io/gorm"
)

type ProductDAO struct {
	ID                string             `gorm:"column:id;primaryKey"`
	SKU               string             `gorm:"column:sku;uniqueIndex"`
	Name              string             `gorm:"column:name"`
	GenericName       string             `gorm:"column:generic_name"`
	ATCCode           string             `gorm:"column:atc_code"`
	DosageForm        string             `gorm:"column:dosage_form"`
	Strength          string             `gorm:"column:strength"`
	PackageSize       int                `gorm:"column:package_size"`
	IsJNVLP           bool               `gorm:"column:is_jnvlp"`
	ManufacturerID    *string            `gorm:"column:manufacturer_id"`
	StorageConditions string             `gorm:"column:storage_conditions"`
	PhotoURL          string             `gorm:"column:photo_url"`
	VenCategory       domain.VenCategory `gorm:"column:ven_category"`
	LeadTimeDays      int                `gorm:"column:lead_time_days"`
	SafetyStockQty    int                `gorm:"column:safety_stock_qty"`
	MaxStockQty       int                `gorm:"column:max_stock_qty"`
	CreatedAt         time.Time          `gorm:"column:created_at"`
	UpdatedAt         time.Time          `gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"column:deleted_at;index"`
}

func (ProductDAO) TableName() string {
	return "products"
}

func (p ProductDAO) ToDomain() domain.Product {
	var deletedAt *time.Time
	if p.DeletedAt.Valid {
		deletedAt = &p.DeletedAt.Time
	}

	return domain.Product{
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
		VenCategory:       p.VenCategory,
		LeadTimeDays:      p.LeadTimeDays,
		SafetyStockQty:    p.SafetyStockQty,
		MaxStockQty:       p.MaxStockQty,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}

func FromProductDomain(p domain.Product) ProductDAO {
	dao := ProductDAO{
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
		VenCategory:       p.VenCategory,
		LeadTimeDays:      p.LeadTimeDays,
		SafetyStockQty:    p.SafetyStockQty,
		MaxStockQty:       p.MaxStockQty,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
	if p.DeletedAt != nil {
		dao.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	}
	return dao
}
