package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"webdevclass-tasks/entity"
	"webdevclass-tasks/model/task1"
	"webdevclass-tasks/utils"
)

func InitAndRun(db *gorm.DB) {
	//增删改查
	var engine = gin.Default()
	var err error
	var v1UserGroup = engine.Group("/task1/user")
	v1UserGroup.POST("/:userID", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID != "" {
			utils.ResponseWithErr(c, http.StatusForbidden, errors.New("user:"+userID+" has existed"))
			return
		}
		var requestBody = new(task1.CreateUser)
		err = c.BindJSON(requestBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		userItem = &entity.User{
			UserID:   userID,
			Password: requestBody.Password,
		}
		var userStatus = &entity.UserStatus{
			UserID: userID,
			Status: entity.UserStatusPendingActive,
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := db.Save(userItem).Error; err != nil {
				return err
			}
			if err := db.Save(userStatus).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v1UserGroup.DELETE("/:userID", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID == "" {
			responseUserDoesNotExist(c, userID)
			return
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(userItem).Delete(&entity.User{}).Error; err != nil {
				return err
			}
			//delete associated data
			if err := tx.Where(entity.UserRole{UserID: userID}).Delete(&entity.UserRole{}).Error; err != nil {
				return err
			}
			if err := tx.Where(entity.UserStatus{UserID: userID}).Delete(&entity.UserStatus{}).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v1UserGroup.PUT("/:userID/password", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID == "" {
			responseUserDoesNotExist(c, userID)
		}
		var requestBody = new(task1.UpdateUserPassword)
		err = c.BindJSON(requestBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		userItem.Password = requestBody.Password
		if err := db.Save(userItem).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v1UserGroup.PUT("/:userID/role", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID == "" {
			responseUserDoesNotExist(c, userID)
			return
		}
		var requestBody = new(task1.UpdateUserRole)
		err = c.BindJSON(requestBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var newRoleItems = make([]*entity.UserRole, 0)
		for _, roleID := range requestBody.Roles {
			newRoleItems = append(newRoleItems, &entity.UserRole{
				UserID: userItem.UserID,
				RoleID: roleID,
			})
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(entity.UserRole{UserID: userID}).Delete(entity.UserRole{}).Error; err != nil {
				return err
			}
			if err := tx.Save(newRoleItems).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v1UserGroup.PUT("/:userID/status", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID == "" {
			responseUserDoesNotExist(c, userID)
			return
		}
		var requestBody = new(task1.UpdateUserStatus)
		err = c.BindJSON(requestBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var userStatus = new(entity.UserStatus)
		//是否使用事务不是根据并发，而是根据连续操作的完整性和连续性来决定
		db.Where(entity.UserStatus{UserID: userID}).First(userStatus)
		userStatus.Status = requestBody.Status
		if err := db.Save(userStatus).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v1UserGroup.GET("/:userID", func(c *gin.Context) {
		var userID = c.Param("userID")
		var userItem = new(entity.User)
		db.Where(entity.User{UserID: userID}).First(userItem)
		if userItem.UserID == "" {
			responseUserDoesNotExist(c, userID)
			return
		}
		//map role id
		var userRoles = make([]*entity.UserRole, 0)
		db.Where(entity.UserRole{UserID: userID}).Find(&userRoles)
		var roleIDs = make([]entity.RoleID, 0)
		for _, roleItem := range userRoles {
			roleIDs = append(roleIDs, roleItem.RoleID)
		}
		//map status
		var userStatus = new(entity.UserStatus)
		db.Where(entity.UserStatus{UserID: userID}).First(userStatus)
		var data = &task1.UserResponse{
			UserID: userItem.UserID,
			Mail:   userItem.Mail,
			Roles:  roleIDs,
			Status: userStatus.Status,
		}
		utils.ResponseSuccess(c, data)
	})
	err = engine.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func responseUserDoesNotExist(c *gin.Context, userID string) {
	utils.ResponseWithErr(c, http.StatusNotFound, errors.New("user:"+userID+" does not exist"))
}
