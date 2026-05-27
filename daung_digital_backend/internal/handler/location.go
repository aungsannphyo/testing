package handler

import (
	"net/http"

	"github.com/daung-digital/location-api/internal/models"
	"github.com/daung-digital/location-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LocationHandler handles HTTP requests for locations
type LocationHandler struct {
	service *service.LocationService
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// CreateLocation handles POST /locations
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req models.LocationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	location, err := h.service.CreateLocation(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://api.example.com/errors/internal-server-error",
			"title":  "Internal Server Error",
			"status": http.StatusInternalServerError,
			"detail": "Failed to create location",
		})
		return
	}

	c.Header("Location", "/v1/locations/"+location.ID.String())
	c.JSON(http.StatusCreated, location)
}

// GetLocation handles GET /locations/:id
func (h *LocationHandler) GetLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid location ID format",
		})
		return
	}

	location, err := h.service.GetLocation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/location-not-found",
			"title":  "Location Not Found",
			"status": http.StatusNotFound,
			"detail": "Location not found",
		})
		return
	}

	c.JSON(http.StatusOK, location)
}

// ListLocations handles GET /locations
func (h *LocationHandler) ListLocations(c *gin.Context) {
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

	response, err := h.service.ListLocations(c.Request.Context(), limit, query.Cursor, query.Sort, query.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://api.example.com/errors/internal-server-error",
			"title":  "Internal Server Error",
			"status": http.StatusInternalServerError,
			"detail": "Failed to list locations",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateLocation handles PATCH /locations/:id
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid location ID format",
		})
		return
	}

	var req models.LocationUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	location, err := h.service.UpdateLocation(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/location-not-found",
			"title":  "Location Not Found",
			"status": http.StatusNotFound,
			"detail": "Location not found",
		})
		return
	}

	c.JSON(http.StatusOK, location)
}

// DeleteLocation handles DELETE /locations/:id
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://api.example.com/errors/invalid-uuid",
			"title":  "Invalid UUID",
			"status": http.StatusBadRequest,
			"detail": "Invalid location ID format",
		})
		return
	}

	if err := h.service.DeleteLocation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://api.example.com/errors/location-not-found",
			"title":  "Location Not Found",
			"status": http.StatusNotFound,
			"detail": "Location not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// FindNearby handles GET /locations/nearby
func (h *LocationHandler) FindNearby(c *gin.Context) {
	var query models.NearbyQuery
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

	response, err := h.service.FindNearby(c.Request.Context(), query.Latitude, query.Longitude, query.Radius, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://api.example.com/errors/internal-server-error",
			"title":  "Internal Server Error",
			"status": http.StatusInternalServerError,
			"detail": "Failed to find nearby locations",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
