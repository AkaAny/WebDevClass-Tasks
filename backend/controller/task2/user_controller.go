package task2

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"net/http"
	"text/template"
	"time"
	"webdevclass-tasks/config"
	"webdevclass-tasks/entity"
	"webdevclass-tasks/model/task2"
	"webdevclass-tasks/service"
	template2 "webdevclass-tasks/template/mail"
	"webdevclass-tasks/template/page"
	"webdevclass-tasks/utils"
)

func RegisterUserGroup(db *gorm.DB, engine *gin.Engine, cfg *config.Config) {
	var v2UserGroup = engine.Group("/task2/user")
	RegisterRegisterHandler(db, v2UserGroup, cfg)
	RegisterActiveHandler(db, v2UserGroup)
	RegisterResetPasswordRequestHandler(db, v2UserGroup, cfg)
	RegisterResetPasswordCallbackHandler(db, v2UserGroup)
	RegisterSelfHandler(db, v2UserGroup)
}

func RegisterRegisterHandler(db *gorm.DB, group gin.IRouter, cfg *config.Config) {
	var handlerFunc = func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID != "" {
			utils.ResponseWithErr(c, http.StatusForbidden,
				errors.New("user:"+userItem.UserID+" has existed"))
			return
		}
		var reqBody = new(task2.RegisterRequest)
		err := c.BindJSON(reqBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		userItem = &entity.User{
			UserID:   userID,
			Mail:     reqBody.Mail,
			Password: reqBody.Password,
		}
		var userStatus = &entity.UserStatus{
			UserID: userItem.UserID,
			Status: entity.UserStatusPendingActive,
		}
		//active code
		var activeItem = &entity.UserActive{
			Code:      uuid.New().String(),
			UserID:    userItem.UserID,
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := db.Save(userItem).Error; err != nil {
				return err
			}
			if err := db.Save(userStatus).Error; err != nil {
				return err
			}
			if err := db.Save(activeItem).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		//send mail
		{
			var mailConfig = cfg.Mail["UserActive"]
			var mailClient = service.NewFromMailConfig(mailConfig)
			//execute template
			t, err := template.New("user_active").Parse(template2.UserActiveHtmlTemplate)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to parse template:%w", err))
				return
			}
			var msgBuf = bytes.NewBuffer(nil)
			err = t.Execute(msgBuf, template2.UserActiveTemplateData{
				UserID: activeItem.UserID,
				Code:   activeItem.Code,
			})
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to execute template:%w", err))
				return
			}
			var mailMessage = gomail.NewMessage()
			mailMessage.SetHeader("From", mailConfig.UserName)
			mailMessage.SetHeader("To", userItem.Mail)
			mailMessage.SetHeader("Subject", "[WebDevClass-Tasks] Active your account")
			mailMessage.SetBody("text/html", msgBuf.String())
			err = mailClient.SendMail(mailMessage)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to send mail"))
				return
			}
		}
		utils.ResponseSuccess(c, nil)
	}
	group.POST("/:userID", handlerFunc)
}

func RegisterActiveHandler(db *gorm.DB, group gin.IRouter) {
	var handlerFunc = func(c *gin.Context) {
		var reqBody = new(task2.UserActiveRequest)
		err := c.BindJSON(reqBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var activeItem = new(entity.UserActive)
		db.Where(entity.UserActive{Code: reqBody.Code}).First(activeItem)
		if activeItem.Code == "" {
			utils.ResponseWithErr(c, http.StatusNotFound,
				errors.New("code:"+reqBody.Code+" does not exist"))
			return
		}
		var userStatus = new(entity.UserStatus)
		db.Where(entity.UserStatus{UserID: activeItem.UserID}).First(userStatus)
		if userStatus.UserID == "" {
			utils.ResponseWithErr(c, http.StatusNotFound,
				errors.New("user:"+activeItem.UserID+" does not exist"))
			return
		}
		userStatus.Status = entity.UserStatusActive
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(activeItem).Delete(&entity.UserActive{}).Error; err != nil {
				return err
			}
			if err := tx.Save(userStatus).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	}
	group.PUT("/active", handlerFunc)
}

func RegisterResetPasswordRequestHandler(db *gorm.DB, group gin.IRouter, cfg *config.Config) {
	var handlerFunc = func(c *gin.Context) {
		var reqBody = new(task2.ResetPasswordRequestRequest)
		err := c.BindJSON(reqBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var userItem = new(entity.User)
		db.Where(entity.User{Mail: reqBody.Mail}).First(userItem)
		if userItem.UserID == "" {
			utils.ResponseWithErr(c, http.StatusNotFound,
				errors.New("user with mail:"+reqBody.Mail+" does not exist"))
			return
		}
		var resetItem = &entity.UserResetPassword{
			Code:      uuid.New().String(),
			UserID:    userItem.UserID,
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		if err := db.Save(resetItem).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError,
				fmt.Errorf("failed to save reset item to db:%w", err))
			return
		}
		{
			var mailConfig = cfg.Mail["ResetPassword"]
			var mailClient = service.NewFromMailConfig(mailConfig)
			//execute template
			t, err := template.New("reset_password").Parse(template2.ResetPasswordHtmlTemplate)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to parse template:%w", err))
				return
			}
			var msgBuf = bytes.NewBuffer(nil)
			err = t.Execute(msgBuf, template2.ResetPasswordTemplateData{
				UserID: resetItem.UserID,
				Code:   resetItem.Code,
			})
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to execute template:%w", err))
				return
			}
			var mailMessage = gomail.NewMessage()
			mailMessage.SetHeader("From", mailConfig.UserName)
			mailMessage.SetHeader("To", userItem.Mail)
			mailMessage.SetHeader("Subject", "[WebDevClass-Tasks] Reset your account's password")
			mailMessage.SetBody("text/html", msgBuf.String())
			err = mailClient.SendMail(mailMessage)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError,
					fmt.Errorf("failed to send mail"))
				return
			}
		}
		utils.ResponseSuccess(c, nil)
	}
	group.POST("/resetPassword/request", handlerFunc)
}

func RegisterResetPasswordCallbackHandler(db *gorm.DB, group gin.IRouter) {
	var handlerFunc = func(c *gin.Context) {
		var reqBody = new(task2.ResetPasswordCallbackRequest)
		err := c.BindJSON(reqBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var resetItem = new(entity.UserResetPassword)
		db.Where(entity.UserResetPassword{Code: reqBody.Code}).First(resetItem)
		if resetItem.Code == "" {
			utils.ResponseWithErr(c, http.StatusNotFound, errors.New("code:"+reqBody.Code+"does not exist"))
			return
		}
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: resetItem.UserID}).First(userItem)
		if userItem.UserID == "" {
			utils.ResponseWithErr(c, http.StatusNotFound, errors.New("user:"+userItem.UserID+"does not exist"))
			return
		}
		userItem.Password = reqBody.NewPassword
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(resetItem).Delete(&entity.UserResetPassword{}).Error; err != nil {
				return err
			}
			if err := tx.Save(userItem).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
	}
	group.PUT("/resetPassword/callback", handlerFunc)
}

func RegisterSelfHandler(db *gorm.DB, group gin.IRouter) {
	group.GET("/self/:type", NewTokenValidate(db), func(c *gin.Context) {
		var retType = c.Param("type")
		var tokenItem = GetTokenFromGinContext(c)
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: tokenItem.UserID}).First(userItem)
		if userItem.UserID == "" {
			utils.ResponseWithErr(c, http.StatusNotFound, errors.New("user:"+userItem.UserID+"does not exist"))
			return
		}
		//map role ids
		var roleItems = make([]*entity.UserRole, 0)
		db.Where(entity.UserRole{UserID: userItem.UserID}).Find(&roleItems)
		var roleIDs = make([]entity.RoleID, 0)
		for _, roleItem := range roleItems {
			roleIDs = append(roleIDs, roleItem.RoleID)
		}
		var userStatus = new(entity.UserStatus)
		db.Where(entity.UserStatus{UserID: userItem.UserID}).First(userStatus)
		switch retType {
		case "json":
			var data = &task2.UserResponse{
				UserID: userItem.UserID,
				Mail:   userItem.Mail,
				Roles:  roleIDs,
				Status: userStatus.Status,
			}
			utils.ResponseSuccess(c, data)
		case "html":
			var data = page.ParseAndExecuteUserInfoTemplate(page.UserInfoTemplateData{
				UserID: userItem.UserID,
				Mail:   userItem.Mail,
				Roles:  roleIDs,
				Status: userStatus.Status,
			})
			c.Data(http.StatusOK, "text/html", data)
		default:
			utils.ResponseWithErr(c, http.StatusInternalServerError, errors.New("failed to execute template"))
		}
	})
}
