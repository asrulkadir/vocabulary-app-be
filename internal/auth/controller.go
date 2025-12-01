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
		auth.POST("/logout", c.Logout)
	}
}

// Login handles user login
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Set HTTP-only cookie with token
	ctx.SetCookie(
		"auth_token",
		response.Token,
		86400,
		"/",
		"",
		true,
		true,
	)

	response.Message = "Login successful"
	ctx.JSON(http.StatusOK, response)
}

// Register handles user registration
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrUserAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Set HTTP-only cookie with token
	ctx.SetCookie(
		"auth_token",
		response.Token,
		86400,
		"/",
		"",
		true,
		true,
	)

	response.Message = "Registration successful"
	ctx.JSON(http.StatusCreated, response)
}

// Logout handles user logout
func (c *Controller) Logout(ctx *gin.Context) {
	// Clear the auth cookie
	ctx.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
