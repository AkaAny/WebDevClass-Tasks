package task2

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webdevclass-tasks/entity"
	"webdevclass-tasks/model/task2"
	"webdevclass-tasks/utils"
)

func RegisterTokenGroup(db *gorm.DB, engine *gin.Engine) {
	var v2TokenGroup = engine.Group("/task2/token")
	v2TokenGroup.POST("/", func(c *gin.Context) {
		var reqBody = new(task2.LoginRequest)
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
		if userItem.Password != reqBody.Password {
			utils.ResponseWithErr(c, http.StatusForbidden,
				errors.New("password does not match"))
			return
		}
		var userStatus = new(entity.UserStatus)
		db.Where(entity.UserStatus{UserID: userItem.UserID}).First(userStatus)
		if userStatus.Status != entity.UserStatusActive {
			utils.ResponseWithErr(c, http.StatusForbidden,
				errors.New("invalid status"))
			return
		}
		var tokenItem = &entity.AccessToken{
			Token:     uuid.New().String(),
			UserID:    userItem.UserID,
			ExpiredAt: time.Now().Add(24 * time.Hour),
		}
		if err := db.Save(tokenItem).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		var data = &task2.LoginResponse{Token: "token " + tokenItem.Token}
		utils.ResponseSuccess(c, data)
	})
}

func NewTokenValidate(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rawHeader = c.GetHeader("Authorization")
		var tokenParts = strings.Split(rawHeader, " ")
		if len(tokenParts) != 2 {
			c.Abort()
			utils.ResponseWithErr(c, http.StatusUnauthorized, errors.New("invalid header:"+rawHeader))
			return
		}
		var tokenValue = tokenParts[1]
		var tokenItem = new(entity.AccessToken)
		db.Where(entity.AccessToken{Token: tokenValue}).First(tokenItem)
		if tokenItem.Token == "" {
			c.Abort()
			utils.ResponseWithErr(c, http.StatusUnauthorized, errors.New("token does not exist"))
			return
		}
		c.Set("TokenItem", tokenItem)
		c.Next()
	}
}

func NewEnforcePermission(db *gorm.DB, expectedRoleIDs ...entity.RoleID) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenItem = GetTokenFromGinContext(c)
		var roleItems = make([]*entity.UserRole, 0)
		db.Where(entity.UserRole{UserID: tokenItem.UserID}).Find(&roleItems)
		var roleIDs = make([]entity.RoleID, 0)
		for _, roleItem := range roleItems {
			roleIDs = append(roleIDs, roleItem.RoleID)
		}
		if !IsSubCollection[entity.RoleID](roleIDs, expectedRoleIDs) {
			utils.ResponseWithErr(c, http.StatusForbidden, errors.New("access denied"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func Contains[E comparable](u []E, item E) bool {
	for _, uItem := range u {
		if uItem == item {
			return true
		}
	}
	return false
}

func IsSubCollection[E comparable](u []E, sub []E) bool {
	for _, subItem := range sub {
		if !Contains(u, subItem) {
			return false
		}
	}
	return true
}

func GetTokenFromGinContext(c *gin.Context) *entity.AccessToken {
	return c.MustGet("TokenItem").(*entity.AccessToken)
}
