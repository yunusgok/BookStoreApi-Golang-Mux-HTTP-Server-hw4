package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/yunusgok/go-patika/library"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/yunusgok/go-patika/helpers"
)

func init() {
	library.InitBooks()
	library.InitRepo()
}

func main() {
	r := mux.NewRouter()

	CORSOptions()
	// tüm response'ları sıkıştırmak için
	r.Use(loggingMiddleware)
	r.Use(authenticationMiddleware)

	s := r.PathPrefix("/books").Subrouter()
	s.HandleFunc("/", ListBookHandler)
	s.HandleFunc("/{name}", NameSearchHandler)

	s.HandleFunc("/id/{id:[0-9]+}", FindByIdHandler)
	s.HandleFunc("/id/{id:[0-9]+}/buy", BuyBookHandler)
	s.HandleFunc("/id/{id:[0-9]+}/delete", DeleteBookHandler)

	p := r.PathPrefix("/book").Subrouter()
	p.HandleFunc("/create", bookCreate).Methods("POST")
	p.HandleFunc("/update", bookUpdate).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:8090",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)
}

// Returns book with given id in body
func FindByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	isInt, value := library.IsInt(vars["id"])
	if !isInt {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: library.ErrInvalidInput.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	book, err := library.FindBook(value)
	if err != nil {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: err.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	d := ApiResponse{
		Data: book,
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

// Process buy operation with given book id and count
func BuyBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	count := r.URL.Query().Get("count")
	idIsInt, idValue := library.IsInt(vars["id"])
	countIsInt, countValue := library.IsInt(count)
	if !idIsInt || !countIsInt {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: library.ErrInvalidInput.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	res, err := library.Buy(idValue, countValue)
	if err != nil {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: err.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	d := ApiResponse{
		Data: res,
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

// Deletes book with given book id
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	idIsInt, idValue := library.IsInt(vars["id"])
	if !idIsInt {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: library.ErrInvalidInput.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	err := library.DeleteBook(idValue)
	if err != nil {
		mr := &helpers.MalformedRequest{Status: http.StatusBadRequest, Msg: err.Error()}
		http.Error(w, mr.Msg, mr.Status)
		return
	}
	d := ApiResponse{
		Data: "Book succesfully deleted",
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

// Returns list of all books
func ListBookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	d := ApiResponse{
		Data: library.ListBooks(),
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

// Returns list of books matches with given name param
func NameSearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	d := ApiResponse{
		Data: library.FindBooks(vars["name"]),
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

// creates book with given book data
func bookCreate(w http.ResponseWriter, r *http.Request) {
	var b library.Book

	err := helpers.DecodeJSONBody(w, r, &b)

	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	library.CreateBook(b)

	fmt.Fprintf(w, "Book: %+v", b)
}

//updates book which matches with ISBN
func bookUpdate(w http.ResponseWriter, r *http.Request) {
	var b library.Book

	err := helpers.DecodeJSONBody(w, r, &b)

	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	library.UpdateBook(b)

	fmt.Fprintf(w, "Book: %+v", b)
}

type ApiResponse struct {
	Data interface{} `json:"data"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.URL.Query())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/book") {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token not found", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func CORSOptions() {
	handlers.AllowedOrigins([]string{"https://www.example.com"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})
}
