package mail

import _ "embed"

//go:embed user_active.html
var UserActiveHtmlTemplate string

type UserActiveTemplateData struct {
	UserID string
	Code   string
}

//go:embed reset_password.html
var ResetPasswordHtmlTemplate string

type ResetPasswordTemplateData struct {
	UserID string
	Code   string
}
