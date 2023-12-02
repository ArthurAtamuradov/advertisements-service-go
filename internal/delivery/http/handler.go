package http

import (
	"log"
	"net/http"
	"strconv"

	"advertisements-service/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AdvertisementHandler struct {
	Service *usecases.AdvertisementService
}

func NewAdvertisementHandler(service *usecases.AdvertisementService) *AdvertisementHandler {
	return &AdvertisementHandler{Service: service}
}

func (h *AdvertisementHandler) GetAllAdvertisements(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	sortBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "asc")

	advertisements, err := h.Service.GetAdvertisementsPage(page, 10, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, advertisements)
}

func (h *AdvertisementHandler) GetAdvertisement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}

	advertisement, err := h.Service.GetAdvertisementByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if advertisement == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Advertisement not found"})
		return
	}

	c.JSON(http.StatusOK, advertisement)
}

// CreateAdvertisement handles the request to create a new advertisement
func (h *AdvertisementHandler) CreateAdvertisement(c *gin.Context) {
	var newAdvertisement usecases.AdvertisementRequest
	if err := c.ShouldBindJSON(&newAdvertisement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := h.Service.CreateAdvertisement(&newAdvertisement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})

	
}


// UpdateAdvertisement handles the request to update an advertisement by ID
func (h *AdvertisementHandler) UpdateAdvertisement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}

	var updatedAdvertisement usecases.AdvertisementUpdateRequest
	if err := c.ShouldBindJSON(&updatedAdvertisement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := h.Service.UpdateAdvertisement(id, &updatedAdvertisement); err != nil {

	log.Printf("Error updating advertisement: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertisement updated successfully"})
}

// DeleteAdvertisement handles the request to delete an advertisement by ID
func (h *AdvertisementHandler) DeleteAdvertisement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}

	if err := h.Service.DeleteAdvertisement(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Advertisement deleted successfully"})
}
