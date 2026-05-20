package domain

import (
	"context"
	"time"
)

type VenCategory string

const (
	VenV VenCategory = "V"
	VenE VenCategory = "E"
	VenN VenCategory = "N"
)

// Product — медицинский препарат / товар.
type Product struct {
	ID                string      `json:"id"`
	SKU               string      `json:"sku"`               // уникальный складской код
	Name              string      `json:"name"`              // торговое наименование
	GenericName       string      `json:"generic_name"`      // МНН (международное непатентованное наименование)
	ATCCode           string      `json:"atc_code"`          // код АТХ
	DosageForm        string      `json:"dosage_form"`       // лекарственная форма (таблетки, капсулы и т.д.)
	Strength          string      `json:"strength"`          // дозировка (например, "500 мг")
	PackageSize       int         `json:"package_size"`      // кол-во в упаковке
	IsJNVLP           bool        `json:"is_jnvlp"`          // входит ли в список ЖНВЛП (Жизненно Необходимые и Важнейшие Лекарственные Препараты)
	ManufacturerID    *string     `json:"manufacturer_id"`   // ID производителя (опционально)
	StorageConditions string      `json:"storage_conditions"` // условия хранения
	PhotoURL          string      `json:"photo_url"`         // ссылка на фото (S3)
	VenCategory       VenCategory `json:"ven_category"`
	LeadTimeDays      int         `json:"lead_time_days"`
	SafetyStockQty    int         `json:"safety_stock_qty"`
	MaxStockQty       int         `json:"max_stock_qty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	DeletedAt         *time.Time  `json:"deleted_at,omitempty"`
}

// ProductFilter — фильтрация при поиске товаров.
type ProductFilter struct {
	Query     string
	IsJNVLP   *bool
	ATCCode   string
	Limit     int
	Offset    int
}

// ProductRepository — интерфейс для работы с товарами.
type ProductRepository interface {
	List(ctx context.Context, filter ProductFilter) ([]Product, int, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	GetBySKU(ctx context.Context, sku string) (*Product, error)
}

// ROPResult — результат расчета точки перезаказа (Reorder Point)
type ROPResult struct {
	ProductID       string  `json:"product_id"`
	SKU             string  `json:"sku"`
	Name            string  `json:"name"`
	CurrentStock    int     `json:"current_stock"`
	SafetyStock     int     `json:"safety_stock"`
	MaxStock        int     `json:"max_stock"`
	MonthlyTurnover int     `json:"monthly_turnover"`
	DailyUsage      float64 `json:"daily_usage"`
	ROP             int     `json:"rop"`
	NeedsReorder    bool    `json:"needs_reorder"`
	ReorderQty      int     `json:"reorder_qty"`
}

// ProductService — бизнес-логика товаров.
type ProductService interface {
	ListProducts(ctx context.Context, filter ProductFilter) ([]Product, int, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	CreateProduct(ctx context.Context, callerRole UserRole, p *Product) (*Product, error)
	UpdateProduct(ctx context.Context, callerRole UserRole, p *Product) (*Product, error)
	DeleteProduct(ctx context.Context, callerRole UserRole, id string) error
	CheckReorderPoint(ctx context.Context, productID string) (*ROPResult, error)
	RunReorderCheckAll(ctx context.Context) ([]ROPResult, error)
}
