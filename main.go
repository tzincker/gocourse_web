package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tzincker/gocourse_web/internal/user"


	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	userSrc := user.NewService()

	port := ":8000"
	userEnd := user.MakeEndpoints(userSrc)
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	url := fmt.Sprintf("127.0.0.1%s", port)

	srv := &http.Server{ 
		Handler: http.TimeoutHandler(router, 5 * time.Second, "Server timed out"),
		Addr: url,
		WriteTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
	}
	log.Printf("Server listening to: %s\n", url)
	err :=  srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
