package library

import (
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name       string
	StockCode  int
	ISBN       int
	PageCount  int
	Price      float64
	StockCount int
	Author     string
	IsDeleted  bool
}

// Book constructor
func NewBook(name string, author string) *Book {
	book := new(Book)
	book.Name = name
	book.Author = author
	//Seed is current time to give randomness
	// page count will be in range 300-400
	book.PageCount = rand.Intn(300) + 100
	// price will be in range 20.00- 220.00
	book.Price = rand.Float64()*200 + 20
	// ISBN will be in range 100000 - 1000000
	book.ISBN = rand.Intn(100000) + 100000
	// stock count  will be in range 0-50
	book.StockCount = rand.Intn(50)
	// stock code  will be in range 100000 - 1000000
	book.StockCode = rand.Intn(100000) + 100000
	// book is initially not deleted
	book.IsDeleted = false
	//id will be incremented for next book
	return book
}

func GiveISBN(b Book) *Book {
	if b.ISBN == 0 {
		b.ISBN = rand.Intn(100000) + 100000
	}
	return &b
}

type Deletable interface {
	Delete()
}

//sets book isDeleted field to trueif not set already
func (book *Book) Delete() error {
	if book.IsDeleted {
		return ErrBookNotFound
	}
	book.IsDeleted = true
	fmt.Printf("Book: %s is deleted", book.Name)
	return nil
}

// buy book with given count if stock is enough
func (book *Book) Buy(count int) (string, error) {
	if book.StockCount < count {
		return "", ErrNotEnoughStock
	}
	book.StockCount -= count
	result := fmt.Sprintf("Book: %s is buyed by user. New stockCount is %d", book.Name, book.StockCount)
	fmt.Printf("Book: %s is buyed by user. New stockCount is %d", book.Name, book.StockCount)
	return result, nil
}
