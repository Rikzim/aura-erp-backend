package services

import (
	"fmt"

	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func GetAllProposals() ([]models.Proposal, error) {
	query := `
		SELECT p.id, p.reference, p.client_id, p.section_id, p.status, p.notes, p.created_at, p.updated_at,
		       c.name, s.name
		FROM proposals p
		LEFT JOIN clients c ON p.client_id = c.id
		LEFT JOIN sections s ON p.section_id = s.id
		ORDER BY p.id ASC
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proposals []models.Proposal
	for rows.Next() {
		var proposal models.Proposal
		if err := rows.Scan(&proposal.ID, &proposal.Reference, &proposal.ClientID, &proposal.SectionID, &proposal.Status, &proposal.Notes, &proposal.CreatedAt, &proposal.UpdatedAt, &proposal.ClientName, &proposal.SectionName); err != nil {
			return nil, err
		}
		proposals = append(proposals, proposal)
	}

	return proposals, nil
}

func SearchProposals(q string, limit int) ([]models.ProposalSearchResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	query := `
		SELECT id, reference
		FROM proposals
		WHERE reference ILIKE $1
		ORDER BY id DESC
		LIMIT $2
	`
	rows, err := config.DB.Query(query, fmt.Sprintf("%%%s%%", q), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProposalSearchResult
	for rows.Next() {
		var r models.ProposalSearchResult
		if err := rows.Scan(&r.ID, &r.Reference); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if results == nil {
		results = []models.ProposalSearchResult{}
	}
	return results, nil
}

func GetProposalByID(id string) (*models.Proposal, error) {
	query := `
		SELECT p.id, p.reference, p.client_id, p.section_id, p.status, p.notes, p.created_at, p.updated_at,
		       c.name, s.name
		FROM proposals p
		LEFT JOIN clients c ON p.client_id = c.id
		LEFT JOIN sections s ON p.section_id = s.id
		WHERE p.id = $1
	`

	var proposal models.Proposal
	err := config.DB.QueryRow(query, id).Scan(
		&proposal.ID, &proposal.Reference, &proposal.ClientID, &proposal.SectionID, &proposal.Status, &proposal.Notes, &proposal.CreatedAt, &proposal.UpdatedAt, &proposal.ClientName, &proposal.SectionName,
	)
	if err != nil {
		return nil, err
	}

	// Get items
	items, err := GetProposalItems(id)
	if err == nil {
		proposal.Items = items
	}

	return &proposal, nil
}

func CreateProposal(data models.ProposalCreate) (*models.Proposal, error) {
	query := `INSERT INTO proposals (reference, client_id, section_id, status, notes)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, reference, client_id, section_id, status, notes, created_at, updated_at`

	var proposal models.Proposal
	status := data.Status
	if status == "" {
		status = "draft"
	}

	err := config.DB.QueryRow(query, data.Reference, data.ClientID, data.SectionID, status, data.Notes).Scan(
		&proposal.ID, &proposal.Reference, &proposal.ClientID, &proposal.SectionID, &proposal.Status, &proposal.Notes, &proposal.CreatedAt, &proposal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &proposal, nil
}

func UpdateProposal(id string, data models.ProposalUpdate) (*models.Proposal, error) {
	query := `UPDATE proposals SET reference = $1, client_id = $2, section_id = $3, status = $4, notes = $5, updated_at = NOW()
	          WHERE id = $6
	          RETURNING id, reference, client_id, section_id, status, notes, created_at, updated_at`

	var proposal models.Proposal
	err := config.DB.QueryRow(query, data.Reference, data.ClientID, data.SectionID, data.Status, data.Notes, id).Scan(
		&proposal.ID, &proposal.Reference, &proposal.ClientID, &proposal.SectionID, &proposal.Status, &proposal.Notes, &proposal.CreatedAt, &proposal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &proposal, nil
}

func DeleteProposal(id string) error {
	query := `DELETE FROM proposals WHERE id = $1`
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

// ── Proposal Items ────────────────────────────

func GetProposalItems(proposalID string) ([]models.ProposalItem, error) {
	query := `
		SELECT pi.id, pi.proposal_id, pi.product_id, pi.quantity, pi.unit_price, pi.notes,
		       pr.name, pr.unit
		FROM proposal_items pi
		LEFT JOIN products pr ON pi.product_id = pr.id
		WHERE pi.proposal_id = $1
		ORDER BY pi.id ASC
	`
	rows, err := config.DB.Query(query, proposalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ProposalItem
	for rows.Next() {
		var item models.ProposalItem
		if err := rows.Scan(&item.ID, &item.ProposalID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.ProductName, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func CreateProposalItem(proposalID string, data models.ProposalItemCreate) (*models.ProposalItem, error) {
	query := `INSERT INTO proposal_items (proposal_id, product_id, quantity, unit_price, notes)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, proposal_id, product_id, quantity, unit_price, notes`

	var item models.ProposalItem
	err := config.DB.QueryRow(query, proposalID, data.ProductID, data.Quantity, data.UnitPrice, data.Notes).Scan(
		&item.ID, &item.ProposalID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateProposalItem(itemID string, data models.ProposalItemUpdate) (*models.ProposalItem, error) {
	query := `UPDATE proposal_items SET product_id = $1, quantity = $2, unit_price = $3, notes = $4
	          WHERE id = $5
	          RETURNING id, proposal_id, product_id, quantity, unit_price, notes`

	var item models.ProposalItem
	err := config.DB.QueryRow(query, data.ProductID, data.Quantity, data.UnitPrice, data.Notes, itemID).Scan(
		&item.ID, &item.ProposalID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func DeleteProposalItem(itemID string) error {
	query := `DELETE FROM proposal_items WHERE id = $1`
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
