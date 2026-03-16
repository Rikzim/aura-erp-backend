package services

import (
	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func GetAllOrders() ([]models.Order, error) {
	query := `
		SELECT o.id, o.reference, o.proposal_id, o.client_id, o.section_id, o.status, o.due_date, o.created_at, o.updated_at,
		       c.name, s.name, p.reference
		FROM orders o
		LEFT JOIN clients c ON o.client_id = c.id
		LEFT JOIN sections s ON o.section_id = s.id
		LEFT JOIN proposals p ON o.proposal_id = p.id
		ORDER BY o.id ASC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.Reference, &order.ProposalID, &order.ClientID, &order.SectionID, &order.Status, &order.DueDate, &order.CreatedAt, &order.UpdatedAt, &order.ClientName, &order.SectionName, &order.ProposalReference); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func GetOrderByID(id string) (*models.Order, error) {
	query := `
		SELECT o.id, o.reference, o.proposal_id, o.client_id, o.section_id, o.status, o.due_date, o.created_at, o.updated_at,
		       c.name, s.name, p.reference
		FROM orders o
		LEFT JOIN clients c ON o.client_id = c.id
		LEFT JOIN sections s ON o.section_id = s.id
		LEFT JOIN proposals p ON o.proposal_id = p.id
		WHERE o.id = $1
	`

	var order models.Order
	err := config.DB.QueryRow(query, id).Scan(
		&order.ID, &order.Reference, &order.ProposalID, &order.ClientID, &order.SectionID, &order.Status, &order.DueDate, &order.CreatedAt, &order.UpdatedAt, &order.ClientName, &order.SectionName, &order.ProposalReference,
	)
	if err != nil {
		return nil, err
	}

	// Get items
	items, err := GetOrderItems(id)
	if err == nil {
		order.Items = items
	}

	return &order, nil
}

func CreateOrder(data models.OrderCreate) (*models.Order, error) {
	query := `INSERT INTO orders (reference, proposal_id, client_id, section_id, status, due_date)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, reference, proposal_id, client_id, section_id, status, due_date, created_at, updated_at`

	var order models.Order
	status := data.Status
	if status == "" {
		status = "pending"
	}

	err := config.DB.QueryRow(query, data.Reference, data.ProposalID, data.ClientID, data.SectionID, status, data.DueDate).Scan(
		&order.ID, &order.Reference, &order.ProposalID, &order.ClientID, &order.SectionID, &order.Status, &order.DueDate, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrder(id string, data models.OrderUpdate) (*models.Order, error) {
	query := `UPDATE orders SET reference = $1, proposal_id = $2, client_id = $3, section_id = $4, status = $5, due_date = $6, updated_at = NOW()
	          WHERE id = $7
	          RETURNING id, reference, proposal_id, client_id, section_id, status, due_date, created_at, updated_at`

	var order models.Order
	err := config.DB.QueryRow(query, data.Reference, data.ProposalID, data.ClientID, data.SectionID, data.Status, data.DueDate, id).Scan(
		&order.ID, &order.Reference, &order.ProposalID, &order.ClientID, &order.SectionID, &order.Status, &order.DueDate, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func DeleteOrder(id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

// ── Order Items ───────────────────────────────

// GetAllOrderItems returns every order_item row across all orders in a single
// query. Used by the calendar page to bulk-hydrate cards without N+1 requests.
func GetAllOrderItems() ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.notes,
		       pr.name, pr.unit
		FROM order_items oi
		LEFT JOIN products pr ON oi.product_id = pr.id
		ORDER BY oi.order_id ASC, oi.id ASC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.ProductName, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.OrderItem{}
	}
	return items, nil
}

func GetOrderItems(orderID string) ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.notes,
		       pr.name, pr.unit
		FROM order_items oi
		LEFT JOIN products pr ON oi.product_id = pr.id
		WHERE oi.order_id = $1
		ORDER BY oi.id ASC
	`
	rows, err := config.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.ProductName, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func CreateOrderItem(orderID string, data models.OrderItemCreate) (*models.OrderItem, error) {
	query := `INSERT INTO order_items (order_id, product_id, quantity, unit_price, notes)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, order_id, product_id, quantity, unit_price, notes`

	var item models.OrderItem
	err := config.DB.QueryRow(query, orderID, data.ProductID, data.Quantity, data.UnitPrice, data.Notes).Scan(
		&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateOrderItem(itemID string, data models.OrderItemUpdate) (*models.OrderItem, error) {
	query := `UPDATE order_items SET product_id = $1, quantity = $2, unit_price = $3, notes = $4
	          WHERE id = $5
	          RETURNING id, order_id, product_id, quantity, unit_price, notes`

	var item models.OrderItem
	err := config.DB.QueryRow(query, data.ProductID, data.Quantity, data.UnitPrice, data.Notes, itemID).Scan(
		&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func DeleteOrderItem(itemID string) error {
	query := `DELETE FROM order_items WHERE id = $1`
	result, err := config.DB.Exec(query, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}
