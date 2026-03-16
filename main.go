package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tzincker/gocourse_web/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	router := mux.NewRouter()
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	fmt.Println(dsn)

	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = db.Debug()
	db.AutoMigrate(&user.User{})

	userSrc := user.NewService()

	userEnd := user.MakeEndpoints(userSrc)
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	url := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	srv := &http.Server{
		Handler:      http.TimeoutHandler(router, 5*time.Second, "Server timed out"),
		Addr:         url,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Printf("Server listening to: %s\n", url)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
