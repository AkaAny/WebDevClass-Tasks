package task1

import "webdevclass-tasks/entity"

type CreateUser struct {
	Mail     string
	Password string
}

type UpdateUserPassword struct {
	Password string
}

type UpdateUserRole struct {
	Roles []entity.RoleID
}

type UpdateUserStatus struct {
	Status entity.Status
}

type UserResponse struct {
	UserID string
	Mail   string
	Roles  []entity.RoleID
	Status entity.Status
}
