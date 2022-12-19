package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"webdevclass-tasks/cmd/common"
	"webdevclass-tasks/controller/task2"
	"webdevclass-tasks/controller/task3"
	"webdevclass-tasks/entity"
)

func main() {
	cfg, db := common.InitConfigAndDB()
	err := db.AutoMigrate(&entity.User{}, &entity.UserRole{}, &entity.UserStatus{}, &entity.UserActive{},
		&entity.UserResetPassword{}, &entity.AccessToken{},
		&entity.Book{}, &entity.BookRent{})
	if err != nil {
		panic(err)
	}
	var engine = gin.Default()
	var corsConfig = cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddExposeHeaders("Authorization")
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))
	task2.RegisterUserGroup(db, engine, cfg)
	task2.RegisterTokenGroup(db, engine)
	task3.RegisterBookGroup(db, engine)
	task3.RegisterBookRentGroup(db, engine)
	err = engine.Run(":8083")
	if err != nil {
		panic(err)
	}
}
