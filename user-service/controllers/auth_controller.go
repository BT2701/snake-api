package controllers

import (
	"context"
	"net/http"

	"snake_api/models"
	"snake_api/services"

	// "snake_api/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// Chuẩn hóa phản hồi
type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func newAPIResponse(status int, data interface{}, err interface{}) *APIResponse {
	return &APIResponse{
		Status: status,
		Data:   data,
		Error:  err,
	}
}

func (ctrl *UserController) Login(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
    }

    // Gọi hàm service để login
    token, err, user := ctrl.service.Login(context.Background(), input.Email, input.Password)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, newAPIResponse(http.StatusUnauthorized, nil, err.Error()))
    }

    // Trả về phản hồi khi login thành công
    return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
        "message": "Login successful",
        "token":   token,
        "user":    user, // Bao gồm thông tin user (nếu cần)
    }, nil))
}


func (ctrl *UserController) SignUp(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.SignUp(context.Background(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, newAPIResponse(http.StatusConflict, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "User created",
		"user":    user,
	}, nil))
}

func (ctrl *UserController) ForgotPassword(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid email"))
	}

	resetToken, err := ctrl.service.ForgotPassword(context.Background(), input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message":     "Password reset email sent successfully",
		"reset_token": resetToken, // Optional: Include reset token for testing purposes
	}, nil))
}

func (ctrl *UserController) ResetPassword(c echo.Context) error {
	var input struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.ResetPassword(context.Background(), input.Token, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
	}, nil))
}

func (ctrl *UserController) GetAllUsers(c echo.Context) error {
	users, err := ctrl.service.GetAllUsers(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, users, nil))
}
