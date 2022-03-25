package library

import (
	"errors"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) FindAll() []Book {
	var books []Book
	r.db.Find(&books)

	return books
}

func (r *BookRepository) FindByISBN(ISBN string) []Book {
	var books []Book
	r.db.Where("ISBN = ?", ISBN).Order("Id desc,name").Find(&books)
	return books
}

func (r *BookRepository) FindByName(name string) []Book {
	var books []Book
	r.db.Where("Name LIKE ?", "%"+name+"%").Find(&books)

	return books
}

func (r *BookRepository) GetById(id int) *Book {
	var Book Book
	result := r.db.First(&Book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &Book
}

func (r *BookRepository) Create(b Book) error {
	result := r.db.Create(GiveISBN(b))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookRepository) Update(b Book) error {
	result := r.db.Where(Book{ISBN: b.ISBN}).Updates(&b)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookRepository) Delete(b Book) error {
	result := r.db.Delete(b)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *BookRepository) DeleteById(id int) error {
	result := r.db.Delete(&Book{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookRepository) Migration() {
	r.db.AutoMigrate(&Book{})
}

func (r *BookRepository) InsertSampleData() {

	for _, b := range Books {
		r.db.Where(Book{ISBN: b.ISBN, Name: b.Name}).Attrs(Book{ISBN: b.ISBN, Name: b.Name}).FirstOrCreate(&b)
	}
}
