package library

import (
	"fmt"
	"math/rand"
)

var ID int = 1

type Book struct {
	id         int
	name       string
	stockCode  int
	ISBN       int
	pageCount  int
	price      float64
	stockCount int
	author     string
	isDeleted  bool
}

// Book constructor
func NewBook(name string, author string) *Book {
	book := new(Book)
	book.name = name
	book.author = author
	book.id = ID
	//Seed is current time to give randomness
	// page count will be in range 300-400
	book.pageCount = rand.Intn(300) + 100
	// price will be in range 20.00- 220.00
	book.price = rand.Float64()*200 + 20
	// ISBN will be in range 100000 - 1000000
	book.ISBN = rand.Intn(100000) + 100000
	// stock count  will be in range 0-50
	book.stockCount = rand.Intn(50)
	// stock code  will be in range 100000 - 1000000
	book.stockCode = rand.Intn(100000) + 100000
	// book is initially not deleted
	book.isDeleted = false
	//id will be incremented for next book
	ID++
	return book
}

type Deletable interface {
	Delete()
}

//sets book isDeleted field to trueif not set already
func (book *Book) Delete() error {
	if book.isDeleted {
		return ErrBookNotFound
	}
	book.isDeleted = true
	fmt.Printf("Book: %s is deleted", book.name)
	return nil
}

// buy book with given count if stock is enough
func (book *Book) Buy(count int) error {
	if book.stockCount < count {
		return ErrNotEnoughStock
	}
	book.stockCount -= count
	fmt.Printf("Book: %s is buyed by user. New stockCount is %d", book.name, book.stockCount)
	return nil
}
