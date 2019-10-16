package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"log"
	"net/http"
	"strconv"
	"github.com/petrijam/bookstore/dao"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Total int `json:"total"`
	TotalPages int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PerPage int `json:"perPage"`
	Count int `json:"count"`
}

type MetaData struct {
	Pag Pagination `json:"pagination"`
}

type JsonDataBook struct {
	Data []dao.Book `json:"data"`
	Meta MetaData `json:"meta"`
}

type JsonDataComment struct {
	Data []dao.Comment `json:"data"`
	Meta MetaData `json:"meta"`
}

func main() {
	if dao.InitDb() == false {
		return
	}
	handleRequests()
}

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:10000/")
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/books", getBooks).Methods("GET")
	myRouter.HandleFunc("/books", postBook).Methods("POST")
	myRouter.HandleFunc("/books/{id}", getBook).Methods("GET")
	myRouter.HandleFunc("/books/{id}", putBook).Methods("PUT")
	myRouter.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	myRouter.HandleFunc("/books/{bookId}/comments", getComments).Methods("GET")
	myRouter.HandleFunc("/books/{bookId}/comments", postComment).Methods("POST")
	myRouter.HandleFunc("/books/{bookId}/comments/{id}", getComment).Methods("GET")
	myRouter.HandleFunc("/books/{bookId}/comments/{id}", putComment).Methods("PUT")
	myRouter.HandleFunc("/books/{bookId}/comments/{id}", deleteComment).Methods("DELETE")
	http.ListenAndServe(":10000", myRouter)
}

func returnError(w http.ResponseWriter, code int, message string) {
	err := Error {
		Code : code,
		Message : message,
	}
	var jsonData map[string]Error
	jsonData = make(map[string]Error)
	jsonData["error"] = err
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(jsonData)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pageNumber,err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil || pageNumber <= 0 {
		returnError(w, 400, "Bad Request. Invalid Page Number.")
		log.Println(err)
		return
	}
	pageSize,err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		returnError(w, 400, "Bad Request. Invalid Page Size.")
		log.Println(err)
		return
	}	
	books := []dao.Book{} 
	var totalRows int
	if err := dao.GetBooks(&totalRows,&pageNumber,&pageSize,&books); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	countRows := len(books)
	countPages := totalRows / pageSize
	if totalRows % pageSize != 0 {
		countPages = countPages + 1
	}
	pag := Pagination {
		Total : totalRows,
		TotalPages : countPages,
		CurrentPage : pageNumber,
		PerPage : pageSize,
		Count : countRows,
	}	
	metaD := MetaData {
		Pag : pag,
	}
	jData := JsonDataBook {
		Data : books,
		Meta : metaD,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jData)
}

func postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnError(w, 400, "Bad Request.")
		log.Println(err)
		return
	}
	var book dao.Book
	json.Unmarshal(reqBody, &book)
	if bookValidation(w, book) == false {
		return
	}
	if err := dao.PostBook(&book); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["id"])
	if err != nil || bookId < 0 {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}	
	var book dao.Book
	book.ID = uint(bookId)
	if err := dao.GetBook(&book); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			returnError(w, 404, "Record Not Found.")
			log.Println(err)
			return
		}

		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func putBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["id"])
	if err != nil || bookId < 0 {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	var book dao.Book
	book.ID = uint(bookId)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnError(w, 400, "Bad Request.")
		log.Println(err)
		return
	}	
	json.Unmarshal(reqBody, &book)

	if bookValidation(w, book) == false {
		return
	}
	if dao.ValidateBookId(book.ID) == false {
		returnError(w, 404, "Record Not Found.")
		log.Println(err)
		return
	}
	if err := dao.PutBook(&book); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["id"])
	if err != nil || bookId < 0 {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	var book dao.Book
	book.ID = uint(bookId)
	if err := dao.DeleteBook(&book); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			returnError(w, 404, "Record Not Found.")
			log.Println(err)
			return
		}
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func bookValidation(w http.ResponseWriter, book dao.Book) bool {
	if book.Title == "" {
		returnError(w, 400, "Bad Request. Title cannot be empty.")
		return false
	}
	if len(book.Title) > 255 {
		returnError(w, 400, "Bad Request. Title cannot be longer then 255 characters.")
		return false
	}
	if len(book.Author) > 255 {
		returnError(w, 400, "Bad Request. Author cannot be longer then 255 characters.")
		return false
	}
	if len(book.Publisher) > 255 {
		returnError(w, 400, "Bad Request. Publisher cannot be longer then 255 characters.")
		return false
	}
	if len(book.Isbn) > 13 {
		returnError(w, 400, "Bad Request. ISBN number cannot be longer then 13 characters.")
		return false
	}
	if dao.ValidateIsbnNumber(book.Isbn,book.ID) == false {
		returnError(w, 400, "Bad Request. ISBN number must be unique.")
		return false
	}
	if book.Pages <= 0 {
		returnError(w, 400, "Bad Request. Number of pages cannot be less than 1.")
		return false
	}
	t := time.Now()
	if book.PublicationDate.After(t) {
		returnError(w, 400, "Bad Request. Publication date cannot be in the future.")
		return false
	}
	return true
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil || bookId < 0 {
		returnError(w,400,"Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil || pageNumber <= 0 {
		returnError(w, 400, "Bad Request. Invalid Page Number.")
		log.Println(err)
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		returnError(w, 400, "Bad Request. Invalid Page Size.")
		log.Println(err)
		return
	}	
	comments := []dao.Comment{} 
	var totalRows int
	if err := dao.GetComments(&totalRows,&pageNumber,&pageSize,&bookId,&comments); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	countRows := len(comments)
	countPages := totalRows / pageSize
	if totalRows % pageSize != 0 {
		countPages = countPages + 1
	}
	pag := Pagination {
		Total : totalRows,
		TotalPages : countPages,
		CurrentPage : pageNumber,
		PerPage : pageSize,
		Count : countRows,
	}	
	metaD := MetaData {
		Pag : pag,
	}
	jData := JsonDataComment {
		Data : comments,
		Meta : metaD,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jData)
}

func postComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil || bookId < 0 {
		returnError(w,400,"Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnError(w, 400, "Bad Request.")
		log.Println(err)
		return
	}
	var comment dao.Comment
	json.Unmarshal(reqBody, &comment)
	comment.BookID = uint(bookId)

	if commentValidation(w, comment) == false {
		return
	}
	if err := dao.PostComment(&comment); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func getComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil || bookId < 0 {
		returnError(w,400,"Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}	
	commentId, err := strconv.Atoi(vars["id"])
	if err != nil || commentId < 0 {
		returnError(w, 400, "Bad Request. Invalid Comment ID.")
		log.Println(err)
		return
	}	
	var comment dao.Comment
	comment.ID = uint(commentId)
	comment.BookID = uint(bookId)
	if err := dao.GetComment(&comment); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			returnError(w, 404, "Record Not Found.")
			log.Println(err)
			return
		}
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func putComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil || bookId < 0 {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	commentId, err := strconv.Atoi(vars["id"])
	if err != nil || commentId < 0 {
		returnError(w, 400, "Bad Request. Invalid Comment ID.")
		log.Println(err)
		return
	}
	var comment dao.Comment
	comment.BookID = uint(bookId)
	comment.ID = uint(commentId)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnError(w, 400, "Bad Request.")
		log.Println(err)
		return
	}
	json.Unmarshal(reqBody, &comment)
	if commentValidation(w, comment) == false {
		return
	}
	if dao.ValidateCommentId(comment.BookID, comment.ID) == false {
		returnError(w, 404, "Record Not Found.")
		log.Println(err)
		return
	}
	if err := dao.PutComment(&comment); err != nil {
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil || bookId < 0 {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		log.Println(err)
		return
	}
	commentId, err := strconv.Atoi(vars["id"])
	if err != nil || commentId < 0 {
		returnError(w, 400, "Bad Request. Invalid Comment ID.")
		log.Println(err)
		return
	}
	var comment dao.Comment
	comment.BookID = uint(bookId)
	comment.ID = uint(commentId)
	if err := dao.DeleteComment(&comment); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			returnError(w, 404, "Record Not Found.")
			log.Println(err)
			return
		}
		returnError(w, 500, err.Error())
		log.Println(err)
		return
	}		
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func commentValidation(w http.ResponseWriter, comment dao.Comment) bool {
	if comment.Author == "" {
		returnError(w, 400, "Bad Request. Author cannot be empty.")
		return false
	}
	if len(comment.Author) > 255 {
		returnError(w, 400, "Bad Request. Author cannot be longer then 255 characters.")
		return false
	}
	if len(comment.CommentText) > 255 {
		returnError(w, 400, "Bad Request. Comment cannot be longer then 255 characters.")
		return false
	}
	if dao.ValidateBookId(comment.BookID) == false {
		returnError(w, 400, "Bad Request. Invalid Book ID.")
		return false
	}
	return true
}
