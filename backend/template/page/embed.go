package page

import (
	"bytes"
	_ "embed"
	template3 "html/template"
	"webdevclass-tasks/entity"
)

//go:embed user_info.html
var UserInfoTemplate string

const UserInfoTemplateName = "userInfo"

type UserInfoTemplateData struct {
	UserID string
	Mail   string
	Roles  []entity.RoleID
	Status entity.Status
}

func ParseAndExecuteUserInfoTemplate(data UserInfoTemplateData) []byte {
	tmpl, err := template3.New(UserInfoTemplateName).Parse(UserInfoTemplate)
	if err != nil {
		panic(err)
	}
	var buf = bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
