package task3

import "time"

type AddBookRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type UpdateBookRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type AddBookRentRequest struct {
	BookID uint `json:"bookID"`
}

type BookRentResponse struct {
	ID         uint          `json:"id"`
	Book       *BookResponse `json:"book"`
	RentBy     string        `json:"rentBy"`
	ExpiredAt  time.Time     `json:"expiredAt"`
	ReturnedAt time.Time     `json:"returnedAt"`
}
