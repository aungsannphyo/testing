package handler

import (
	"net/http"

	"github.com/daung-digital/location-api/internal/models"
	"github.com/daung-digital/location-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RoutePlanHandler handles HTTP requests for route plans
type RoutePlanHandler struct {
	service *service.RoutePlanService
}

// NewRoutePlanHandler creates a new route plan handler
func NewRoutePlanHandler(service *service.RoutePlanService) *RoutePlanHandler {
	return &RoutePlanHandler{service: service}
}

// CreateRoutePlan handles POST /route-plans
func (h *RoutePlanHandler) CreateRoutePlan(c *gin.Context) {
	var req models.RoutePlanCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	routePlan, err := h.service.CreateRoutePlan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://api.example.com/errors/internal-server-error",
			"title":  "Internal Server Error",
			"status": http.StatusInternalServerError,
			"detail": "Failed to create route plan",
		})
		return
	}

	c.Header("Location", "/v1/route-plans/"+routePlan.ID.String())
	c.JSON(http.StatusCreated, routePlan)
}

// GetRoutePlan handles GET /route-plans/:id
func (h *RoutePlanHandler) GetRoutePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid route plan ID format",
		})
		return
	}

	routePlan, err := h.service.GetRoutePlan(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/route-plan-not-found",
			"title":  "Route Plan Not Found",
			"status": http.StatusNotFound,
			"detail": "Route plan not found",
		})
		return
	}

	c.JSON(http.StatusOK, routePlan)
}

// ListRoutePlans handles GET /route-plans
func (h *RoutePlanHandler) ListRoutePlans(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	routePlans, nextCursor, hasMore, err := h.service.ListRoutePlans(c.Request.Context(), limit, query.Cursor, query.Sort, query.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://api.example.com/errors/internal-server-error",
			"title":  "Internal Server Error",
			"status": http.StatusInternalServerError,
			"detail": "Failed to list route plans",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": routePlans,
		"pagination": gin.H{
			"next_cursor": nextCursor,
			"has_more":    hasMore,
		},
	})
}

// UpdateRoutePlan handles PATCH /route-plans/:id
func (h *RoutePlanHandler) UpdateRoutePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid route plan ID format",
		})
		return
	}

	var req models.RoutePlanUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	routePlan, err := h.service.UpdateRoutePlan(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/route-plan-not-found",
			"title":  "Route Plan Not Found",
			"status": http.StatusNotFound,
			"detail": "Route plan not found",
		})
		return
	}

	c.JSON(http.StatusOK, routePlan)
}

// DeleteRoutePlan handles DELETE /route-plans/:id
func (h *RoutePlanHandler) DeleteRoutePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid route plan ID format",
		})
		return
	}

	if err := h.service.DeleteRoutePlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/route-plan-not-found",
			"title":  "Route Plan Not Found",
			"status": http.StatusNotFound,
			"detail": "Route plan not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
