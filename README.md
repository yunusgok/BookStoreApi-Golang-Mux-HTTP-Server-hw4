# HW2- Library- YUNUS EMRE GÖK

This project is a basic library application

## Usage
### List
To list all books in the library
```bash
go run main.go list
```
Output
```b
Moby Dick
War and Peace
Hamlet
The Odyssey
Madame Bovary
The Divine Comedy
Lolita
The Brothers Karamazov
.
.
.
```
### Search
To search books with specific name, author or ISBN no.
Prints name and ISBN of the books
```bash
go run main.go search <Words to be searched>
```
##### Example 1
```bash
go run main.go search the
```
```
The Odyssey -- ISBN: 131445 
The Divine Comedy -- ISBN: 132888 
The Brothers Karamazov -- ISBN: 160631
The Catcher in the Rye -- ISBN: 164324
The Adventures of Huckleberry Finn -- ISBN: 104538
The Iliad -- ISBN: 189002
To the Lighthouse -- ISBN: 107463
The Sound and the Fury -- ISBN: 172546
The Grapes of Wrath -- ISBN: 118591
The Trial -- ISBN: 144885
The Red and the Black -- ISBN: 192818
The Stories of Anton Chekhov -- ISBN: 177578
The Stranger -- ISBN: 104208
```

##### Example 2
```bash
go run main.go search 104208
```
```
The Stranger -- ISBN: 104208 
```
### Buy
To buy books with id and count. Gives error if stockCount is less than count
```bash
go run main.go buy <id> <count>
```
##### Example 1
```bash
go run main.go buy 1 5 
```

```
Book: War and Peace is buyed by user. New stockCount is 37
```
##### Example 2
```bash
go run main.go buy 1 75 
```
```
not enough stock
```
### Delete
To delete a book with id. Gives error if book is not exist
```bash
go run main.go delete  <id> 
```
##### Example 1
```bash
go run main.go delete 1
```
```
Book: War and Peace is deleted
```
##### Example 2
```bash
go run main.go delete 76
```
```
Book not found
```


## License
[MIT](https://choosealicense.com/licenses/mit/)