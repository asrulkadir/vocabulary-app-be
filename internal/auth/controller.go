package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests for auth
type Controller struct {
	service Service
}

// NewController creates a new auth controller
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// RegisterRoutes registers auth routes
func RegisterRoutes(router *gin.Engine, c *Controller) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/register", c.Register)
	}
}

// Login handles user login
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Register handles user registration
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
