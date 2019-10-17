package dao

import (
	"time"
	"github.com/jinzhu/gorm"	
)

type Book struct {
	gorm.Model
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Isbn            string    `json:"isbn"`
	Publisher       string    `json:"publisher"`
	PublicationDate time.Time `json:"publicationDate"`
	Pages           int       `json:"pages"`
}

func GetBooks(totalRowsP *int,pageNumberP *int,pageSizeP *int, booksP *[]Book) error {
	if err := db.Model(&Book{}).Count(totalRowsP).Error; err != nil {
		return err
	}
	if err := db.Select("id, created_at, updated_at, deleted_at, title, author, isbn, publisher, publication_date, pages").Limit(*pageSizeP).Offset((*pageNumberP-1) * *pageSizeP).Find(booksP).Error; err != nil {
		return err
	}
	return nil
}

func PostBook(bookP *Book) error {
	if err := db.Create(bookP).Error; err != nil {
		return err
	}
	return nil
}

func GetBook(bookP *Book) error {
	if err := db.Select("id, created_at, updated_at, deleted_at, title, author, isbn, publisher, publication_date, pages").First(bookP).Error; err != nil {
		return err
	}
	return nil
}

func PutBook(bookP *Book) error {
	if err := db.Save(bookP).Error; err != nil {
		return err
	}
	return nil
}

func DeleteBook(bookP *Book) error {
	if err := db.Select("id, created_at, updated_at, deleted_at, title, author, isbn, publisher, publication_date, pages").First(bookP).Error; err != nil {
		return err
	}			
	if err := db.Unscoped().Delete(bookP).Error; err != nil {
		return err
	}
	return nil
}

func ValidateIsbnNumber(isbn string, id uint) bool {
	if db.Select("id").Where("isbn = ? AND id != ?",isbn,id).First(&Book{}).RecordNotFound() {
		return true
	}
	return false
}

func ValidateBookId(id uint) bool {
	if db.Select("id").Where("id = ? ",id).First(&Book{}).RecordNotFound() {
		return false
	}
	return true
}