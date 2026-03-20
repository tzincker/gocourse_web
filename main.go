package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tzincker/gocourse_web/internal/course"
	"github.com/tzincker/gocourse_web/internal/enrollment"
	"github.com/tzincker/gocourse_web/internal/user"
	"github.com/tzincker/gocourse_web/pkg/bootstrap"
)

func main() {
	_ = godotenv.Load()
	url := bootstrap.Url()

	router := mux.NewRouter()

	log := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewRepo(log, db)
	userSrv := user.NewService(log, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	courseRepo := course.NewRepo(log, db)
	courseSrv := course.NewService(log, courseRepo)
	courseEnd := course.MakeEndpoints(courseSrv)

	enrollRepo := enrollment.NewRepo(log, db)
	enrollSrv := enrollment.NewService(log, userSrv, courseSrv, enrollRepo)
	enrollEnd := enrollment.MakeEndpoints(enrollSrv)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	router.HandleFunc("/courses", courseEnd.Create).Methods("POST")
	router.HandleFunc("/courses", courseEnd.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEnd.Get).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEnd.Update).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseEnd.Delete).Methods("DELETE")

	router.HandleFunc("/enrollments", enrollEnd.Create).Methods("POST")

	srv := &http.Server{
		Handler:      http.TimeoutHandler(router, 5*time.Second, "Server timed out"),
		Addr:         url,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Printf("Server listening to: %s\n", url)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
