package user

import (
	"log"

	"github.com/tzincker/gocourse_web/internal/domain"
)

type (
	Filters struct {
		FirstName string
		LastName  string
	}

	Service interface {
		Create(firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Get(id string) (*domain.User, error)
		Delete(id string) error
		Update(id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName, lastName, email, phone string) (*domain.User, error) {
	log.Println("Create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	u, err := s.repo.Create(&user)

	if err != nil {
		s.log.Println(err)
	}

	return u, err
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	log.Println("Get all users service")

	users, err := s.repo.GetAll(filters, offset, limit)

	if err != nil {
		s.log.Println(err)
	}

	return users, err
}

func (s service) Get(id string) (*domain.User, error) {
	log.Println("Get user service")

	u, err := s.repo.Get(id)

	if err != nil {
		s.log.Println(err)
	}

	return u, err
}

func (s service) Delete(id string) error {
	log.Println("Delete user service")

	err := s.repo.Delete(id)

	if err != nil {
		s.log.Println(err)
	}

	return nil
}

func (s service) Update(
	id string,
	firstName *string,
	lastName *string,
	email *string,
	phone *string,
) error {
	log.Println("Update user service")
	err := s.repo.Update(id, firstName, lastName, email, phone)
	if err != nil {
		s.log.Println(err)
	}

	return nil
}

func (s service) Count(filters Filters) (int64, error) {
	log.Println("Get all users count service")
	count, err := s.repo.Count(filters)
	if err != nil {
		s.log.Println(err)
	}

	return count, err
}
