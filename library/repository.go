package library

import (
	"errors"
	"fmt"

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
	var cities []Book
	r.db.Find(&cities)

	return cities
}

func (r *BookRepository) FindByCountryCode(countryCode string) []Book {
	var cities []Book
	r.db.Where("CountryCode = ?", countryCode).Order("Id desc,name").Find(&cities)

	// Struct
	//r.db.Where(&Book{CountryCode: countryCode}).First(&cities)
	//r.db.Where(map[string]interface{}{"CountryCode": countryCode, "Code": "01"}).Find(&cities)
	//r.db.Where([]int64{20, 21, 22}).Find(&cities) // ID IN(20,21,22)

	return cities
}

func (r *BookRepository) FindByCountryCodeOrBookCode(code string) []Book {
	var cities []Book
	r.db.Where("CountryCode = ?", code).Or("Code = ?", code).Find(&cities)
	return cities
}

func (r *BookRepository) FindByName(name string) []Book {
	var cities []Book
	r.db.Where("Name LIKE ?", "%"+name+"%").Find(&cities)

	return cities
}

func (r *BookRepository) FindByNameWithRawSQL(name string) []Book {
	var cities []Book
	r.db.Raw("SELECT * FROM Book WHERE Name LIKE ?", "%"+name+"%").Scan(&cities)

	return cities
}

func (r *BookRepository) GetById(id int) Book {
	var Book Book
	result := r.db.First(&Book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Book not found with id : %d", id)
		return Book
	}
	return Book
}

func (r *BookRepository) Create(c Book) error {
	result := r.db.Create(c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookRepository) Update(c Book) error {
	result := r.db.Save(c)
	//r.db.Model(&c).Update("name", "deneme")

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BookRepository) Delete(c Book) error {
	result := r.db.Delete(c)

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
	//https://gorm.io/docs/migration.html#content-inner
}

func (r *BookRepository) InsertSampleData() {

	for _, b := range Books {
		r.db.Where(Book{ISBN: b.ISBN, Name: b.Name}).Attrs(Book{ISBN: b.ISBN, Name: b.Name}).FirstOrCreate(&b)
	}
}
