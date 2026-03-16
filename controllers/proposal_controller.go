package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

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

func GetAllProposals(c *gin.Context) {
	proposals, err := services.GetAllProposals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, proposals)
}

func GetProposalByID(c *gin.Context) {
	id := c.Param("id")

	proposal, err := services.GetProposalByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	c.JSON(http.StatusOK, proposal)
}

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

func DeleteProposalItem(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteProposalItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposal item deleted successfully"})
}
