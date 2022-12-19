package main

import (
	"webdevclass-tasks/cmd/common"
	controller "webdevclass-tasks/controller/task1"
	"webdevclass-tasks/entity"
)

func main() {
	_, db := common.InitConfigAndDB()
	err := db.AutoMigrate(&entity.User{}, &entity.UserRole{}, &entity.UserStatus{})
	if err != nil {
		panic(err)
	}
	//增删改查
	controller.InitAndRun(db)
}
