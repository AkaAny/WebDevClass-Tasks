package task3

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
	"webdevclass-tasks/controller/task2"
	"webdevclass-tasks/entity"
	"webdevclass-tasks/model/task3"
	"webdevclass-tasks/utils"
)

func RegisterBookGroup(db *gorm.DB, engine *gin.Engine) {
	var v3BookGroup = engine.Group("/task3/book")
	v3BookGroup.POST("/",
		task2.NewTokenValidate(db), task2.NewEnforcePermission(db, entity.RoleManageBook),
		func(c *gin.Context) {
			var reqBody = new(task3.AddBookRequest)
			err := c.BindJSON(reqBody)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusBadRequest, err)
				return
			}
			var bookItem = new(entity.Book)
			db.Where(entity.Book{Name: reqBody.Name}).First(bookItem)
			if bookItem.ID != 0 {
				utils.ResponseWithErr(c, http.StatusBadRequest,
					errors.New("book with name:"+reqBody.Name+" does not exist"))
				return
			}
			bookItem = &entity.Book{
				Name:   reqBody.Name,
				Author: reqBody.Author,
			}
			if err := db.Save(bookItem).Error; err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError, err)
				return
			}
			utils.ResponseSuccess(c, nil)
		})
	v3BookGroup.DELETE("/:id",
		task2.NewTokenValidate(db), task2.NewEnforcePermission(db, entity.RoleManageBook),
		func(c *gin.Context) {
			bookID, err := parseBookID(c, "id")
			if err != nil {
				utils.ResponseWithErr(c, http.StatusBadRequest, err)
				return
			}
			var bookItem = new(entity.Book)
			db.Where(entity.Book{ID: bookID}).First(bookItem)
			if bookItem.ID == 0 {
				utils.ResponseWithErr(c, http.StatusNotFound, fmt.Errorf("book:%d does not exist", bookID))
				return
			}
			err = db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Where(entity.BookRent{BookID: bookItem.ID}).Delete(&entity.BookRent{}).Error; err != nil {
					return err
				}
				if err := tx.Where(bookItem).Delete(&entity.Book{}).Error; err != nil {
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
	v3BookGroup.PUT("/:id",
		task2.NewTokenValidate(db), task2.NewEnforcePermission(db, entity.RoleManageBook),
		func(c *gin.Context) {
			bookID, err := parseBookID(c, "id")
			if err != nil {
				utils.ResponseWithErr(c, http.StatusBadRequest, err)
				return
			}
			var bookItem = new(entity.Book)
			db.Where(entity.Book{ID: bookID}).First(bookItem)
			if bookItem.ID == 0 {
				utils.ResponseWithErr(c, http.StatusNotFound, fmt.Errorf("book:%d does not exist", bookID))
				return
			}
			var reqBody = new(task3.AddBookRequest)
			err = c.BindJSON(reqBody)
			if err != nil {
				utils.ResponseWithErr(c, http.StatusBadRequest, err)
				return
			}
			bookItem.Name = reqBody.Name
			bookItem.Author = reqBody.Author
			if err := db.Save(bookItem).Error; err != nil {
				utils.ResponseWithErr(c, http.StatusInternalServerError, err)
				return
			}
			utils.ResponseSuccess(c, nil)
		})
	v3BookGroup.GET("/", func(c *gin.Context) {
		var bookItems = make([]*entity.Book, 0)
		db.Find(&bookItems)
		var dataList = make([]*task3.BookResponse, 0)
		for _, bookItem := range bookItems {
			var data = mapBookItem(bookItem)
			dataList = append(dataList, data)
		}
		utils.ResponseSuccess(c, dataList)
	})
}

func RegisterBookRentGroup(db *gorm.DB, engine *gin.Engine) {
	var v3BookRentGroup = engine.Group("/task3/book/rent")
	v3BookRentGroup.POST("/", task2.NewTokenValidate(db), func(c *gin.Context) {
		var tokenItem = task2.GetTokenFromGinContext(c)
		var reqBody = new(task3.AddBookRentRequest)
		err := c.BindJSON(reqBody)
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var rentItem = new(entity.BookRent)
		db.Where(entity.BookRent{
			BookID: reqBody.BookID,
		}).Last(rentItem)
		if rentItem.ID != 0 { //这本书已经被借出且没有归还
			utils.ResponseWithErr(c, http.StatusForbidden,
				fmt.Errorf("book:%d has been rented", reqBody.BookID))
			return
		}
		rentItem = &entity.BookRent{
			BookID:    reqBody.BookID,
			RentBy:    tokenItem.UserID,
			ExpiredAt: time.Now().Add(24 * time.Hour * 30),
		}
		if err := db.Save(rentItem).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v3BookRentGroup.PUT("/return/:id", task2.NewTokenValidate(db), func(c *gin.Context) {
		recordID, err := parseRecordID(c, "id")
		if err != nil {
			utils.ResponseWithErr(c, http.StatusBadRequest, err)
			return
		}
		var rentItem = new(entity.BookRent)
		db.Where(entity.BookRent{Model: gorm.Model{ID: recordID}}).First(rentItem)
		if rentItem.ID == 0 {
			utils.ResponseWithErr(c, http.StatusNotFound,
				fmt.Errorf("book rent record:%d does not exist", recordID))
			return
		}
		if err := db.Where(rentItem).Delete(&entity.BookRent{}).Error; err != nil {
			utils.ResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
		utils.ResponseSuccess(c, nil)
	})
	v3BookRentGroup.GET("/user/self", task2.NewTokenValidate(db), func(c *gin.Context) {
		var tokenItem = task2.GetTokenFromGinContext(c)
		var rentItems = make([]*entity.BookRent, 0)
		db.Unscoped().Where(entity.BookRent{RentBy: tokenItem.UserID}).Find(&rentItems)
		var dataList = make([]*task3.BookRentResponse, 0)
		for _, rentItem := range rentItems {
			var bookItem = new(entity.Book)
			db.Where(entity.Book{ID: rentItem.BookID}).First(bookItem)
			var data = &task3.BookRentResponse{
				ID:         rentItem.ID,
				Book:       mapBookItem(bookItem),
				RentBy:     rentItem.RentBy,
				ExpiredAt:  rentItem.ExpiredAt,
				ReturnedAt: rentItem.DeletedAt.Time,
			}
			dataList = append(dataList, data)
		}
		utils.ResponseSuccess(c, dataList)
	})
}

func parseBookID(c *gin.Context, paramKey string) (uint, error) {
	var bookIDStr = c.Param(paramKey)
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(bookID), nil
}

func parseRecordID(c *gin.Context, paramKey string) (uint, error) {
	var bookIDStr = c.Param(paramKey)
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(bookID), nil
}

func mapBookItem(item *entity.Book) *task3.BookResponse {
	var data = &task3.BookResponse{
		ID:     item.ID,
		Name:   item.Name,
		Author: item.Author,
	}
	return data
}
