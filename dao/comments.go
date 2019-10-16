package dao

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	BookID      uint   `json:"bookId"`
	Author      string `json:"author"`
	CommentText string `json:"commentText"`
}

func GetComments(totalRowsP *int,pageNumberP *int,pageSizeP *int,bookId *int,commentsP *[]Comment) error {
	if err := db.Where("book_id = ?",*bookId).Model(&Comment{}).Count(totalRowsP).Error; err != nil {
		return err
	}
	if err := db.Where("book_id = ?",*bookId).Select("id, created_at, updated_at, deleted_at, book_id, author, comment_text").Limit(*pageSizeP).Offset((*pageNumberP-1) * *pageSizeP).Find(commentsP).Error; err != nil {
		return err
	}
	return nil
}

func PostComment(commentP *Comment) error {
	if err := db.Create(commentP).Error; err != nil {
		return err
	}
	return nil
}

func GetComment(commentP *Comment) error {
	if err := db.Select("id, created_at, updated_at, deleted_at, book_id, author, comment_text").Where("book_id = ? AND id = ?",commentP.BookID, commentP.ID).First(commentP).Error; err != nil {
		return err
	}
	return nil
}

func PutComment(commentP *Comment) error {
	if err := db.Save(commentP).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentP *Comment) error {
	if err := db.Select("id, created_at, updated_at, deleted_at, book_id, author, comment_text").Where("book_id = ? AND id = ?",commentP.BookID, commentP.ID).First(commentP).Error; err != nil {
		return err
	}			
	if err := db.Unscoped().Delete(commentP).Error; err != nil {
		return err
	}
	return nil
}

func ValidateCommentId(bookId uint , id uint) bool {
	if db.Select("id").Where("book_id = ? AND id = ?", bookId, id).First(&Comment{}).RecordNotFound() {
		return false
	}
	return true
}