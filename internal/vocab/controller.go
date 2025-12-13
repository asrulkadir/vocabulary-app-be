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
		vocab.GET("/stats", c.GetStats)
		vocab.GET("/:id", c.GetByID)
		vocab.PUT("/:id", c.Update)
		vocab.DELETE("/:id", c.Delete)
	}

	// Test-specific routes
	test := router.Group("/api/test")
	test.Use(middleware.AuthMiddleware(jwtSecret))
	{
		test.GET("/vocabularies", c.GetRandomForTest)
		test.GET("/vocabularies/:id/options", c.GetTestOptions)
		test.POST("/vocabularies/:id/answer", c.SubmitTestAnswer)
	}
}

// getUserID extracts user ID from context (set by auth middleware)
func getUserID(ctx *gin.Context) string {
	userID, exists := ctx.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}

// Create handles vocabulary creation
func (c *Controller) Create(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
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
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	search := ctx.Query("search")
	status := ctx.Query("status")

	// Validate status if provided
	if status != "" && status != "all" && status != "learning" && status != "memorized" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid status. Use: all, learning, or memorized")
		return
	}

	// Treat "all" as empty to get all statuses
	if status == "all" {
		status = ""
	}

	response, err := c.service.GetByUserID(ctx.Request.Context(), userID, page, pageSize, search, status)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vocabularies")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabularies retrieved successfully", response)
}

// GetStats handles getting vocabulary statistics
func (c *Controller) GetStats(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stats, err := c.service.GetVocabStats(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vocabulary stats")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabulary stats retrieved successfully", stats)
}

// GetByID handles getting a vocabulary by ID
func (c *Controller) GetByID(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
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
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
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
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
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

// GetRandomForTest handles getting a random vocabulary for testing
func (c *Controller) GetRandomForTest(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	status := ctx.DefaultQuery("status", "all")
	// Validate status
	if status != "all" && status != "learning" && status != "memorized" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid status. Use: all, learning, or memorized")
		return
	}

	vocab, err := c.service.GetRandomForTest(ctx.Request.Context(), userID, status)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrNoVocabsAvailable:
			utils.ErrorResponse(ctx, http.StatusNotFound, "No vocabularies available for testing")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vocabulary for test")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Vocabulary retrieved successfully", vocab)
}

// GetTestOptions handles getting multiple-choice options for a vocabulary
func (c *Controller) GetTestOptions(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vocabID := ctx.Param("id")
	if vocabID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid vocabulary ID")
		return
	}

	options, err := c.service.GetTestOptions(ctx.Request.Context(), userID, vocabID)
	if err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get test options")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test options retrieved successfully", options)
}

// SubmitTestAnswer handles validating a test answer
func (c *Controller) SubmitTestAnswer(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == "" {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req TestResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result, err := c.service.ValidateTestAnswer(ctx.Request.Context(), userID, id, req.Input)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrVocabNotFound:
			utils.ErrorResponse(ctx, http.StatusNotFound, "Vocabulary not found")
		case ErrUnauthorized:
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to validate test answer")
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test answer validated successfully", result)
}
