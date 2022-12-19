package task2

import "webdevclass-tasks/entity"

type ResetPasswordRequestRequest struct {
	Mail string `json:"mail"`
}

type ResetPasswordCallbackRequest struct {
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

type RegisterRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type UserActiveRequest struct {
	Code string `json:"code"`
}

type UserResponse struct {
	UserID string          `json:"userID"`
	Mail   string          `json:"mail"`
	Roles  []entity.RoleID `json:"roles"`
	Status entity.Status   `json:"status"`
}
