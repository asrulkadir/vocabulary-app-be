package auth

import (
	"net/http"

	"vocabulary-app-be/pkg/utils"

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
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrInvalidCredentials:
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid email or password")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
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

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}

// Register handles user registration
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		switch err {
		case ErrUserAlreadyExists:
			utils.ErrorResponse(ctx, http.StatusConflict, "User already exists")
		default:
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
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

	utils.SuccessResponse(ctx, http.StatusCreated, "Registration successful", response)
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

	utils.SuccessResponse(ctx, http.StatusOK, "Logout successful", nil)
}
