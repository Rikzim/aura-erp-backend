package services

import (
	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func GetStatistics() (*models.StatisticsResponse, error) {
	stats := &models.StatisticsResponse{}

	// ── Simple counts ────────────────────────────────────────────────────────
	err := config.DB.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM orders)    AS total_orders,
			(SELECT COUNT(*) FROM proposals) AS total_proposals,
			(SELECT COUNT(*) FROM clients)   AS total_clients,
			(SELECT COUNT(*) FROM products)  AS total_products,
			(SELECT COUNT(*) FROM sections)  AS total_sections,

			(SELECT COUNT(*) FROM orders WHERE status IN ('pending','in_production'))      AS active_orders,
			(SELECT COUNT(*) FROM orders WHERE status IN ('completed','delivered'))        AS completed_orders,

			COALESCE((SELECT SUM(oi.quantity * oi.unit_price)
			          FROM order_items oi), 0)                                             AS total_order_revenue,

			COALESCE((SELECT SUM(pi.quantity * pi.unit_price)
			          FROM proposal_items pi), 0)                                          AS total_proposal_value,

			COALESCE((SELECT SUM(oi.quantity * oi.unit_price) / NULLIF(COUNT(DISTINCT oi.order_id), 0)
			          FROM order_items oi), 0)                                             AS avg_order_value,

			COALESCE(
				(SELECT COUNT(*) FILTER (WHERE status = 'approved')::float
				 FROM proposals)
				/ NULLIF((SELECT COUNT(*) FROM proposals)::float, 0) * 100
			, 0)                                                                           AS conversion_rate
	`).Scan(
		&stats.TotalOrders,
		&stats.TotalProposals,
		&stats.TotalClients,
		&stats.TotalProducts,
		&stats.TotalSections,
		&stats.ActiveOrders,
		&stats.CompletedOrders,
		&stats.TotalOrderRevenue,
		&stats.TotalProposalValue,
		&stats.AvgOrderValue,
		&stats.ConversionRate,
	)
	if err != nil {
		return nil, err
	}

	// ── Orders by status ─────────────────────────────────────────────────────
	rows, err := config.DB.Query(`
		SELECT status, COUNT(*) AS count
		FROM orders
		GROUP BY status
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sc models.StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		stats.OrdersByStatus = append(stats.OrdersByStatus, sc)
	}
	if stats.OrdersByStatus == nil {
		stats.OrdersByStatus = []models.StatusCount{}
	}

	// ── Proposals by status ──────────────────────────────────────────────────
	rows2, err := config.DB.Query(`
		SELECT status, COUNT(*) AS count
		FROM proposals
		GROUP BY status
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var sc models.StatusCount
		if err := rows2.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		stats.ProposalsByStatus = append(stats.ProposalsByStatus, sc)
	}
	if stats.ProposalsByStatus == nil {
		stats.ProposalsByStatus = []models.StatusCount{}
	}

	// ── Orders by section ────────────────────────────────────────────────────
	rows3, err := config.DB.Query(`
		SELECT COALESCE(s.name, 'No Section') AS label, COUNT(o.id) AS count
		FROM orders o
		LEFT JOIN sections s ON o.section_id = s.id
		GROUP BY s.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows3.Close()
	for rows3.Next() {
		var lc models.LabelCount
		if err := rows3.Scan(&lc.Label, &lc.Count); err != nil {
			return nil, err
		}
		stats.OrdersBySection = append(stats.OrdersBySection, lc)
	}
	if stats.OrdersBySection == nil {
		stats.OrdersBySection = []models.LabelCount{}
	}

	// ── Proposals by section ─────────────────────────────────────────────────
	rows4, err := config.DB.Query(`
		SELECT COALESCE(s.name, 'No Section') AS label, COUNT(p.id) AS count
		FROM proposals p
		LEFT JOIN sections s ON p.section_id = s.id
		GROUP BY s.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows4.Close()
	for rows4.Next() {
		var lc models.LabelCount
		if err := rows4.Scan(&lc.Label, &lc.Count); err != nil {
			return nil, err
		}
		stats.ProposalsBySection = append(stats.ProposalsBySection, lc)
	}
	if stats.ProposalsBySection == nil {
		stats.ProposalsBySection = []models.LabelCount{}
	}

	// ── Orders per month (last 6 months) ─────────────────────────────────────
	rows5, err := config.DB.Query(`
		SELECT
			EXTRACT(YEAR  FROM created_at)::int AS year,
			EXTRACT(MONTH FROM created_at)::int AS month,
			COUNT(*)                            AS count
		FROM orders
		WHERE created_at >= date_trunc('month', NOW()) - INTERVAL '5 months'
		GROUP BY year, month
		ORDER BY year, month
	`)
	if err != nil {
		return nil, err
	}
	defer rows5.Close()
	for rows5.Next() {
		var mc models.MonthCount
		if err := rows5.Scan(&mc.Year, &mc.Month, &mc.Count); err != nil {
			return nil, err
		}
		stats.OrdersPerMonth = append(stats.OrdersPerMonth, mc)
	}
	if stats.OrdersPerMonth == nil {
		stats.OrdersPerMonth = []models.MonthCount{}
	}

	// ── Proposals per month (last 6 months) ──────────────────────────────────
	rows6, err := config.DB.Query(`
		SELECT
			EXTRACT(YEAR  FROM created_at)::int AS year,
			EXTRACT(MONTH FROM created_at)::int AS month,
			COUNT(*)                            AS count
		FROM proposals
		WHERE created_at >= date_trunc('month', NOW()) - INTERVAL '5 months'
		GROUP BY year, month
		ORDER BY year, month
	`)
	if err != nil {
		return nil, err
	}
	defer rows6.Close()
	for rows6.Next() {
		var mc models.MonthCount
		if err := rows6.Scan(&mc.Year, &mc.Month, &mc.Count); err != nil {
			return nil, err
		}
		stats.ProposalsPerMonth = append(stats.ProposalsPerMonth, mc)
	}
	if stats.ProposalsPerMonth == nil {
		stats.ProposalsPerMonth = []models.MonthCount{}
	}

	// ── New clients per month (last 6 months) ────────────────────────────────
	rows7, err := config.DB.Query(`
		SELECT
			EXTRACT(YEAR  FROM created_at)::int AS year,
			EXTRACT(MONTH FROM created_at)::int AS month,
			COUNT(*)                            AS count
		FROM clients
		WHERE created_at >= date_trunc('month', NOW()) - INTERVAL '5 months'
		GROUP BY year, month
		ORDER BY year, month
	`)
	if err != nil {
		return nil, err
	}
	defer rows7.Close()
	for rows7.Next() {
		var mc models.MonthCount
		if err := rows7.Scan(&mc.Year, &mc.Month, &mc.Count); err != nil {
			return nil, err
		}
		stats.ClientsPerMonth = append(stats.ClientsPerMonth, mc)
	}
	if stats.ClientsPerMonth == nil {
		stats.ClientsPerMonth = []models.MonthCount{}
	}

	// ── Top 5 clients by order revenue ───────────────────────────────────────
	rows8, err := config.DB.Query(`
		SELECT
			c.id,
			c.name,
			COALESCE(SUM(oi.quantity * oi.unit_price), 0) AS revenue
		FROM clients c
		JOIN orders o ON o.client_id = c.id
		JOIN order_items oi ON oi.order_id = o.id
		GROUP BY c.id, c.name
		ORDER BY revenue DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows8.Close()
	for rows8.Next() {
		var cr models.ClientRevenue
		if err := rows8.Scan(&cr.ClientID, &cr.ClientName, &cr.Revenue); err != nil {
			return nil, err
		}
		stats.TopClientsByRevenue = append(stats.TopClientsByRevenue, cr)
	}
	if stats.TopClientsByRevenue == nil {
		stats.TopClientsByRevenue = []models.ClientRevenue{}
	}

	// ── Top 5 products by quantity ordered ───────────────────────────────────
	rows9, err := config.DB.Query(`
		SELECT
			p.id,
			p.name,
			COALESCE(SUM(oi.quantity), 0) AS total_qty
		FROM products p
		JOIN order_items oi ON oi.product_id = p.id
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows9.Close()
	for rows9.Next() {
		var pq models.ProductQty
		if err := rows9.Scan(&pq.ProductID, &pq.ProductName, &pq.TotalQty); err != nil {
			return nil, err
		}
		stats.TopProductsByQty = append(stats.TopProductsByQty, pq)
	}
	if stats.TopProductsByQty == nil {
		stats.TopProductsByQty = []models.ProductQty{}
	}

	return stats, nil
}
