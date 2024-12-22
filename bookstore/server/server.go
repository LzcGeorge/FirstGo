package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	bss := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/book", bss.createBook).Methods("POST")
	router.HandleFunc("/book/{id}", bss.getBook).Methods("GET")
	router.HandleFunc("/book/{id}", bss.updateBook).Methods("POST")
	router.HandleFunc("book", bss.listBooks).Methods("GET")
	router.HandleFunc("/book/{id}", bss.deleteBook).Methods("DELETE")

	bss.srv.Handler = middleware.Logging(middleware.Validating(router))
	return bss
}

func (bss *BookStoreServer) createBook(w http.ResponseWriter, r *http.Request) {
	// 解码
	dec := json.NewDecoder(r.Body)
	var b store.Book
	if err := dec.Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 加入到 BookStoreServer中
	if err := bss.s.Create(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bss *BookStoreServer) getBook(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	b, err := bss.s.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response(w, b)
}

func (bss *BookStoreServer) updateBook(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	var b store.Book
	if err := dec.Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b.Id = id
	if err := bss.s.Update(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bss *BookStoreServer) listBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bss.s.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response(w, books)
}

func (bss *BookStoreServer) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := bss.s.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func response(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}
