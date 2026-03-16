package routes

import (
	"net/http"
	"time"

	"aura-erp/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.GET("/verify", controllers.Verify)
			auth.POST("/logout", controllers.Logout)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", controllers.GetAllUsers)
			users.GET("/:id", controllers.GetUserByID)
			users.POST("", controllers.CreateUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}

		// Client routes
		clients := api.Group("/clients")
		{
			clients.GET("", controllers.GetAllClients)
			clients.GET("/search", controllers.SearchClients)
			clients.GET("/:id", controllers.GetClientByID)
			clients.POST("", controllers.CreateClient)
			clients.PUT("/:id", controllers.UpdateClient)
			clients.DELETE("/:id", controllers.DeleteClient)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", controllers.GetAllProducts)
			products.GET("/search", controllers.SearchProducts)
			products.GET("/:id", controllers.GetProductByID)
			products.POST("", controllers.CreateProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

		// Section routes
		sections := api.Group("/sections")
		{
			sections.GET("", controllers.GetAllSections)
			sections.GET("/search", controllers.SearchSections)
			sections.GET("/:id", controllers.GetSectionByID)
			sections.POST("", controllers.CreateSection)
			sections.PUT("/:id", controllers.UpdateSection)
			sections.DELETE("/:id", controllers.DeleteSection)
		}

		// Proposal routes — no item sub-routes here to avoid wildcard conflicts
		proposals := api.Group("/proposals")
		{
			proposals.GET("", controllers.GetAllProposals)
			proposals.GET("/search", controllers.SearchProposals)
			proposals.GET("/:id", controllers.GetProposalByID)
			proposals.POST("", controllers.CreateProposal)
			proposals.PUT("/:id", controllers.UpdateProposal)
			proposals.DELETE("/:id", controllers.DeleteProposal)
		}

		// Proposal item routes — top-level group, no wildcard conflict
		// GET/POST  /api/proposal-items?proposal_id=X  (list + create)
		// PUT/DELETE /api/proposal-items/:id            (update + delete)
		proposalItems := api.Group("/proposal-items")
		{
			proposalItems.GET("", controllers.GetProposalItems)
			proposalItems.POST("", controllers.CreateProposalItem)
			proposalItems.PUT("/:id", controllers.UpdateProposalItem)
			proposalItems.DELETE("/:id", controllers.DeleteProposalItem)
		}

		// Order routes — no item sub-routes here to avoid wildcard conflicts
		orders := api.Group("/orders")
		{
			orders.GET("", controllers.GetAllOrders)
			orders.GET("/:id", controllers.GetOrderByID)
			orders.POST("", controllers.CreateOrder)
			orders.PUT("/:id", controllers.UpdateOrder)
			orders.DELETE("/:id", controllers.DeleteOrder)
		}

		// Order item routes — top-level group, no wildcard conflict
		// GET/POST  /api/order-items?order_id=X  (list + create)
		// PUT/DELETE /api/order-items/:id         (update + delete)
		orderItems := api.Group("/order-items")
		{
			orderItems.GET("", controllers.GetOrderItems)
			orderItems.POST("", controllers.CreateOrderItem)
			orderItems.PUT("/:id", controllers.UpdateOrderItem)
			orderItems.DELETE("/:id", controllers.DeleteOrderItem)
		}

		// Statistics route
		api.GET("/statistics", controllers.GetStatistics)

		// Audit log routes
		auditLog := api.Group("/audit-log")
		{
			auditLog.GET("", controllers.GetAllAuditLogs)
			auditLog.GET("/:entityType/:entityId", controllers.GetAuditLogsByEntity)
			auditLog.POST("", controllers.CreateAuditLog)
		}
	}
}
