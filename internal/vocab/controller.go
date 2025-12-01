package vocab

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests for vocabulary
type Controller struct {
	service Service
}

// NewController creates a new vocabulary controller
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// RegisterRoutes registers vocabulary routes
func RegisterRoutes(router *gin.Engine, c *Controller) {
	vocab := router.Group("/api/vocabularies")
	// TODO: Add auth middleware
	{
		vocab.POST("", c.Create)
		vocab.GET("", c.GetAll)
		vocab.GET("/:id", c.GetByID)
		vocab.PUT("/:id", c.Update)
		vocab.DELETE("/:id", c.Delete)
	}
}

// getUserID extracts user ID from context (set by auth middleware)
func getUserID(ctx *gin.Context) int64 {
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// Create handles vocabulary creation
func (c *Controller) Create(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateVocabRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vocab, err := c.service.Create(ctx.Request.Context(), userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vocabulary"})
		return
	}

	ctx.JSON(http.StatusCreated, vocab)
}

// GetAll handles getting all vocabularies for a user
func (c *Controller) GetAll(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.GetByUserID(ctx.Request.Context(), userID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vocabularies"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID handles getting a vocabulary by ID
func (c *Controller) GetByID(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	vocab, err := c.service.GetByID(ctx.Request.Context(), userID, id)
	if err != nil {
		switch err {
		case ErrVocabNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vocabulary not found"})
		case ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vocabulary"})
		}
		return
	}

	ctx.JSON(http.StatusOK, vocab)
}

// Update handles vocabulary update
func (c *Controller) Update(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateVocabRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vocab, err := c.service.Update(ctx.Request.Context(), userID, id, &req)
	if err != nil {
		switch err {
		case ErrVocabNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vocabulary not found"})
		case ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vocabulary"})
		}
		return
	}

	ctx.JSON(http.StatusOK, vocab)
}

// Delete handles vocabulary deletion
func (c *Controller) Delete(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), userID, id); err != nil {
		switch err {
		case ErrVocabNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Vocabulary not found"})
		case ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vocabulary"})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
