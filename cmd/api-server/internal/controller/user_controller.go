package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-api-server-sample/cmd/api-server/internal/application"
	"go-api-server-sample/cmd/api-server/internal/controller/dtos"
	"go-api-server-sample/internal/domain/entities"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService application.UserServiceInterface
}

// NewUserController creates a new UserController
func NewUserController(userService application.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.CreateUserRequest true "User information"
// @Success 201 {object} dtos.UserResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 409 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dtos.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.handleValidationError(ctx, err)
		return
	}

	user, err := c.userService.CreateUser(ctx.Request.Context(), req.ToEntity())
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dtos.FromEntity(user)
	ctx.JSON(http.StatusCreated, response)
}

// GetUsers handles GET /api/v1/users
// @Summary Get users list
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to retrieve (default: 10, max: 100)" minimum(1) maximum(100) default(10)
// @Param offset query int false "Number of users to skip (default: 0)" minimum(0) default(0)
// @Success 200 {object} dtos.UserListResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /api/v1/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	req := &application.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
	}

	result, err := c.userService.ListUsers(ctx.Request.Context(), req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := &dtos.UserListResponse{
		Users:  dtos.FromEntities(result.Users),
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetUserByID handles GET /api/v1/users/{id}
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := c.parseUserID(ctx)
	if err != nil {
		return // Error already handled
	}

	user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dtos.FromEntity(user)
	ctx.JSON(http.StatusOK, response)
}

// UpdateUser handles PUT /api/v1/users/{id}
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param user body dtos.UpdateUserRequest true "Updated user information"
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 409 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := c.parseUserID(ctx)
	if err != nil {
		return // Error already handled
	}

	var req dtos.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.handleValidationError(ctx, err)
		return
	}

	user, err := c.userService.UpdateUser(ctx.Request.Context(), id, req.ToEntity())
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dtos.FromEntity(user)
	ctx.JSON(http.StatusOK, response)
}

// DeleteUser handles DELETE /api/v1/users/{id}
// @Summary Delete user
// @Description Soft delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Success 204
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := c.parseUserID(ctx)
	if err != nil {
		return // Error already handled
	}

	err = c.userService.DeleteUser(ctx.Request.Context(), id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// parseUserID parses and validates the user ID from the URL path
func (c *UserController) parseUserID(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		errorResp := dtos.NewErrorResponse("VALIDATION_ERROR", "無効なユーザーIDです")
		ctx.JSON(http.StatusBadRequest, errorResp)
		if err != nil {
			return 0, err
		}
		return 0, errors.New("invalid user ID")
	}
	return uint(id), nil
}

// handleError handles domain and application errors
func (c *UserController) handleError(ctx *gin.Context, err error) {
	// Map domain errors to appropriate HTTP status codes
	switch err {
	case entities.ErrUserNotFound:
		errorResp := dtos.MapDomainErrorToResponse(err)
		ctx.JSON(http.StatusNotFound, errorResp)
	case entities.ErrUserEmailExists:
		errorResp := dtos.MapDomainErrorToResponse(err)
		ctx.JSON(http.StatusConflict, errorResp)
	case entities.ErrInvalidEmail, entities.ErrUserNameRequired, entities.ErrUserNameTooLong, entities.ErrInvalidUserID:
		errorResp := dtos.MapDomainErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, errorResp)
	default:
		// Log the actual error for debugging
		// logger.Error("Unexpected error", "error", err)
		errorResp := dtos.NewErrorResponse("INTERNAL_ERROR", "内部サーバーエラーが発生しました")
		ctx.JSON(http.StatusInternalServerError, errorResp)
	}
}

// handleValidationError handles JSON binding validation errors
func (c *UserController) handleValidationError(ctx *gin.Context, _ error) {
	// For now, return a generic validation error
	// In a more sophisticated implementation, you would parse the validation errors
	// and return detailed field-specific error information
	errorResp := dtos.NewErrorResponse("VALIDATION_ERROR", "リクエストのバリデーションに失敗しました")
	ctx.JSON(http.StatusBadRequest, errorResp)
}