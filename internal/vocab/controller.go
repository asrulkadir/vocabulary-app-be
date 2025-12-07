package vocab

import (
	"net/http"
	"strconv"

	"vocabulary-app-be/pkg/middleware"
	"vocabulary-app-be/pkg/utils"

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
func RegisterRoutes(router *gin.Engine, c *Controller, jwtSecret string) {
	vocab := router.Group("/api/vocabularies")
	// Add auth middleware
	vocab.Use(middleware.AuthMiddleware(jwtSecret))
	{
		vocab.POST("", c.Create)
		vocab.GET("", c.GetAll)
		vocab.GET("/:id", c.GetByID)
		vocab.PUT("/:id", c.Update)
		vocab.DELETE("/:id", c.Delete)
		vocab.POST("/:id/test-result", c.UpdateTestResult)
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
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateVocabRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vocab, err := c.service.Create(ctx.Request.Context(), userID, &req)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create vocabulary")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Vocabulary created successfully", vocab)
}

// GetAll handles getting all vocabularies for a user
func (c *Controller) GetAll(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.GetByUserID(ctx.Request.Context(), userID, page, pageSize)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vocabularies")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabularies retrieved successfully", response)
}

// GetByID handles getting a vocabulary by ID
func (c *Controller) GetByID(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	vocab, err := c.service.GetByID(ctx.Request.Context(), userID, id)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrVocabNotFound:
			utils.ErrorResponse(ctx, http.StatusNotFound, "Vocabulary not found")
		case ErrUnauthorized:
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vocabulary")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabulary retrieved successfully", vocab)
}

// Update handles vocabulary update
func (c *Controller) Update(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req UpdateVocabRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vocab, err := c.service.Update(ctx.Request.Context(), userID, id, &req)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrVocabNotFound:
			utils.ErrorResponse(ctx, http.StatusNotFound, "Vocabulary not found")
		case ErrUnauthorized:
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update vocabulary")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabulary updated successfully", vocab)
}

// Delete handles vocabulary deletion
func (c *Controller) Delete(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), userID, id); err != nil {
		ctx.Error(err)
		switch err {
		case ErrVocabNotFound:
			utils.ErrorResponse(ctx, http.StatusNotFound, "Vocabulary not found")
		case ErrUnauthorized:
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete vocabulary")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "Vocabulary deleted successfully", nil)
}

// UpdateTestResult handles updating test result for a vocabulary
func (c *Controller) UpdateTestResult(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req TestResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vocab, err := c.service.UpdateTestResult(ctx.Request.Context(), userID, id, req.Passed)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrVocabNotFound:
			utils.ErrorResponse(ctx, http.StatusNotFound, "Vocabulary not found")
		case ErrUnauthorized:
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update test result")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test result updated successfully", vocab)
}
