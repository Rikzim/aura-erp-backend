package models

// StatisticsResponse is the single payload returned by GET /api/statistics.
// All heavy aggregation is done in SQL — the frontend makes exactly one request.
type StatisticsResponse struct {
	// Counts
	TotalOrders    int `json:"total_orders"`
	TotalProposals int `json:"total_proposals"`
	TotalClients   int `json:"total_clients"`
	TotalProducts  int `json:"total_products"`
	TotalSections  int `json:"total_sections"`

	// Active/completed
	ActiveOrders    int `json:"active_orders"`    // pending + in_production
	CompletedOrders int `json:"completed_orders"` // completed + delivered

	// Revenue
	TotalOrderRevenue  float64 `json:"total_order_revenue"`
	TotalProposalValue float64 `json:"total_proposal_value"`
	AvgOrderValue      float64 `json:"avg_order_value"`

	// Conversion (approved proposals / total proposals * 100)
	ConversionRate float64 `json:"conversion_rate"`

	// Status breakdowns
	OrdersByStatus    []StatusCount `json:"orders_by_status"`
	ProposalsByStatus []StatusCount `json:"proposals_by_status"`

	// Section breakdowns
	OrdersBySection    []LabelCount `json:"orders_by_section"`
	ProposalsBySection []LabelCount `json:"proposals_by_section"`

	// Time series — last 6 calendar months
	OrdersPerMonth    []MonthCount `json:"orders_per_month"`
	ProposalsPerMonth []MonthCount `json:"proposals_per_month"`
	ClientsPerMonth   []MonthCount `json:"clients_per_month"`

	// Top-N lists
	TopClientsByRevenue []ClientRevenue `json:"top_clients_by_revenue"`
	TopProductsByQty    []ProductQty    `json:"top_products_by_qty"`
}

// StatusCount holds a status string and how many records have that status.
type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// LabelCount holds a human-readable label and a count (used for sections).
type LabelCount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// MonthCount represents a single month bucket in a time series.
// Year and Month are kept separate so the frontend can re-format as needed.
type MonthCount struct {
	Year  int `json:"year"`
	Month int `json:"month"` // 1–12
	Count int `json:"count"`
}

// ClientRevenue pairs a client name with the total revenue generated from their orders.
type ClientRevenue struct {
	ClientID   int     `json:"client_id"`
	ClientName string  `json:"client_name"`
	Revenue    float64 `json:"revenue"`
}

// ProductQty pairs a product name with the total quantity ordered across all orders.
type ProductQty struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalQty    float64 `json:"total_qty"`
}
