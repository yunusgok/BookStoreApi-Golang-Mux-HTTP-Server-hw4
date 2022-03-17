package library

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yunusgok/go-patika/infrastructure"

	"github.com/yunusgok/go-patika/csv_helper"
)

var Books []*Book
var BookNames []string

var ErrBookNotFound = errors.New("Book not found")

var ErrBookDeleted = errors.New("Book is deleted and can not be found")

var ErrNotEnoughStock = errors.New("not enough stock")

var (
	bookRepository *BookRepository
)

//initilize books with book names and authors in list
func InitBooks() {
	BooksList, err := csv_helper.ReadCsv("books.csv")
	if err != nil {
		fmt.Println("book not readed")
		panic(err)
	}
	for _, b := range BooksList {
		Books = append(Books, NewBook(b[0], b[1]))
	}
}
func InitRepo() {
	db := infrastructure.NewPostgresDB("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	bookRepository = NewBookRepository(db)
	bookRepository.Migration()
	bookRepository.InsertSampleData()
}

// List all books line by line by their name
func ListBooks() {
	books := bookRepository.FindAll()
	for _, b := range books {
		fmt.Printf("Book: %s -- Author: %s -- ISBN: %d\n", b.Name, b.Author, b.ISBN)
	}
}

// List given books line by line by their name
func ListGivenBooks(books []Book) {
	for _, b := range books {
		fmt.Printf("Book: %s -- Author: %s -- ISBN: %d\n", b.Name, b.Author, b.ISBN)
	}
}

// Searches given words in books and return matched books names
func FindBooks(word string) []Book {

	// word is turned to lowercase to search case insensitive
	searchWord := strings.ToLower(word)
	// check word is integer so ISBN number can be searched
	var result []Book
	isInteger, value := IsInt(searchWord)
	if !isInteger {
		result = bookRepository.FindByName(word)
	} else {
		result = bookRepository.FindByISBN(word)
		book := bookRepository.GetById(value)
		if book != nil {
			result = append(result, *book)
		}
	}
	if len(result) == 0 {
		fmt.Printf("No book found with given parameter '%s'", word)
	}

	return result
}

//Find book with id
func FindBook(id int) (Book, error) {

	book := bookRepository.GetById(id)
	if book == nil {
		return *new(Book), ErrBookNotFound
	}
	if book.IsDeleted {
		return *new(Book), ErrBookDeleted
	}

	return *book, nil
}

//Buy book if enoubh count exist in stock
func Buy(id int, count int) {
	book, err := FindBook(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err2 := book.Buy(count)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	bookRepository.Update(book)

}

//check given string is int and return to value
// s-> stirng to be checked
// return (check bool, value int)
// check -> whether string is int
// value -> value of the string
func IsInt(s string) (bool, int) {
	if value, err := strconv.Atoi(s); err == nil {
		return true, value
	}
	return false, 0
}

// deletes book if exist
func DeleteBook(id int) {
	book, err := FindBook(id)
	if err == nil {
		err2 := book.Delete()
		if err2 != nil {
			fmt.Println(err2.Error())
			return
		}
	} else {
		fmt.Println(err.Error())
		return
	}
	bookRepository.Update(book)

}
