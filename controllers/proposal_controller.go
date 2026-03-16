package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// SearchProposals godoc
// @Summary List proposals
// @Description Get all proposals in the system
// @Tags proposals
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Proposal
// @Failure 500 {object} map[string]string
// @Router /proposals/search [get]
func SearchProposals(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "25")
	limit := 25
	if n, err := strconv.Atoi(limitStr); err == nil {
		limit = n
	}

	results, err := services.SearchProposals(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
// GetAllProposals godoc
// @Summary List proposals
// @Description Get all proposals in the system
// @Tags proposals
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Proposal
// @Failure 500 {object} map[string]string
// @Router /proposals [get]
func GetAllProposals(c *gin.Context) {
	proposals, err := services.GetAllProposals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, proposals)
}
// GetProposalByID godoc
// @Summary Get proposal by ID
// @Description Get proposal by ID
// @Tags proposals
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Proposal
// @Failure 500 {object} map[string]string
// @Router /proposals/{id} [get]
func GetProposalByID(c *gin.Context) {
	id := c.Param("id")

	proposal, err := services.GetProposalByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	c.JSON(http.StatusOK, proposal)
}
// CreateProposal godoc
// @Summary Create a proposal
// @Description Create a new proposal
// @Tags proposals
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Proposal
// @Failure 500 {object} map[string]string
// @Router /proposals [post]
func CreateProposal(c *gin.Context) {
	var data models.ProposalCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Reference == "" || data.ClientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields reference and client_id are required"})
		return
	}

	proposal, err := services.CreateProposal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, proposal)
}
// UpdateProposal godoc
// @Summary Update a proposal
// @Description Update a proposal
// @Tags proposals
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Proposal
// @Failure 500 {object} map[string]string
// @Router /proposals/{id} [put]
func UpdateProposal(c *gin.Context) {
	id := c.Param("id")

	var data models.ProposalUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Reference == "" || data.ClientID == 0 || data.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields reference, client_id, and status are required"})
		return
	}

	proposal, err := services.UpdateProposal(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	c.JSON(http.StatusOK, proposal)
}

// DeleteProposal godoc
// @Summary Delete a proposal
// @Description Remove a proposal from the system
// @Tags proposals
// @Accept  json
// @Produce  json
// @Param id path int true "Proposal ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /proposals/{id} [delete]
func DeleteProposal(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteProposal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposal deleted successfully"})
}

// ── Proposal Items ────────────────────────────

// GetProposalItems godoc
// @Summary List proposal items
// @Description Get all items associated with a specific proposal
// @Tags proposals
// @Accept  json
// @Produce  json
// @Param proposal_id query int true "Proposal ID"
// @Success 200 {array} models.ProposalItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposal-items [get]
func GetProposalItems(c *gin.Context) {
	proposalID := c.Query("proposal_id")
	if proposalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "proposal_id query param is required"})
		return
	}

	items, err := services.GetProposalItems(proposalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateProposalItem godoc
// @Summary Create a proposal item
// @Description Add a new item to a proposal
// @Tags proposals
// @Accept  json
// @Produce  json
// @Param proposal_id query int true "Proposal ID"
// @Param item body models.ProposalItemCreate true "Proposal item details"
// @Success 201 {object} models.ProposalItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposal-items [post]
func CreateProposalItem(c *gin.Context) {
	proposalID := c.Query("proposal_id")
	if proposalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "proposal_id query param is required"})
		return
	}

	var data models.ProposalItemCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.ProductID == 0 || data.Quantity == 0 || data.UnitPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields product_id, quantity, and unit_price are required"})
		return
	}

	item, err := services.CreateProposalItem(proposalID, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateProposalItem godoc
// @Summary Update a proposal item
// @Description Update details of an existing proposal item
// @Tags proposals
// @Accept  json
// @Produce  json
// @Param id path int true "Proposal Item ID"
// @Param item body models.ProposalItemUpdate true "Updated proposal item details"
// @Success 200 {object} models.ProposalItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /proposal-items/{id} [put]
func UpdateProposalItem(c *gin.Context) {
	id := c.Param("id")

	var data models.ProposalItemUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.ProductID == 0 || data.Quantity == 0 || data.UnitPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields product_id, quantity, and unit_price are required"})
		return
	}

	item, err := services.UpdateProposalItem(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteProposalItem godoc
// @Summary Delete a proposal item
// @Description Remove an item from a proposal
// @Tags proposals
// @Accept  json
// @Produce  json
// @Param id path int true "Proposal Item ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /proposal-items/{id} [delete]
func DeleteProposalItem(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteProposalItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposal item deleted successfully"})
}
