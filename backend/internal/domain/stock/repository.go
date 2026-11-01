package stock

type Product struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenantId"`
	Name          string  `json:"name"`
	SKU           string  `json:"sku"`
	Category      string  `json:"category"`
	PriceCents    int     `json:"priceCents"`
	CostCents     int     `json:"costCents"`
	StockQuantity float64 `json:"stockQuantity"`
	MinStockAlert float64 `json:"minStockAlert"`
	IsSupply      bool    `json:"isSupply"`
	Unit          string  `json:"unit"`
	Active        bool    `json:"active"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type ServiceSupply struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenantId"`
	ServiceID        string  `json:"serviceId"`
	ProductID        string  `json:"productId"`
	ProductName      string  `json:"productName,omitempty"`
	QuantityConsumed float64 `json:"quantityConsumed"`
	CreatedAt        string  `json:"createdAt"`
}

type StockMovement struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenantId"`
	ProductID     string  `json:"productId"`
	ProductName   string  `json:"productName,omitempty"`
	Type          string  `json:"type"` // "in", "out", "adjustment", "consumption", "sale"
	Quantity      float64 `json:"quantity"`
	UnitCostCents int     `json:"unitCostCents"`
	Notes         string  `json:"notes"`
	AppointmentID string  `json:"appointmentId,omitempty"`
	CreatedAt     string  `json:"createdAt"`
}

type Repository interface {
	CreateProduct(input Product) (Product, error)
	UpdateProduct(input Product) (Product, error)
	GetProductByID(tenantID string, id string) (Product, error)
	ListProducts(tenantID string, activeOnly bool) ([]Product, error)
	DeleteProduct(tenantID string, id string) error

	CreateStockMovement(input StockMovement) (StockMovement, error)
	ListStockMovements(tenantID string, productID string, limit int) ([]StockMovement, error)

	SetServiceSupplies(tenantID string, serviceID string, supplies []ServiceSupply) error
	ListServiceSupplies(tenantID string, serviceID string) ([]ServiceSupply, error)
	DeductSuppliesForAppointment(tenantID string, appointmentID string, serviceID string) error
}

